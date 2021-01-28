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
	binanceDefaultWebsocketFutureURL          = "wss://fstream.binance.com/stream"
	binanceDefaultWebsocketFutureURLListenKey = "wss://fstream.binance.com/ws"
)

var listenKeyFuture string

// WsConnectFuture initiates a websocket connection
func (b *Binance) WsConnectFuture() error {
	if !b.WebsocketFuture.IsEnabled() || !b.IsEnabled() {
		return errors.New(stream.WebsocketNotEnabled)
	}
	var dialer websocket.Dialer
	var err error
	if b.WebsocketFuture.CanUseAuthenticatedEndpoints() {
		listenKeyFuture, err = b.GetWsAuthStreamKeyFuture()
		if err != nil {
			b.WebsocketFuture.SetCanUseAuthenticatedEndpoints(false)
			log.Errorf(log.ExchangeSys,
				"%v unable to connect to authenticated WebsocketFuture. Error: %s",
				b.Name,
				err)
		} else {
			authPayload := binanceDefaultWebsocketFutureURLListenKey + "/" + listenKeyFuture
			err = b.WebsocketFuture.SetWebsocketURL(authPayload, false, false)
			if err != nil {
				return err
			}
		}
	}

	err = b.WebsocketFuture.Conn.Dial(&dialer, http.Header{})
	if err != nil {
		return fmt.Errorf("%v - Unable to connect to WebsocketFuture. Error: %s",
			b.Name,
			err)
	}

	if b.WebsocketFuture.CanUseAuthenticatedEndpoints() {
		go b.KeepAuthKeyAliveFuture()
	}
	b.WebsocketFuture.Conn.SetupPingHandler(stream.PingHandler{
		UseGorillaHandler: true,
		MessageType:       websocket.PongMessage,
		Delay:             pingDelay,
	})

	// enabledPairs, err := b.GetEnabledPairs(asset.Future)
	// if err != nil {
	// 	return err
	// }

	// for i := range enabledPairs {
	// 	err = b.SeedLocalCache(enabledPairs[i])
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	go b.wsReadDataFuture()

	subs, err := b.GenerateSubscriptionsFuture()
	if err != nil {
		return err
	}

	return b.WebsocketFuture.SubscribeToChannels(subs)
}

// GetWsAuthStreamKeyFuture will retrieve a key to use for authorised WS streaming
func (b *Binance) GetWsAuthStreamKeyFuture() (string, error) {
	var resp UserAccountStream
	path := futureApiURL + userAccountFutureStream
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

// MaintainWsAuthStreamKeyFuture will keep the key alive
func (b *Binance) MaintainWsAuthStreamKeyFuture() error {
	var err error
	if listenKey == "" {
		listenKey, err = b.GetWsAuthStreamKeyFuture()
		return err
	}
	path := futureApiURL + userAccountFutureStream
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

// KeepAuthKeyAliveFuture will continuously send messages to
// keep the WS auth key active
func (b *Binance) KeepAuthKeyAliveFuture() {
	b.WebsocketFuture.Wg.Add(1)
	defer b.WebsocketFuture.Wg.Done()
	ticks := time.NewTicker(time.Minute * 30)
	for {
		select {
		case <-b.WebsocketFuture.ShutdownC:
			ticks.Stop()
			return
		case <-ticks.C:
			err := b.MaintainWsAuthStreamKeyFuture()
			if err != nil {
				b.WebsocketFuture.DataHandler <- err
				log.Warnf(log.ExchangeSys,
					b.Name+" - Unable to renew auth websocketFuture token, may experience shutdown")
			}
		}
	}
}

// wsReadData receives and passes on websocket messages for processing
func (b *Binance) wsReadDataFuture() {
	b.WebsocketFuture.Wg.Add(1)
	defer b.WebsocketFuture.Wg.Done()

	for {
		resp := b.WebsocketFuture.Conn.ReadMessage()
		if resp.Raw == nil {
			return
		}
		err := b.wsHandleDataFuture(resp.Raw)
		if err != nil {
			b.WebsocketFuture.DataHandler <- err
		}
	}
}

func (b *Binance) wsHandleDataFuture(respRaw []byte) error {
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

		pairs, err := b.GetEnabledPairs(asset.Future)
		if err != nil {
			return err
		}

		format, err := b.GetPairFormat(asset.Future, true)

		if err != nil {
			return err
		}

		switch e {
		case "ACCOUNT_UPDATE":
			// var data wsAccountInfo
			// err := json.Unmarshal(respRaw, &data)
			// if err != nil {
			// 	return fmt.Errorf("%v - Could not convert to outboundAccountInfo structure %s",
			// 		b.Name,
			// 		err)
			// }
			// b.WebsocketFuture.DataHandler <- data
			fmt.Printf("账户更新事件:%+v\n", string(respRaw))
		case "ORDER_TRADE_UPDATE":

			fmt.Printf("订单/交易 更新推送:%+v\n", string(respRaw))
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
			b.WebsocketFuture.DataHandler <- MarkPriceStreamResponse{
				Symbol:               pair,
				EventType:            _stream.EventType,
				EventTime:            time.Unix(0, _stream.EventTime*int64(time.Millisecond)),
				MarkPrice:            _stream.MarkPrice,
				IndexPrice:           _stream.IndexPrice,
				EstimatedSettlePrice: _stream.EstimatedSettlePrice,
				LastFundingRate:      _stream.LastFundingRate,
				NextFundingTime:      time.Unix(0, _stream.NextFundingTime*int64(time.Millisecond)),
				AssetType:            asset.Future,
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
			b.WebsocketFuture.DataHandler <- stream.KlineData{
				Timestamp:  time.Unix(0, kline.EventTime*int64(time.Millisecond)),
				Pair:       pair,
				AssetType:  asset.Future,
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

// GenerateSubscriptionsFuture generates the default subscription set
func (b *Binance) GenerateSubscriptionsFuture() ([]stream.ChannelSubscription, error) {
	var channels = []string{"@markPrice", "@kline_1m", "@forceOrder"}
	var subscriptions []stream.ChannelSubscription
	// assets := b.GetAssetTypes()
	assetType := asset.Future
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

// SubscribeFuture subscribes to a set of channels
func (b *Binance) SubscribeFuture(channelsToSubscribe []stream.ChannelSubscription) error {
	payload := WsPayload{
		Method: "SUBSCRIBE",
	}
	for i := range channelsToSubscribe {
		payload.Params = append(payload.Params, channelsToSubscribe[i].Channel)
	}
	err := b.WebsocketFuture.Conn.SendJSONMessage(payload)
	if err != nil {
		return err
	}

	b.WebsocketFuture.AddSuccessfulSubscriptions(channelsToSubscribe...)
	return nil
}

// UnsubscribeFuture unsubscribes from a set of channels
func (b *Binance) UnsubscribeFuture(channelsToUnsubscribe []stream.ChannelSubscription) error {
	payload := WsPayload{
		Method: "UNSUBSCRIBE",
	}
	for i := range channelsToUnsubscribe {
		payload.Params = append(payload.Params, channelsToUnsubscribe[i].Channel)
	}
	err := b.WebsocketFuture.Conn.SendJSONMessage(payload)
	if err != nil {
		return err
	}
	b.WebsocketFuture.RemoveSuccessfulUnsubscriptions(channelsToUnsubscribe...)
	return nil
}
