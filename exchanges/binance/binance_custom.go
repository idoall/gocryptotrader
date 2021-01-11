package binance

import (
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
	"github.com/idoall/gocryptotrader/exchanges/order"
)

// PositionRiskFuture 用户持仓风险V2 (USER_DATA)
func (b *Binance) PositionRiskFuture(symbol string) ([]PositionRiskResponse, error) {

	path := futureApiURL + binanceFuturePositionRisk

	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))

	var result []PositionRiskResponse
	var resp []interface{}
	var err error
	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}
	for _, v := range resp {
		p := PositionRiskResponse{}

		mapObj := v.(map[string]interface{})

		p.Symbol = mapObj["symbol"].(string)
		if p.PositionAmt, err = strconv.ParseFloat(mapObj["positionAmt"].(string), 64); err != nil {
			return nil, err
		}
		if p.EntryPrice, err = strconv.ParseFloat(mapObj["entryPrice"].(string), 64); err != nil {
			return nil, err
		}
		if p.MarkPrice, err = strconv.ParseFloat(mapObj["markPrice"].(string), 64); err != nil {
			return nil, err
		}
		if p.LiquidationPrice, err = strconv.ParseFloat(mapObj["liquidationPrice"].(string), 64); err != nil {
			return nil, err
		}
		if p.Leverage, err = strconv.ParseInt(mapObj["leverage"].(string), 10, 64); err != nil {
			return nil, err
		}
		if p.MaxNotionalValue, err = strconv.ParseInt(mapObj["maxNotionalValue"].(string), 10, 64); err != nil {
			return nil, err
		}
		p.MarginType = MarginType(mapObj["marginType"].(string))
		if p.IsolatedMargin, err = strconv.ParseFloat(mapObj["isolatedMargin"].(string), 64); err != nil {
			return nil, err
		}
		if p.IsAutoAddMargin, err = strconv.ParseBool(mapObj["isAutoAddMargin"].(string)); err != nil {
			return nil, err
		}
		p.PositionSide = PositionSide(mapObj["positionSide"].(string))

		result = append(result, p)
	}

	return result, nil
}

// CancelExistingOrderFuture sends a cancel order to Binance
// 取消订单
func (b *Binance) CancelExistingOrderFuture(symbol string, orderID int64, origClientOrderID string) (CancelOrderResponse, error) {
	var resp CancelOrderResponse

	path := futureApiURL + binanceFutureCancelOrder

	params := url.Values{}
	params.Set("symbol", symbol)

	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}

	if origClientOrderID != "" {
		params.Set("origClientOrderId", origClientOrderID)
	}

	return resp, b.SendAuthHTTPRequest(http.MethodDelete, path, params, limitOrder, &resp)
}

// IncomeFuture 获取账户损益资金流水
func (b *Binance) IncomeFuture(req FutureIncomeRequest) ([]FutureIncomeResponse, error) {

	path := futureApiURL + binanceFutureIncome

	params := url.Values{}
	if req.Symbol != "" {
		params.Set("symbol", strings.ToUpper(req.Symbol))
	}
	if req.IncomeType != IncomeType_ALL {
		params.Set("incomeType", string(req.IncomeType))
	}
	if req.StartTime != 0 {
		params.Set("startTime", strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		params.Set("endTime", strconv.FormatInt(req.EndTime, 10))
	}
	if req.Limit != 0 {
		params.Set("limit", strconv.FormatInt(req.Limit, 10))
	}

	var result []FutureIncomeResponse
	var resp []interface{}
	var err error
	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}

	for _, v := range resp {
		p := FutureIncomeResponse{}

		mapObj := v.(map[string]interface{})

		p.Symbol = mapObj["symbol"].(string)
		p.IncomeType = IncomeType(mapObj["incomeType"].(string))
		if p.Income, err = strconv.ParseFloat(mapObj["income"].(string), 64); err != nil {
			return nil, err
		}
		p.Asset = mapObj["asset"].(string)
		p.Info = mapObj["info"].(string)
		p.Time = time.Unix(0, int64(mapObj["time"].(float64))*int64(time.Millisecond))
		p.TranId = int64(mapObj["tranId"].(float64))
		p.TradeId = mapObj["tradeId"].(string)

		result = append(result, p)
	}

	return result, nil
}

// Transfer 用户万向划转
func (b *Binance) Transfer(transferType TransferType, symbolBase string, amount float64) (tranId int64, err error) {
	path := fmt.Sprintf("%s%s", apiURL, binanceTransfer)

	params := url.Values{}
	params.Set("type", string(transferType))
	params.Set("asset", symbolBase)
	params.Set("amount", strconv.FormatFloat(amount, 'f', -1, 64))

	type response struct {
		TranId int64 `json:'tranId'`
	}
	var resp response
	err = b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, &resp)
	return resp.TranId, err
}

// FutureLeverage 调整开仓杠杆
func (b *Binance) FutureLeverage(symbol string, leverage int) (*FutureLeverageResponse, error) {
	path := fmt.Sprintf("%s%s", futureApiURL, binanceFutureLeverage)

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("leverage", strconv.FormatInt(int64(leverage), 10))

	var resp interface{}
	err := b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, &resp)
	if err != nil {
		return nil, err
	}

	mapObj := resp.(map[string]interface{})

	result := &FutureLeverageResponse{}
	result.Symbol = mapObj["symbol"].(string)
	result.Leverage = int(mapObj["leverage"].(float64))
	if result.MaxNotionalValue, err = strconv.ParseInt(mapObj["maxNotionalValue"].(string), 10, 64); err != nil {
		return nil, err
	}
	return result, err
}

// QueryOrderFuture returns information on a past order
func (b *Binance) QueryOrderFuture(symbol string, orderID int64, origClientOrderID string) (FutureQueryOrderData, error) {

	type Response struct {
		AvgPrice      float64                    `json:"avgPrice,string"`    // 平均成交价
		ClientOrderID string                     `json:"clientOrderId"`      // 用户自定义的订单号
		CumQuote      float64                    `json:"cumQuote,string"`    //成交金额
		ExecutedQty   float64                    `json:"executedQty,string"` //成交量
		OrderID       int64                      `json:"orderId"`
		OrigQty       float64                    `json:"origQty,string"` // 原始委托数量
		OrigType      string                     `json:"origType"`
		Price         float64                    `json:"price,string"`
		ReduceOnly    bool                       `json:"reduceOnly"` // 是否仅减仓
		Side          order.Side                 `json:"side"`
		PositionSide  PositionSide               `json:"positionSide"` // 持仓方向
		Status        order.Status               `json:"status"`
		StopPrice     float64                    `json:"stopPrice,string"` // 触发价，对`TRAILING_STOP_MARKET`无效
		ClosePosition bool                       `json:"closePosition"`    // 是否条件全平仓
		Symbol        string                     `json:"symbol"`
		Time          float64                    `json:"time"`                 // 订单时间
		TimeInForce   RequestParamsTimeForceType `json:"timeInForce"`          // 有效方法
		Type          string                     `json:"type"`                 //订单类型
		ActivatePrice float64                    `json:"activatePrice,string"` // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
		PriceRate     float64                    `json:"priceRate,string"`     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
		UpdateTime    int64                      `json:"updateTime"`
		WorkingType   WorkingType                `json:"workingType"`  // 条件价格触发类型
		PriceProtect  bool                       `json:"priceProtect"` // 是否开启条件单触发保护
	}
	var resp Response

	path := futureApiURL + binanceFutureQueryOrder

	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	if origClientOrderID != "" {
		params.Set("origClientOrderId", origClientOrderID)
	}
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}

	var result FutureQueryOrderData
	if err := b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}
	result = FutureQueryOrderData{
		AvgPrice:      resp.AvgPrice,
		ClientOrderID: resp.ClientOrderID,
		CumQuote:      resp.CumQuote,
		ExecutedQty:   resp.ExecutedQty,
		OrderID:       resp.OrderID,
		OrigQty:       resp.OrigQty,
		OrigType:      resp.OrigType,
		Price:         resp.Price,
		ReduceOnly:    resp.ReduceOnly,
		Side:          resp.Side,
		PositionSide:  resp.PositionSide,
		Status:        resp.Status,
		StopPrice:     resp.StopPrice,
		ClosePosition: resp.ClosePosition,
		Symbol:        resp.Symbol,
		TimeInForce:   resp.TimeInForce,
		Type:          resp.Type,
		ActivatePrice: resp.ActivatePrice,
		PriceRate:     resp.PriceRate,
		WorkingType:   resp.WorkingType,
		PriceProtect:  resp.PriceProtect,
		Time:          time.Unix(0, int64(resp.Time)*int64(time.Millisecond)),
		UpdateTime:    time.Unix(0, resp.UpdateTime*int64(time.Millisecond)),
	}

	return result, nil
}

// OpenOrdersFuture Current open orders. Get all open orders on a symbol.
// Careful when accessing this with no symbol: The number of requests counted against the rate limiter
// is significantly higher
func (b *Binance) OpenOrdersFuture(symbol string) ([]FutureQueryOrderData, error) {
	type Response struct {
		AvgPrice      float64                    `json:"avgPrice,string"`    // 平均成交价
		ClientOrderID string                     `json:"clientOrderId"`      // 用户自定义的订单号
		CumQuote      float64                    `json:"cumQuote,string"`    //成交金额
		ExecutedQty   float64                    `json:"executedQty,string"` //成交量
		OrderID       int64                      `json:"orderId"`
		OrigQty       float64                    `json:"origQty,string"` // 原始委托数量
		OrigType      string                     `json:"origType"`
		Price         float64                    `json:"price,string"`
		ReduceOnly    bool                       `json:"reduceOnly"` // 是否仅减仓
		Side          order.Side                 `json:"side"`
		PositionSide  PositionSide               `json:"positionSide"` // 持仓方向
		Status        order.Status               `json:"status"`
		StopPrice     float64                    `json:"stopPrice,string"` // 触发价，对`TRAILING_STOP_MARKET`无效
		ClosePosition bool                       `json:"closePosition"`    // 是否条件全平仓
		Symbol        string                     `json:"symbol"`
		Time          float64                    `json:"time"`                 // 订单时间
		TimeInForce   RequestParamsTimeForceType `json:"timeInForce"`          // 有效方法
		Type          string                     `json:"type"`                 //订单类型
		ActivatePrice float64                    `json:"activatePrice,string"` // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
		PriceRate     float64                    `json:"priceRate,string"`     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
		UpdateTime    int64                      `json:"updateTime"`
		WorkingType   WorkingType                `json:"workingType"`  // 条件价格触发类型
		PriceProtect  bool                       `json:"priceProtect"` // 是否开启条件单触发保护
	}
	var respList []Response
	var result []FutureQueryOrderData

	path := futureApiURL + binanceFutureOpenOrders

	params := url.Values{}

	if symbol != "" {
		params.Set("symbol", strings.ToUpper(symbol))
	}

	if err := b.SendAuthHTTPRequest(http.MethodGet, path, params, openOrdersLimit(symbol), &respList); err != nil {
		return result, err
	}
	for _, resp := range respList {
		result = append(result, FutureQueryOrderData{
			AvgPrice:      resp.AvgPrice,
			ClientOrderID: resp.ClientOrderID,
			CumQuote:      resp.CumQuote,
			ExecutedQty:   resp.ExecutedQty,
			OrderID:       resp.OrderID,
			OrigQty:       resp.OrigQty,
			OrigType:      resp.OrigType,
			Price:         resp.Price,
			ReduceOnly:    resp.ReduceOnly,
			Side:          resp.Side,
			PositionSide:  resp.PositionSide,
			Status:        resp.Status,
			StopPrice:     resp.StopPrice,
			ClosePosition: resp.ClosePosition,
			Symbol:        resp.Symbol,
			TimeInForce:   resp.TimeInForce,
			Type:          resp.Type,
			ActivatePrice: resp.ActivatePrice,
			PriceRate:     resp.PriceRate,
			WorkingType:   resp.WorkingType,
			PriceProtect:  resp.PriceProtect,
			Time:          time.Unix(0, int64(resp.Time)*int64(time.Millisecond)),
			UpdateTime:    time.Unix(0, resp.UpdateTime*int64(time.Millisecond)),
		})
	}

	return result, nil
}

// NewFutureOrder sends a new order to Binance
func (b *Binance) NewOrderFuture(o *FutureNewOrderRequest) (resp *FutureNewOrderResponse, err error) {

	if resp, err = b.newOrderFuture(o); err != nil {
		return resp, err
	}
	return resp, nil
}

func (b *Binance) newOrderFuture(o *FutureNewOrderRequest) (result *FutureNewOrderResponse, err error) {
	path := fmt.Sprintf("%s%s", futureApiURL, binanceFutureNewOrder)

	params := url.Values{}
	params.Set("symbol", o.Symbol)
	params.Set("side", string(o.Side))
	params.Set("type", string(o.Type))
	params.Set("positionSide", string(o.PositionSide))

	params.Set("quantity", strconv.FormatFloat(o.Quantity, 'f', -1, 64))

	if o.Type == BinanceRequestParamsOrderLimit {
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
	type Response struct {
		Symbol        string  `json:"symbol"` //交易对
		OrderID       int64   `json:"orderId"`
		ClientOrderID string  `json:"clientOrderId"`
		AvgPrice      string  `json:"avgPrice, string"` //平均成交价
		Price         float64 `json:"price,string"`     //委托价格
		OrigQty       float64 `json:"origQty,string"`   //原始委托数量
		CumQty        float64 `json:"cumQty,string"`
		CumQuote      float64 `json:"cumQuote,string"` //成交金额
		// The cumulative amount of the quote that has been spent (with a BUY order) or received (with a SELL order).
		ExecutedQty   float64                    `json:"executedQty,string"` //成交量
		Status        order.Status               `json:"status"`             //订单状态
		TimeInForce   RequestParamsTimeForceType `json:"timeInForce"`        //有效方法
		Type          order.Type                 `json:"type"`               //订单类型
		Side          order.Side                 `json:"side"`               //买卖方向
		PositionSide  string                     `json:"positionSide"`       //持仓方向
		StopPrice     string                     `json:"stopPrice"`          //触发价，对`TRAILING_STOP_MARKET`无效
		ClosePosition bool                       `json:"closePosition"`      //是否条件全平仓
		OrigType      string                     `json:"origType"`           //触发前订单类型
		ActivatePrice string                     `json:"activatePrice"`      //跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
		PriceRate     string                     `json:"priceRate"`          //跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段

		UpdateTime   int64       `json:"updateTime"`   // 更新时间
		WorkingType  WorkingType `json:"workingType"`  // 条件价格触发类型
		PriceProtect bool        `json:"priceProtect"` // 是否开启条件单触发保护
	}
	var resp Response
	if err = b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, &resp); err != nil {
		return result, err
	}

	result = &FutureNewOrderResponse{
		Symbol:        resp.Symbol,
		OrderID:       resp.OrderID,
		ClientOrderID: resp.ClientOrderID,
		Price:         resp.Price,
		OrigQty:       resp.OrigQty,
		CumQty:        resp.CumQty,
		CumQuote:      resp.CumQuote,
		ExecutedQty:   resp.ExecutedQty,
		Status:        resp.Status,
		TimeInForce:   resp.TimeInForce,
		Type:          resp.Type,
		Side:          resp.Side,
		PositionSide:  resp.PositionSide,
		StopPrice:     resp.StopPrice,
		ClosePosition: resp.ClosePosition,
		OrigType:      resp.OrigType,
		ActivatePrice: resp.ActivatePrice,
		PriceRate:     resp.PriceRate,
		UpdateTime:    time.Unix(0, resp.UpdateTime*int64(time.Millisecond)),
	}
	if result.AvgPrice, err = strconv.ParseFloat(resp.AvgPrice, 64); err != nil {
		return nil, err
	}
	return result, nil
}

// GetFutureFundingRate 查询资金费率历史
func (b *Binance) GetFutureFundingRate(req FutureFundingRateRequest) ([]FutureFundingRateResponeItem, error) {

	params := url.Values{}
	if req.Symbol.String() != "" {
		params.Set("symbol", req.Symbol.String())
	}
	if req.Limit != 0 {
		params.Set("limit", strconv.FormatInt(req.Limit, 10))
	}
	if req.StartTime != 0 {
		params.Set("startTime", strconv.FormatInt(req.StartTime, 10))
	}
	if req.EndTime != 0 {
		params.Set("endTime", strconv.FormatInt(req.EndTime, 10))
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
func (b *Binance) GetHistoricCandlesFuture(pair currency.Pair, contractType ContractType, start, end time.Time, interval kline.Interval) (kline.Item, error) {
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
		contractType: contractType,
		StartTime:    start.Unix() * 1000,
		EndTime:      end.Unix() * 1000,
		Limit:        int(b.Features.Enabled.Kline.ResultLimit),
	}

	ret := kline.Item{
		Exchange: b.Name,
		Pair:     pair,
		Asset:    asset.Future,
		Interval: interval,
	}

	candles, err := b.GetSpotKlineFuture(req)
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
