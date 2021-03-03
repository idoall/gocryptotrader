package binance

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/idoall/gocryptotrader/common"
	"github.com/idoall/gocryptotrader/currency"
	"github.com/idoall/gocryptotrader/exchanges/asset"
	"github.com/idoall/gocryptotrader/exchanges/request"
	"github.com/idoall/gocryptotrader/exchanges/stream"
	"github.com/idoall/gocryptotrader/log"
)

const (
	binanceDefaultWebsocketPerpURL          = "wss://dstream.binance.com/stream"
	binanceDefaultWebsocketPerpURLListenKey = "wss://dstream.binance.com/ws"
)

var listenKeyPerp string

// WsConnectPerp initiates a websocket connection
func (b *Binance) WsConnectPerp() error {
	if !b.WebsocketPerp.IsEnabled() || !b.IsEnabled() {
		return errors.New(stream.WebsocketNotEnabled)
	}
	var dialer websocket.Dialer
	var err error
	if b.WebsocketPerp.CanUseAuthenticatedEndpoints() {
		listenKeyPerp, err = b.GetWsAuthStreamKeyPerp()
		if err != nil {
			b.WebsocketPerp.SetCanUseAuthenticatedEndpoints(false)
			log.Errorf(log.ExchangeSys,
				"%v unable to connect to authenticated WebsocketPerp. Error: %s",
				b.Name,
				err)
		} else {
			authPayload := binanceDefaultWebsocketPerpURLListenKey + "/" + listenKeyPerp
			err = b.WebsocketPerp.SetWebsocketURL(authPayload, false, false)
			if err != nil {
				return err
			}
		}
	}

	err = b.WebsocketPerp.Conn.Dial(&dialer, http.Header{})
	if err != nil {
		return fmt.Errorf("%v - Unable to connect to WebsocketPerp. Error: %s",
			b.Name,
			err)
	}

	if b.WebsocketPerp.CanUseAuthenticatedEndpoints() {
		go b.KeepAuthKeyAlivePerp()
	}
	b.WebsocketPerp.Conn.SetupPingHandler(stream.PingHandler{
		UseGorillaHandler: true,
		MessageType:       websocket.PongMessage,
		Delay:             pingDelay,
	})

	// enabledPairs, err := b.GetEnabledPairs(asset.PerpetualContract)
	// if err != nil {
	// 	return err
	// }

	// for i := range enabledPairs {
	// 	err = b.SeedLocalCache(enabledPairs[i])
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	go b.wsReadDataPerp()

	subs, err := b.GenerateSubscriptionsPerp()
	if err != nil {
		return err
	}

	return b.WebsocketPerp.SubscribeToChannels(subs)
}

// GetWsAuthStreamKeyPerp will retrieve a key to use for authorised WS streaming
func (b *Binance) GetWsAuthStreamKeyPerp() (string, error) {
	var resp UserAccountStream
	path := perpetualApiURL + userAccountPerpStream
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = b.API.Credentials.Key
	err := b.SendPayload(context.Background(), &request.Item{
		Method:        http.MethodPost,
		Path:          path,
		Headers:       headers,
		Body:          bytes.NewBuffer(nil),
		Result:        &resp,
		AuthRequest:   true,
		Verbose:       b.Verbose,
		HTTPDebugging: b.HTTPDebugging,
		HTTPRecording: b.HTTPRecording,
	})
	if err != nil {
		return "", err
	}
	return resp.ListenKey, nil
}

// MaintainWsAuthStreamKeyPerp will keep the key alive
func (b *Binance) MaintainWsAuthStreamKeyPerp() error {
	var err error
	if listenKey == "" {
		listenKey, err = b.GetWsAuthStreamKeyPerp()
		return err
	}
	path := perpetualApiURL + userAccountPerpStream
	params := url.Values{}
	params.Set("listenKey", listenKey)
	path = common.EncodeURLValues(path, params)
	headers := make(map[string]string)
	headers["X-MBX-APIKEY"] = b.API.Credentials.Key
	return b.SendPayload(context.Background(), &request.Item{
		Method:        http.MethodPut,
		Path:          path,
		Headers:       headers,
		Body:          bytes.NewBuffer(nil),
		AuthRequest:   true,
		Verbose:       b.Verbose,
		HTTPDebugging: b.HTTPDebugging,
		HTTPRecording: b.HTTPRecording,
	})
}

// KeepAuthKeyAlivePerp will continuously send messages to
// keep the WS auth key active
func (b *Binance) KeepAuthKeyAlivePerp() {
	b.WebsocketPerp.Wg.Add(1)
	defer b.WebsocketPerp.Wg.Done()
	ticks := time.NewTicker(time.Minute * 30)
	for {
		select {
		case <-b.WebsocketPerp.ShutdownC:
			ticks.Stop()
			return
		case <-ticks.C:
			err := b.MaintainWsAuthStreamKeyPerp()
			if err != nil {
				b.WebsocketPerp.DataHandler <- err
				log.Warnf(log.ExchangeSys,
					b.Name+" - Unable to renew auth websocketPerp token, may experience shutdown")
			}
		}
	}
}

// wsReadData receives and passes on websocket messages for processing
func (b *Binance) wsReadDataPerp() {
	b.WebsocketPerp.Wg.Add(1)
	defer b.WebsocketPerp.Wg.Done()

	for {
		resp := b.WebsocketPerp.Conn.ReadMessage()
		if resp.Raw == nil {
			return
		}
		err := b.wsHandleDataPerp(resp.Raw)
		if err != nil {
			b.WebsocketPerp.DataHandler <- err
		}
	}
}

func (b *Binance) wsHandleDataPerp(respRaw []byte) error {
	var multiStreamData map[string]interface{}
	err := json.Unmarshal(respRaw, &multiStreamData)
	if err != nil {
		return err
	}
	if method, ok := multiStreamData["method"].(string); ok {
		// TODO handle subscription handling
		if strings.EqualFold(method, "subscribe") {
			return nil
		}
		if strings.EqualFold(method, "unsubscribe") {
			return nil
		}
	}
	if e, ok := multiStreamData["e"].(string); ok {

		pairs, err := b.GetEnabledPairs(asset.PerpetualContract)
		if err != nil {
			return err
		}

		format, err := b.GetPairFormat(asset.PerpetualContract, true)

		if err != nil {
			return err
		}

		switch e {
		case "ACCOUNT_UPDATE":
			var data AccountUpdateStream
			err := json.Unmarshal(respRaw, &data)
			if err != nil {
				return fmt.Errorf("%v - Could not convert to ACCOUNT_UPDATE structure %s",
					b.Name,
					err)
			}

			var o AccountUpdateStreamResponse
			o.EventType = data.EventType
			o.EventTime = time.Unix(0, data.EventTime*int64(time.Millisecond))
			o.TimeStamp = time.Unix(0, data.TimeStamp*int64(time.Millisecond))
			o.AccountUpdateEvent.EventCause = data.AccountUpdateEvent.EventCause
			o.Exchange = b.GetName()
			o.AssetType = asset.PerpetualContract
			for _, v := range data.AccountUpdateEvent.Balance {
				o.AccountUpdateEvent.Balance = append(o.AccountUpdateEvent.Balance, AccountUpdateEventBalance{
					Asset:         v.Asset,
					RealyBalance:  v.RealyBalance,
					WalletBalance: v.WalletBalance,
				})
			}

			for _, v := range data.AccountUpdateEvent.Position {
				pair, err := currency.NewPairFromFormattedPairs(v.Symbol, pairs, format)
				if err != nil {
					return err
				}

				marginType := MarginType_CROSSED
				if strings.EqualFold(v.MarginType, "isolated") {
					marginType = MarginType_ISOLATED
				}
				o.AccountUpdateEvent.Position = append(o.AccountUpdateEvent.Position, AccountUpdateEventPosition{
					Symbol:                pair,
					PositionAmt:           v.PositionAmt,
					EntryPrice:            v.EntryPrice,
					RealizedProfitAndLoss: v.RealizedProfitAndLoss,
					UnRealizedProfit:      v.UnRealizedProfit,
					MarginType:            marginType,
					IsolatedMargin:        v.IsolatedMargin,
					PositionSide:          PositionSide(v.PositionSide),
				})
			}

			// fmt.Printf("账户更新事件:%+v\n", string(respRaw))
			// fmt.Printf("账户更新事件:%+v\n", o)
			b.WebsocketPerp.DataHandler <- o
		case "ORDER_TRADE_UPDATE":

			// fmt.Printf("订单/交易 更新推送:%+v\n", string(respRaw))
		case "markPriceUpdate":
			var _stream MarkPriceStream
			err := json.Unmarshal(respRaw, &_stream)
			if err != nil {

				return fmt.Errorf("%v - Could not convert to a MarkPriceStream structure %s",
					b.Name,
					err)
			}
			pair, err := currency.NewPairFromFormattedPairs(_stream.Symbol, pairs, format)
			if err != nil {
				return err
			}
			b.WebsocketPerp.DataHandler <- MarkPriceStreamResponse{
				Symbol:               pair,
				EventType:            _stream.EventType,
				EventTime:            time.Unix(0, _stream.EventTime*int64(time.Millisecond)),
				MarkPrice:            _stream.MarkPrice,
				IndexPrice:           _stream.IndexPrice,
				EstimatedSettlePrice: _stream.EstimatedSettlePrice,
				LastFundingRate:      _stream.LastFundingRate,
				NextFundingTime:      time.Unix(0, _stream.NextFundingTime*int64(time.Millisecond)),
				AssetType:            asset.PerpetualContract,
				Exchange:             b.Name,
			}
		case "kline":
			var kline KlineStream
			err := json.Unmarshal(respRaw, &kline)
			if err != nil {

				return fmt.Errorf("%v - Could not convert to a KlineStream structure %s",
					b.Name,
					err)
			}

			pair, err := currency.NewPairFromFormattedPairs(kline.Symbol, pairs, format)
			if err != nil {
				return err
			}
			b.WebsocketPerp.DataHandler <- stream.KlineData{
				Timestamp:  time.Unix(0, kline.EventTime*int64(time.Millisecond)),
				Pair:       pair,
				AssetType:  asset.PerpetualContract,
				Exchange:   b.Name,
				StartTime:  time.Unix(0, kline.Kline.StartTime*int64(time.Millisecond)),
				CloseTime:  time.Unix(0, kline.Kline.CloseTime*int64(time.Millisecond)),
				Interval:   kline.Kline.Interval,
				OpenPrice:  kline.Kline.OpenPrice,
				ClosePrice: kline.Kline.ClosePrice,
				HighPrice:  kline.Kline.HighPrice,
				LowPrice:   kline.Kline.LowPrice,
				Volume:     kline.Kline.Volume,
			}

		}
	}

	return nil
}

// GenerateSubscriptionsPerp generates the default subscription set
func (b *Binance) GenerateSubscriptionsPerp() ([]stream.ChannelSubscription, error) {
	var channels = []string{"@markPrice", "@kline_1m", "@forceOrder"}
	var subscriptions []stream.ChannelSubscription
	// assets := b.GetAssetTypes()
	assetType := asset.PerpetualContract
	// for x := range assets {
	pairs, err := b.GetEnabledPairs(assetType)
	if err != nil {
		return nil, err
	}

	for y := range pairs {
		for z := range channels {
			lp := pairs[y].Lower()
			lp.Delimiter = ""
			subscriptions = append(subscriptions, stream.ChannelSubscription{
				Channel:  lp.String() + channels[z],
				Currency: pairs[y],
				Asset:    assetType,
			})
		}
	}
	// }

	return subscriptions, nil
}

// SubscribePerp subscribes to a set of channels
func (b *Binance) SubscribePerp(channelsToSubscribe []stream.ChannelSubscription) error {
	payload := WsPayload{
		Method: "SUBSCRIBE",
	}
	for i := range channelsToSubscribe {
		payload.Params = append(payload.Params, channelsToSubscribe[i].Channel)
	}
	err := b.WebsocketPerp.Conn.SendJSONMessage(payload)
	if err != nil {
		return err
	}

	b.WebsocketPerp.AddSuccessfulSubscriptions(channelsToSubscribe...)
	return nil
}

// UnsubscribePerp unsubscribes from a set of channels
func (b *Binance) UnsubscribePerp(channelsToUnsubscribe []stream.ChannelSubscription) error {
	payload := WsPayload{
		Method: "UNSUBSCRIBE",
	}
	for i := range channelsToUnsubscribe {
		payload.Params = append(payload.Params, channelsToUnsubscribe[i].Channel)
	}
	err := b.WebsocketPerp.Conn.SendJSONMessage(payload)
	if err != nil {
		return err
	}
	b.WebsocketPerp.RemoveSuccessfulUnsubscriptions(channelsToUnsubscribe...)
	return nil
}
