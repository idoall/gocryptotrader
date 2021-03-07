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
	"github.com/idoall/gocryptotrader/exchanges/order"
)

// SpotSubmitOrder submits a new order
func (b *Binance) SpotSubmitOrder(s *SpotSubmit) (order.SubmitResponse, error) {
	var submitOrderResponse order.SubmitResponse

	var sideType string
	if s.Side == order.Buy {
		sideType = order.Buy.String()
	} else {
		sideType = order.Sell.String()
	}

	timeInForce := s.TimeInForce
	var requestParamsOrderType RequestParamsOrderType
	switch s.Type {
	case order.Market:
		timeInForce = ""
		requestParamsOrderType = BinanceRequestParamsOrderMarket
	case order.Limit:
		requestParamsOrderType = BinanceRequestParamsOrderLimit
	default:
		submitOrderResponse.IsOrderPlaced = false
		return submitOrderResponse, errors.New("unsupported order type")
	}

	fPair, err := b.FormatExchangeCurrency(s.Symbol, s.AssetType)
	if err != nil {
		return submitOrderResponse, err
	}

	var orderRequest = NewOrderRequest{
		Symbol:           fPair.String(),
		Side:             sideType,
		Price:            s.Price,
		Quantity:         s.Amount,
		TradeType:        requestParamsOrderType,
		TimeInForce:      timeInForce,
		NewClientOrderID: s.NewClientOrderId,
		StopPrice:        s.StopPrice,
		IcebergQty:       s.IcebergQty,
	}

	response, err := b.NewOrder(&orderRequest)
	if err != nil {
		return submitOrderResponse, err
	}

	if response.OrderID > 0 {
		submitOrderResponse.OrderID = strconv.FormatInt(response.OrderID, 10)
	}
	if response.ExecutedQty == response.OrigQty {
		submitOrderResponse.FullyMatched = true
	}
	submitOrderResponse.IsOrderPlaced = true

	for i := range response.Fills {
		submitOrderResponse.Trades = append(submitOrderResponse.Trades, order.TradeHistory{
			Price:    response.Fills[i].Price,
			Amount:   response.Fills[i].Qty,
			Fee:      response.Fills[i].Commission,
			FeeAsset: response.Fills[i].CommissionAsset,
		})
	}

	return submitOrderResponse, nil
}

// Ping 测试服务器连通性
func (b *Binance) Ping(assetType asset.Item) (ping bool, err error) {

	var path string
	if assetType == asset.Spot {
		path = b.API.Endpoints.URL + pingServer
	} else if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, "ping")
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, "ping")
	} else {
		return false, fmt.Errorf("Error assetType")
	}
	if err = b.SendHTTPRequest(path, limitDefault, nil); err != nil {
		return false, err
	}
	return true, nil
}

// ExchangeInfo returns exchange information. Check binance_types for more
// information
func (b *Binance) ExchangeInfo(assetType asset.Item) (ExchangeInfo, error) {
	var resp ExchangeInfo

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractExchangeInfo)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractExchangeInfo)
	} else {
		return resp, fmt.Errorf("Error assetType")
	}

	return resp, b.SendHTTPRequest(path, limitDefault, &resp)
}

// Account 账户信息V2 (USER_DATA)
func (b *Binance) Account(assetType asset.Item) (*AccountInfoFuture, error) {

	params := url.Values{}
	var resp AccountInfoFuture
	var err error
	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion2, binanceContractAccount)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractAccount)
	} else {
		return &resp, fmt.Errorf("Error assetType")
	}

	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ADLQuantile 持仓ADL队列估算
func (b *Binance) ADLQuantile(assetType asset.Item, symbol currency.Pair) (*AdlQuantileResponse, error) {
	var result *AdlQuantileResponse
	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractAdlQuantile)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractAdlQuantile)
	} else {
		return result, fmt.Errorf("Error assetType")
	}

	params := url.Values{}
	params.Set("symbol", symbol.String())

	var resp interface{}
	var err error
	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}
	p := &AdlQuantileResponse{}
	mapObj := resp.(map[string]interface{})

	if mapObj["symbol"] == nil {
		return nil, nil
	}
	p.Symbol = mapObj["symbol"].(string)

	adlObj := mapObj["adlQuantile"].(map[string]interface{})
	p.AdlQuantile.LONG = adlObj["LONG"].(float64)
	p.AdlQuantile.SHORT = adlObj["SHORT"].(float64)

	if adlObj["HEDGE"] != nil {
		p.AdlQuantile.HEDGE = adlObj["HEDGE"].(float64)
	}
	if mapObj["BOTH"] != nil {
		p.AdlQuantile.BOTH = adlObj["BOTH"].(float64)
	}

	return p, nil
}

// Income 获取账户损益资金流水
func (b *Binance) Income(assetType asset.Item, req IncomeRequest) ([]IncomeResponse, error) {

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

	var result []IncomeResponse
	var resp []interface{}
	var err error

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractIncome)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractIncome)
	} else {
		return result, fmt.Errorf("Error assetType")
	}

	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}

	for _, v := range resp {
		p := IncomeResponse{}

		mapObj := v.(map[string]interface{})

		p.Symbol = mapObj["symbol"].(string)
		p.IncomeType = IncomeType(mapObj["incomeType"].(string))
		if p.Income, err = strconv.ParseFloat(mapObj["income"].(string), 64); err != nil {
			return nil, err
		}
		p.Asset = mapObj["asset"].(string)
		p.Info = mapObj["info"].(string)
		p.Time = time.Unix(0, int64(mapObj["time"].(float64))*int64(time.Millisecond))
		if mapObj["tranId"] == nil {
			p.TranId = 0
		} else {
			p.TranId = int64(mapObj["tranId"].(float64))
		}
		if mapObj["tradeId"].(string) == "" {
			p.TradeId = 0
		} else {
			if p.TradeId, err = strconv.ParseInt(mapObj["tradeId"].(string), 10, 64); err != nil {
				return nil, err
			}
		}

		result = append(result, p)
	}

	return result, nil
}

// Leverage 调整开仓杠杆
func (b *Binance) Leverage(assetType asset.Item, symbol string, leverage int) (*FutureLeverageResponse, error) {

	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("leverage", strconv.FormatInt(int64(leverage), 10))

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceLeverage)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceLeverage)
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

	var resp interface{}
	err := b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, &resp)
	if err != nil {
		return nil, err
	}

	mapObj := resp.(map[string]interface{})

	result := &FutureLeverageResponse{}
	result.Symbol = mapObj["symbol"].(string)
	result.Leverage = int(mapObj["leverage"].(float64))
	if mapObj["maxNotionalValue"] != nil {
		if result.MaxNotionalValue, err = strconv.ParseInt(mapObj["maxNotionalValue"].(string), 10, 64); err != nil {
			return nil, err
		}
	}
	return result, err
}

// OpenOrdersContract Current open orders. Get all open orders on a symbol.
// Careful when accessing this with no symbol: The number of requests counted against the rate limiter
// is significantly higher
func (b *Binance) OpenOrdersContract(assetType asset.Item, symbol string) ([]FutureQueryOrderData, error) {
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

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractOpenOrders)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractOpenOrders)
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

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

// ForceOrders 用户强平单历史 (USER_DATA)
func (b *Binance) ForceOrders(assetType asset.Item, symbol string, startTime, endTime, limit int64) ([]FutureQueryOrderData, error) {
	type Response struct {
		OrderID       int64                      `json:"orderId"`
		Symbol        string                     `json:"symbol"`
		Status        order.Status               `json:"status"`
		ClientOrderID string                     `json:"clientOrderId"` // 用户自定义的订单号
		Price         float64                    `json:"price,string"`
		AvgPrice      float64                    `json:"avgPrice,string"`    // 平均成交价
		OrigQty       float64                    `json:"origQty,string"`     // 原始委托数量
		ExecutedQty   float64                    `json:"executedQty,string"` //成交量
		CumQuote      float64                    `json:"cumQuote,string"`    //成交金额
		TimeInForce   RequestParamsTimeForceType `json:"timeInForce"`        // 有效方法
		OrderType     string                     `json:"origType"`
		ReduceOnly    bool                       `json:"reduceOnly"`    // 是否仅减仓
		ClosePosition bool                       `json:"closePosition"` // 是否条件全平仓
		Side          order.Side                 `json:"side"`
		PositionSide  PositionSide               `json:"positionSide"`     // 持仓方向
		StopPrice     float64                    `json:"stopPrice,string"` // 触发价，对`TRAILING_STOP_MARKET`无效
		WorkingType   WorkingType                `json:"workingType"`      // 条件价格触发类型
		Time          float64                    `json:"time"`             // 订单时间
		UpdateTime    int64                      `json:"updateTime"`
	}
	var respList []Response
	var result []FutureQueryOrderData

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractForceOrder)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractForceOrder)
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

	params := url.Values{}

	if symbol != "" {
		params.Set("symbol", strings.ToUpper(symbol))
	}
	if startTime != 0 {
		params.Set("startTime", strconv.FormatInt(startTime, 10))
	}
	if endTime != 0 {
		params.Set("endTime", strconv.FormatInt(startTime, 10))
	}
	if limit != 0 {
		params.Set("limit", strconv.FormatInt(limit, 10))
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
			OrigType:      resp.OrderType,
			Price:         resp.Price,
			ReduceOnly:    resp.ReduceOnly,
			Side:          resp.Side,
			PositionSide:  resp.PositionSide,
			Status:        resp.Status,
			StopPrice:     resp.StopPrice,
			ClosePosition: resp.ClosePosition,
			Symbol:        resp.Symbol,
			TimeInForce:   resp.TimeInForce,
			WorkingType:   resp.WorkingType,
			Time:          time.Unix(0, int64(resp.Time)*int64(time.Millisecond)),
			UpdateTime:    time.Unix(0, resp.UpdateTime*int64(time.Millisecond)),
		})
	}
	return result, nil
}

// // NewFutureOrder sends a new order to Binance
// func (b *Binance) NewOrderFuture(o *NewOrderContractRequest) (resp *FutureNewOrderResponse, err error) {

// 	if resp, err = b.newOrderFuture(o); err != nil {
// 		return resp, err
// 	}
// 	return resp, nil
// }

// NewOrderContract sends a new order to Binance
func (b *Binance) NewOrderContract(assetType asset.Item, o *NewOrderContractRequest) (result *NewOrderContractResponse, err error) {

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractNewOrder)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractNewOrder)
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

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

	result = &NewOrderContractResponse{
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

// QueryOrderContract returns information on a past order
func (b *Binance) QueryOrderContract(assetType asset.Item, symbol string, orderID int64, origClientOrderID string) (FutureQueryOrderData, error) {

	type Response struct {
		AvgPrice      float64                    `json:"avgPrice,string"`    // 平均成交价
		ClientOrderID string                     `json:"clientOrderId"`      // 用户自定义的订单号
		CumBase       float64                    `json:"cumBase,string"`     //成交额(标的数量)
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
	var result FutureQueryOrderData

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractQueryOrder)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractQueryOrder)
	} else {
		return result, fmt.Errorf("Error assetType")
	}

	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))
	if origClientOrderID != "" {
		params.Set("origClientOrderId", origClientOrderID)
	}
	if orderID != 0 {
		params.Set("orderId", strconv.FormatInt(orderID, 10))
	}

	if err := b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}
	result = FutureQueryOrderData{
		AvgPrice:      resp.AvgPrice,
		ClientOrderID: resp.ClientOrderID,
		CumBase:       resp.CumBase,
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

// CancelExistingOrderContract sends a cancel order to Binance
// 取消订单
func (b *Binance) CancelExistingOrderContract(assetType asset.Item, symbol string, orderID int64, origClientOrderID string) (CancelOrderResponse, error) {
	var resp CancelOrderResponse

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceCancelOrder)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceCancelOrder)
	} else {
		return resp, fmt.Errorf("Error assetType")
	}

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

// GetPremiumIndex 最新标记价格和资金费率
func (b *Binance) GetPremiumIndex(assetType asset.Item, symbol currency.Pair) (*PreminuIndexResponse, error) {

	params := url.Values{}
	params.Set("symbol", symbol.String())

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s?%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractPreminuIndex, params.Encode())
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s?%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractPreminuIndex, params.Encode())
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

	var resp interface{}
	var err error
	if err = b.SendHTTPRequest(path, limitDefault, &resp); err != nil {
		return nil, err
	}

	p := new(PreminuIndexResponse)

	if assetType == asset.Future {
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
	} else if assetType == asset.PerpetualContract {
		mapObjArr := resp.([]interface{})
		mapObj := mapObjArr[0].(map[string]interface{})

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
	}

	return p, nil
}

// GetFundingRate 查询资金费率历史
func (b *Binance) GetFundingRate(assetType asset.Item, req FundingRateRequest) ([]FundingRateResponeItem, error) {

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

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s?%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceContractFundingRate, params.Encode())
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s?%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceContractFundingRate, params.Encode())
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

	var resp []interface{}
	var err error
	if err = b.SendHTTPRequest(path, limitDefault, &resp); err != nil {
		return nil, err
	}

	var result []FundingRateResponeItem
	for _, v := range resp {
		p := FundingRateResponeItem{}

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

// PositionRisk 用户持仓风险V2 (USER_DATA)
func (b *Binance) PositionRisk(assetType asset.Item, symbol string) ([]PositionRiskResponse, error) {

	params := url.Values{}
	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion2, binancePositionRisk)
		params.Set("symbol", strings.ToUpper(symbol))
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binancePositionRisk)
		params.Set("pair", strings.ToUpper(symbol))
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

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
		if p.UnRealizedProfit, err = strconv.ParseFloat(mapObj["unRealizedProfit"].(string), 64); err != nil {
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
		if mapObj["maxNotionalValue"] != nil {
			if p.MaxNotionalValue, err = strconv.ParseInt(mapObj["maxNotionalValue"].(string), 10, 64); err != nil {
				return nil, err
			}
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

// MarginType 变换逐全仓模式
func (b *Binance) MarginType(assetType asset.Item, symbol currency.Pair, marginType MarginType) (flag bool, err error) {

	params := url.Values{}
	params.Set("symbol", symbol.String())
	params.Set("marginType", string(marginType))

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binanceMarginType)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binanceMarginType)
	} else {
		return false, fmt.Errorf("Error assetType")
	}

	type response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	var resp response
	err = b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, &resp)
	if strings.Index(err.Error(), "{\"code\":-4046,\"msg\":\"No need to change margin type.\"}") != -1 {
		return true, nil
	} else if !strings.EqualFold(err.Error(), "success") {
		return false, err
	}
	return true, nil
}

// PositionMargin 调整逐仓保证金
func (b *Binance) PositionMargin(assetType asset.Item, req PositionMarginRequest) (bool, error) {

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binancePositionMargin)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binancePositionMargin)
	} else {
		return false, fmt.Errorf("Error assetType")
	}

	params := url.Values{}
	params.Set("symbol", req.Symbol.String())
	params.Set("amount", strconv.FormatFloat(req.Amount, 'f', -1, 64))
	params.Set("type", strconv.FormatInt(int64(req.Type), 10))
	params.Set("positionSide", string(req.PositionSide))

	var resp interface{}
	err := b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, &resp)
	if strings.EqualFold(err.Error(), "Successfully modify position margin.") {
		return true, nil
	}
	return false, err
}

// PositionMarginHistory 逐仓保证金变动历史 (TRADE)
func (b *Binance) PositionMarginHistory(assetType asset.Item, req PositionMarginHistoryRequest) ([]PositionMarginHistoryResponse, error) {

	var path string
	if assetType == asset.Future { // U本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", futureApiURL, binanceFutureRESTBasePath, binanceAPIVersion, binancePositionMarginHistory)
	} else if assetType == asset.PerpetualContract { // 币本位合约
		path = fmt.Sprintf("%s/%s/v%s/%s", perpetualApiURL, binancePerpetualRESTBasePath, binanceAPIVersion, binancePositionMarginHistory)
	} else {
		return nil, fmt.Errorf("Error assetType")
	}

	params := url.Values{}
	params.Set("symbol", req.Symbol.String())
	if req.Type != 0 {
		params.Set("type", strconv.FormatInt(int64(req.Type), 10))
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

	var result []PositionMarginHistoryResponse
	var resp []interface{}
	var err error
	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}

	for _, v := range resp {
		p := PositionMarginHistoryResponse{}

		mapObj := v.(map[string]interface{})

		p.Asset = mapObj["asset"].(string)
		p.Symbol = mapObj["symbol"].(string)
		if p.Amount, err = strconv.ParseFloat(mapObj["amount"].(string), 64); err != nil {
			return nil, err
		}

		p.Type = PositionMarginType(int(mapObj["type"].(float64)))
		p.PositionSide = PositionSide(mapObj["positionSide"].(string))
		p.Time = time.Unix(0, int64(mapObj["time"].(float64))*int64(time.Millisecond))

		result = append(result, p)
	}

	return result, nil
}

// // FutureAccount 账户信息V2 (USER_DATA)
// func (b *Binance) FutureAccount() (*AccountInfoFuture, error) {

// 	path := futureApiURL + binanceFutureAccount

// 	params := url.Values{}

// 	var resp AccountInfoFuture
// 	var err error
// 	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
// 		return nil, err
// 	}
// 	return &resp, nil
// }

// ADLQuantile 持仓ADL队列估算
// func (b *Binance) ADLQuantile(symbol currency.Pair) (*AdlQuantileResponse, error) {

// 	path := futureApiURL + binanceFutureAdlQuantile

// 	params := url.Values{}
// 	params.Set("symbol", symbol.String())

// 	var result *AdlQuantileResponse
// 	var resp interface{}
// 	var err error
// 	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
// 		return result, err
// 	}
// 	p := &AdlQuantileResponse{}
// 	mapObj := resp.(map[string]interface{})

// 	if mapObj["symbol"] == nil {
// 		return nil, nil
// 	}
// 	p.Symbol = mapObj["symbol"].(string)

// 	adlObj := mapObj["adlQuantile"].(map[string]interface{})
// 	p.AdlQuantile.LONG = adlObj["LONG"].(float64)
// 	p.AdlQuantile.SHORT = adlObj["SHORT"].(float64)

// 	if adlObj["HEDGE"] != nil {
// 		p.AdlQuantile.HEDGE = adlObj["HEDGE"].(float64)
// 	}
// 	if mapObj["BOTH"] != nil {
// 		p.AdlQuantile.BOTH = adlObj["BOTH"].(float64)
// 	}

// 	return p, nil
// }

// GetCommissionRate 用户手续费率
func (b *Binance) GetCommissionRate(assetType asset.Item, symbol string) (CommissionRateResponse, error) {

	path := futureApiURL + binanceFutureTradeFee
	if assetType == asset.Spot {
		path = b.API.Endpoints.URL + binanceSpotTradeFee
	}

	params := url.Values{}
	params.Set("symbol", strings.ToUpper(symbol))

	var result CommissionRateResponse
	var resp []interface{}
	var err error
	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return result, err
	}

	// for _, v := range resp {
	// 	p := PositionRiskResponse{}

	// 	mapObj := v.(map[string]interface{})

	// 	p.Symbol = mapObj["symbol"].(string)
	// 	if p.PositionAmt, err = strconv.ParseFloat(mapObj["positionAmt"].(string), 64); err != nil {
	// 		return nil, err
	// 	}
	// 	if p.EntryPrice, err = strconv.ParseFloat(mapObj["entryPrice"].(string), 64); err != nil {
	// 		return nil, err
	// 	}
	// 	if p.MarkPrice, err = strconv.ParseFloat(mapObj["markPrice"].(string), 64); err != nil {
	// 		return nil, err
	// 	}
	// 	if p.LiquidationPrice, err = strconv.ParseFloat(mapObj["liquidationPrice"].(string), 64); err != nil {
	// 		return nil, err
	// 	}
	// 	if p.Leverage, err = strconv.ParseInt(mapObj["leverage"].(string), 10, 64); err != nil {
	// 		return nil, err
	// 	}
	// 	if p.MaxNotionalValue, err = strconv.ParseInt(mapObj["maxNotionalValue"].(string), 10, 64); err != nil {
	// 		return nil, err
	// 	}
	// 	p.MarginType = MarginType(mapObj["marginType"].(string))
	// 	if p.IsolatedMargin, err = strconv.ParseFloat(mapObj["isolatedMargin"].(string), 64); err != nil {
	// 		return nil, err
	// 	}
	// 	if p.IsAutoAddMargin, err = strconv.ParseBool(mapObj["isAutoAddMargin"].(string)); err != nil {
	// 		return nil, err
	// 	}
	// 	p.PositionSide = PositionSide(mapObj["positionSide"].(string))

	// 	result = append(result, p)
	// }

	return result, nil
}

// IncomeFuture 获取账户损益资金流水
// func (b *Binance) IncomeFuture(req FutureIncomeRequest) ([]FutureIncomeResponse, error) {

// 	path := futureApiURL + binanceFutureIncome

// 	params := url.Values{}
// 	if req.Symbol != "" {
// 		params.Set("symbol", strings.ToUpper(req.Symbol))
// 	}
// 	if req.IncomeType != IncomeType_ALL {
// 		params.Set("incomeType", string(req.IncomeType))
// 	}
// 	if req.StartTime != 0 {
// 		params.Set("startTime", strconv.FormatInt(req.StartTime, 10))
// 	}
// 	if req.EndTime != 0 {
// 		params.Set("endTime", strconv.FormatInt(req.EndTime, 10))
// 	}
// 	if req.Limit != 0 {
// 		params.Set("limit", strconv.FormatInt(req.Limit, 10))
// 	}

// 	var result []FutureIncomeResponse
// 	var resp []interface{}
// 	var err error
// 	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
// 		return result, err
// 	}

// 	for _, v := range resp {
// 		p := FutureIncomeResponse{}

// 		mapObj := v.(map[string]interface{})

// 		p.Symbol = mapObj["symbol"].(string)
// 		p.IncomeType = IncomeType(mapObj["incomeType"].(string))
// 		if p.Income, err = strconv.ParseFloat(mapObj["income"].(string), 64); err != nil {
// 			return nil, err
// 		}
// 		p.Asset = mapObj["asset"].(string)
// 		p.Info = mapObj["info"].(string)
// 		p.Time = time.Unix(0, int64(mapObj["time"].(float64))*int64(time.Millisecond))
// 		if mapObj["tranId"] == nil {
// 			p.TranId = 0
// 		} else {
// 			p.TranId = int64(mapObj["tranId"].(float64))
// 		}
// 		if mapObj["tradeId"].(string) == "" {
// 			p.TradeId = 0
// 		} else {
// 			if p.TradeId, err = strconv.ParseInt(mapObj["tradeId"].(string), 10, 64); err != nil {
// 				return nil, err
// 			}
// 		}

// 		result = append(result, p)
// 	}

// 	return result, nil
// }

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
// func (b *Binance) FutureLeverage(symbol string, leverage int) (*FutureLeverageResponse, error) {
// 	path := fmt.Sprintf("%s%s", futureApiURL, binanceFutureLeverage)

// 	params := url.Values{}
// 	params.Set("symbol", symbol)
// 	params.Set("leverage", strconv.FormatInt(int64(leverage), 10))

// 	var resp interface{}
// 	err := b.SendAuthHTTPRequest(http.MethodPost, path, params, limitOrder, &resp)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mapObj := resp.(map[string]interface{})

// 	result := &FutureLeverageResponse{}
// 	result.Symbol = mapObj["symbol"].(string)
// 	result.Leverage = int(mapObj["leverage"].(float64))
// 	if result.MaxNotionalValue, err = strconv.ParseInt(mapObj["maxNotionalValue"].(string), 10, 64); err != nil {
// 		return nil, err
// 	}
// 	return result, err
// }

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
