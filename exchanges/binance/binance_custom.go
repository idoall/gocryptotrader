package binance

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/idoall/gocryptotrader/common/convert"
	"github.com/idoall/gocryptotrader/currency"
	"github.com/idoall/gocryptotrader/exchanges/asset"
	"github.com/idoall/gocryptotrader/exchanges/kline"
)

// FutureLeverage 调整开仓杠杆
func (b *Binance) FutureLeverage(symbol string, leverage int, resp *FutureLeverageResponse) error {
	path := fmt.Sprintf("%s%s", futureApiURL, binanceFutureNewOrder)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("leverage", strconv.FormatInt(int64(leverage), 10))

	return b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, resp)
}

// QueryOrderFuture returns information on a past order
func (b *Binance) QueryOrderFuture(symbol, origClientOrderID string, orderID int64) (FutureQueryOrderData, error) {
	var resp FutureQueryOrderData

	path := futureApiURL + binanceFutureQueryOrder

	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	if origClientOrderID != "" {
		params.Set("origClientOrderId", origClientOrderID)
	}
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}

	if err := b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return resp, err
	}

	if resp.Code != 0 {
		return resp, errors.New(resp.Msg)
	}
	return resp, nil
}

// OpenOrdersFuture Current open orders. Get all open orders on a symbol.
// Careful when accessing this with no symbol: The number of requests counted against the rate limiter
// is significantly higher
func (b *Binance) OpenOrdersFuture(symbol string) ([]FutureQueryOrderData, error) {
	var resp []FutureQueryOrderData

	path := futureApiURL + binanceFutureOpenOrders

	params := url.Values{}

	if symbol != "" {
		params.Set("symbol", strings.ToUpper(symbol))
	}

	if err := b.SendAuthHTTPRequest(http.MethodGet, path, params, openOrdersLimit(symbol), &resp); err != nil {
		return resp, err
	}

	return resp, nil
}

// NewFutureOrder sends a new order to Binance
func (b *Binance) NewOrderFuture(o *FutureNewOrderRequest) (FutureNewOrderResponse, error) {
	var resp FutureNewOrderResponse
	if err := b.newOrderFuture(o, &resp); err != nil {
		return resp, err
	}

	if resp.Code != 0 {
		return resp, errors.New(resp.Msg)
	}

	return resp, nil
}

func (b *Binance) newOrderFuture(o *FutureNewOrderRequest, resp *FutureNewOrderResponse) error {
	path := fmt.Sprintf("%s%s", futureApiURL, binanceFutureNewOrder)

	params := url.Values{}
	params.Set("symbol", o.Symbol)
	params.Set("side", o.Side)
	params.Set("type", string(o.TradeType))

	params.Set("quantity", strconv.FormatFloat(o.Quantity, 'f', -1, 64))

	if o.TradeType == BinanceRequestParamsOrderLimit {
		params.Set("price", strconv.FormatFloat(o.Price, 'f', -1, 64))
	}
	if o.TimeInForce != "" {
		params.Set("timeInForce", string(o.TimeInForce))
	}

	if o.NewClientOrderID != "" {
		params.Set("newClientOrderID", o.NewClientOrderID)
	}

	if o.StopPrice != 0 {
		params.Set("stopPrice", strconv.FormatFloat(o.StopPrice, 'f', -1, 64))
	}

	if o.NewOrderRespType != "" {
		params.Set("newOrderRespType", o.NewOrderRespType)
	}
	return b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, resp)
}

// GetFutureFundingRate 查询资金费率历史
func (b *Binance) GetFutureFundingRate(symbol currency.Pair, start, end, limit int64) ([]FutureFundingRateResponeItem, error) {

	params := url.Values{}
	params.Set("symbol", symbol.String())
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
	}
	if start != 0 {
		params.Set("startTime", strconv.FormatInt(start, 10))
	}
	if end != 0 {
		params.Set("endTime", strconv.FormatInt(end, 10))
	}

	path := fmt.Sprintf("%s%s?%s", futureApiURL, binanceFutureFundingRate, params.Encode())

	var resp []interface{}
	var err error
	if err = b.SendHTTPRequest(path, limitDefault, &resp); err != nil {
		return nil, err
	}

	var result []FutureFundingRateResponeItem
	for _, v := range resp {
		p := FutureFundingRateResponeItem{}

		mapObj := v.(map[string]interface{})

		p.Symbol = mapObj["symbol"].(string)

		if p.FundingRate, err = strconv.ParseFloat(mapObj["fundingRate"].(string), 64); err != nil {
			return nil, err
		}
		p.FundingTime = time.Unix(0, int64(mapObj["fundingTime"].(float64))*int64(time.Millisecond))

		result = append(result, p)
	}

	return result, nil
}

// GetHistoricCandlesFuture returns candles between a time period for a set time interval
func (b *Binance) GetHistoricCandlesFuture(pair currency.Pair, a asset.Item, start, end time.Time, interval kline.Interval) (kline.Item, error) {
	// if err := b.ValidateKline(pair, a, interval); err != nil {
	// 	return kline.Item{}, err
	// }

	// if kline.TotalCandlesPerInterval(start, end, interval) > b.Features.Enabled.Kline.ResultLimit {
	// 	return kline.Item{}, errors.New(kline.ErrRequestExceedsExchangeLimits)
	// }

	// fpair, err := b.FormatExchangeCurrency(pair, a)
	// if err != nil {
	// 	return kline.Item{}, err
	// }
	req := KlinesContractRequestParams{
		Interval:     b.FormatExchangeKlineInterval(interval),
		Pair:         pair.String(),
		contractType: a,
		StartTime:    start.Unix() * 1000,
		EndTime:      end.Unix() * 1000,
		Limit:        int(b.Features.Enabled.Kline.ResultLimit),
	}

	ret := kline.Item{
		Exchange: b.Name,
		Pair:     pair,
		Asset:    a,
		Interval: interval,
	}

	candles, err := b.GetFutureSpotKline(req)
	if err != nil {
		return kline.Item{}, err
	}

	for x := range candles {
		ret.Candles = append(ret.Candles, kline.Candle{
			Time:   candles[x].OpenTime,
			Open:   candles[x].Open,
			High:   candles[x].High,
			Low:    candles[x].Low,
			Close:  candles[x].Close,
			Volume: candles[x].Volume,
		})
	}

	ret.SortCandlesByTimestamp(false)
	return ret, nil
}

// GetFuturePremiumIndex 最新标记价格和资金费率
func (b *Binance) GetFuturePremiumIndex(symbol currency.Pair) (*PreminuIndexResponse, error) {

	params := url.Values{}
	params.Set("symbol", symbol.String())

	path := fmt.Sprintf("%s%s?%s", futureApiURL, binanceFuturePreminuIndex, params.Encode())

	var resp interface{}
	var err error
	if err = b.SendHTTPRequest(path, limitDefault, &resp); err != nil {
		return nil, err
	}

	p := new(PreminuIndexResponse)

	mapObj := resp.(map[string]interface{})

	if p.MarkPrice, err = strconv.ParseFloat(mapObj["markPrice"].(string), 64); err != nil {
		return nil, err
	}
	if p.IndexPrice, err = strconv.ParseFloat(mapObj["indexPrice"].(string), 64); err != nil {
		return nil, err
	}
	if p.LastFundingRate, err = strconv.ParseFloat(mapObj["lastFundingRate"].(string), 64); err != nil {
		return nil, err
	}
	if p.InterestRate, err = strconv.ParseFloat(mapObj["interestRate"].(string), 64); err != nil {
		return nil, err
	}
	p.NextFundingTime = time.Unix(0, int64(mapObj["nextFundingTime"].(float64))*int64(time.Millisecond))
	p.Time = time.Unix(0, int64(mapObj["time"].(float64))*int64(time.Millisecond))

	return p, nil
}

// GetSpotKlineFuture returns candle stick data
// 获取 K 线数据
// symbol:
// limit:
// interval
func (b *Binance) GetSpotKlineFuture(arg KlinesContractRequestParams) ([]CandleStick, error) {
	var resp interface{}
	var klineData []CandleStick

	params := url.Values{}
	params.Set("pair", arg.Pair)
	params.Set("contractType", string(arg.contractType))
	params.Set("interval", arg.Interval)
	if arg.Limit != 0 {
		params.Set("limit", strconv.Itoa(arg.Limit))
	}
	if arg.StartTime != 0 {
		params.Set("startTime", strconv.FormatInt(arg.StartTime, 10))
	}
	if arg.EndTime != 0 {
		params.Set("endTime", strconv.FormatInt(arg.EndTime, 10))
	}

	path := fmt.Sprintf("%s%s?%s", futureApiURL, binanceFutureCandleStick, params.Encode())

	if err := b.SendHTTPRequest(path, limitDefault, &resp); err != nil {
		return klineData, err
	}

	for _, responseData := range resp.([]interface{}) {
		var candle CandleStick
		for i, individualData := range responseData.([]interface{}) {
			switch i {
			case 0:
				tempTime := individualData.(float64)
				var err error
				candle.OpenTime, err = convert.TimeFromUnixTimestampFloat(tempTime)
				if err != nil {
					return klineData, err
				}
			case 1:
				candle.Open, _ = strconv.ParseFloat(individualData.(string), 64)
			case 2:
				candle.High, _ = strconv.ParseFloat(individualData.(string), 64)
			case 3:
				candle.Low, _ = strconv.ParseFloat(individualData.(string), 64)
			case 4:
				candle.Close, _ = strconv.ParseFloat(individualData.(string), 64)
			case 5:
				candle.Volume, _ = strconv.ParseFloat(individualData.(string), 64)
			case 6:
				tempTime := individualData.(float64)
				var err error
				candle.CloseTime, err = convert.TimeFromUnixTimestampFloat(tempTime)
				if err != nil {
					return klineData, err
				}
			case 7:
				candle.QuoteAssetVolume, _ = strconv.ParseFloat(individualData.(string), 64)
			case 8:
				candle.TradeCount = individualData.(float64)
			case 9:
				candle.TakerBuyAssetVolume, _ = strconv.ParseFloat(individualData.(string), 64)
			case 10:
				candle.TakerBuyQuoteAssetVolume, _ = strconv.ParseFloat(individualData.(string), 64)
			}
		}
		klineData = append(klineData, candle)
	}
	return klineData, nil
}

// Ping 测试服务器连通性
func (b *Binance) Ping() (ping bool, err error) {

	path := b.API.Endpoints.URL + pingServer
	if err = b.SendHTTPRequest(path, limitDefault, nil); err != nil {
		return false, err
	}
	return true, nil
}

// GetAccountSnapshot 查询每日资产快照 (USER_DATA)
func (b *Binance) GetAccountSnapshot(arg AccountSnapshotRequest) (snapshot []AccountSnapshotResponse, err error) {

	type respObj struct {
		Code        int    `json:"code"`
		Msg         string `json:"msg"`
		SnapshotVos []struct {
			Type       string `json:"type"`
			UpdateTime int    `json:"updateTime"`
			Data       struct {
				TotalAssetOfBtc string `json:"totalAssetOfBtc"`
				Balances        []struct {
					Asset  string `json:"asset"`
					Free   string `json:"free"`
					Locked string `json:"locked"`
				}
			} `json:"data"`
		} `json:"snapshotVos"`
	}
	var resp respObj

	params := url.Values{}
	params.Set("type", strings.ToUpper(arg.Type.String()))
	if arg.Limit == 0 {
		params.Set("limit", "10")
	} else {
		params.Set("limit", strconv.FormatInt(arg.Limit, 10))
	}
	if arg.StartTime != 0 {
		params.Set("startTime", strconv.FormatInt(arg.StartTime, 10))
	}
	if arg.EndTime != 0 {
		params.Set("endTime", strconv.FormatInt(arg.EndTime, 10))
	}

	path := b.API.Endpoints.URL + accountSnapshot

	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return
	}

	for _, v := range resp.SnapshotVos {
		var totalAssetOfBtc float64
		if totalAssetOfBtc, err = strconv.ParseFloat(v.Data.TotalAssetOfBtc, 64); err != nil {
			return
		}
		var free float64
		if free, err = strconv.ParseFloat(v.Data.Balances[0].Free, 64); err != nil {
			return
		}
		// var locked int64
		// if locked, err = strconv.ParseInt(v.Data.Balances[0].Locked, 10, 64); err != nil {
		// 	return
		// }
		var locked float64
		if locked, err = strconv.ParseFloat(v.Data.Balances[0].Locked, 64); err != nil {
			return
		}
		var updateTime time.Time
		if updateTime, err = convert.TimeFromUnixTimestampFloat(float64(v.UpdateTime)); err != nil {
			return
		}
		snapshot = append(snapshot, AccountSnapshotResponse{
			Asset:           asset.Item(v.Type),
			Symbol:          v.Data.Balances[0].Asset,
			TotalAssetOfBtc: totalAssetOfBtc,
			Free:            free,
			Locked:          locked,
			UpdateTime:      updateTime,
		})
	}

	return
}
