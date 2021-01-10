package binance

import (
	"time"

	"github.com/idoall/gocryptotrader/exchanges/asset"
)

const (
	pingServer              = "/api/v3/ping"
	perpetualContractApiURL = "https://dapi.binance.com"
	futureApiURL            = "https://fapi.binance.com"

	// PERPETUAL
	PerpetualExchangeInfo = "/dapi/v1/exchangeInfo"
	futureExchangeInfo    = "/fapi/v1/exchangeInfo"

	binancePerpetualCandleStick         = "/dapi/v1/klines"
	binancePerpetualContractCandleStick = "/dapi/v1/continuousKlines"
	binanceFutureCandleStick            = "/fapi/v1/continuousKlines"

	// binanceFuturePreminuIndex 最新标记价格和资金费率
	binanceFuturePreminuIndex = "/fapi/v1/premiumIndex"
	// binanceFutureFundingRate 查询资金费率历史
	binanceFutureFundingRate = "/fapi/v1/fundingRate"
	//下单 (TRADE)
	binanceFutureNewOrder = "/fapi/v1/order"
	//查询订单 (TRADE)
	binanceFutureQueryOrder = "/fapi/v1/order"
	// 查看当前全部挂单
	binanceFutureOpenOrders = "/fapi/v1/openOrders"
	// 调整开仓杠杆
	binanceFutureLeverage = "/fapi/v1/leverage"
)

type FutureLeverageResponse struct {
	Leverage         int    `json:"leverage,string"`  // 平均成交价
	MaxNotionalValue int64  `json:"maxNotionalValue"` // 用户自定义的订单号
	Symbol           string `json:"symbol"`           //成交金额
}

// FutureQueryOrderData holds query order data
type FutureQueryOrderData struct {
	AvgPrice      float64 `json:"avgPrice,string"`    // 平均成交价
	ClientOrderID string  `json:"clientOrderId"`      // 用户自定义的订单号
	CumQuote      float64 `json:"cumQuote,string"`    //成交金额
	ExecutedQty   float64 `json:"executedQty,string"` //成交量
	OrderID       int64   `json:"orderId"`
	OrigQty       float64 `json:"origQty,string"` // 原始委托数量
	OrigType      string  `json:"origType"`
	Price         float64 `json:"price,string"`
	ReduceOnly    bool    `json:"reduceOnly"` // 是否仅减仓
	Side          string  `json:"side"`
	PositionSide  string  `json:"positionSide"` // 持仓方向
	Status        string  `json:"status"`
	StopPrice     float64 `json:"stopPrice,string"` // 触发价，对`TRAILING_STOP_MARKET`无效
	ClosePosition bool    `json:"closePosition"`    // 是否条件全平仓
	Symbol        string  `json:"symbol"`
	Time          float64 `json:"time"`                 // 订单时间
	TimeInForce   string  `json:"timeInForce"`          // 有效方法
	Type          string  `json:"type"`                 //订单类型
	ActivatePrice float64 `json:"activatePrice,string"` // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     float64 `json:"priceRate,string"`     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	UpdateTime    int64   `json:"updateTime"`
	WorkingType   string  `json:"workingType"`  // 条件价格触发类型
	PriceProtect  bool    `json:"priceProtect"` // 是否开启条件单触发保护

	Code int    `json:"code"`
	Msg  string `json:"msg"`
	// // StopPrice           float64 `json:"stopPrice,string"`
	// IcebergQty          float64 `json:"icebergQty,string"`
	// IsWorking           bool    `json:"isWorking"`
	// CummulativeQuoteQty float64 `json:"cummulativeQuoteQty,string"`
	// OrderListID         int64   `json:"orderListId"`
	// OrigQuoteOrderQty   float64 `json:"origQuoteOrderQty,string"`
	// UpdateTime          int64   `json:"updateTime"`
}

// FutureNewOrderRequest request type
type FutureNewOrderRequest struct {
	// Symbol (currency pair to trade)
	Symbol string
	// Side Buy or Sell
	Side string
	// 持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
	PositionSide string
	// TradeType (market or limit order)
	TradeType RequestParamsOrderType
	// true, false; 非双开模式下默认false；双开模式下不接受此参数； 使用closePosition不支持此参数。
	ReduceOnly string
	// 下单数量,使用closePosition不支持此参数
	Quantity float64
	//委托价格
	Price float64
	//用户自定义的订单号，不可以重复出现在挂单中。如空缺系统会自动赋值。必须满足正则规则 ^[a-zA-Z0-9-_]{1,36}$
	NewClientOrderID string
	StopPrice        float64 // Used with STOP_LOSS, STOP_LOSS_LIMIT, TAKE_PROFIT, and TAKE_PROFIT_LIMIT orders.
	// true, false；触发后全部平仓，仅支持STOP_MARKET和TAKE_PROFIT_MARKET；不与quantity合用；自带只平仓效果，不与reduceOnly 合用
	ClosePosition float64
	// 追踪止损激活价格，仅TRAILING_STOP_MARKET 需要此参数, 默认为下单当前市场价格(支持不同workingType)
	ActivationPrice float64
	// 追踪止损回调比例，可取值范围[0.1, 5],其中 1代表1% ,仅TRAILING_STOP_MARKET 需要此参数
	CallbackRate float64
	// TimeInForce specifies how long the order remains in effect.
	// Examples are (Good Till Cancel (GTC), Immediate or Cancel (IOC) and Fill Or Kill (FOK))
	TimeInForce RequestParamsTimeForceType
	//  触发类型: MARK_PRICE(标记价格), CONTRACT_PRICE(合约最新价). 默认 CONTRACT_PRICE
	WorkingType string
	//  条件单触发保护："TRUE","FALSE", 默认"FALSE". 仅 STOP, STOP_MARKET, TAKE_PROFIT, TAKE_PROFIT_MARKET 需要此参数
	priceProtect     string
	NewOrderRespType string
}

// FutureNewOrderResponse is the return structured response from the exchange
type FutureNewOrderResponse struct {
	Code          int     `json:"code"`
	Msg           string  `json:"msg"`
	Symbol        string  `json:"symbol"` //交易对
	OrderID       int64   `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	AvgPrice      int64   `json:"avgPrice"`       //平均成交价
	Price         float64 `json:"price,string"`   //委托价格
	OrigQty       float64 `json:"origQty,string"` //原始委托数量
	CumQty        float64 `json:"cumQty,string"`
	CumQuote      float64 `json:"cumQuote,string"` //成交金额
	// The cumulative amount of the quote that has been spent (with a BUY order) or received (with a SELL order).
	ExecutedQty   float64 `json:"executedQty,string"` //成交量
	Status        string  `json:"status"`             //订单状态
	TimeInForce   string  `json:"timeInForce"`        //有效方法
	Type          string  `json:"type"`               //订单类型
	Side          string  `json:"side"`               //买卖方向
	PositionSide  string  `json:"positionSide"`       //持仓方向
	StopPrice     string  `json:"stopPrice"`          //触发价，对`TRAILING_STOP_MARKET`无效
	ClosePosition string  `json:"closePosition"`      //是否条件全平仓
	OrigType      string  `json:"origType"`           //触发前订单类型
	ActivatePrice string  `json:"activatePrice"`      //跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     string  `json:"priceRate"`          //跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段

	UpdateTime   string `json:"updateTime"`   // 更新时间
	WorkingType  string `json:"workingType"`  // 条件价格触发类型
	PriceProtect bool   `json:"priceProtect"` // 是否开启条件单触发保护
}

type FutureFundingRateResponeItem struct {
	Symbol      string    `json:"symbol"`
	FundingRate float64   `json:"fundingRate"`
	FundingTime time.Time `json:"fundingTime"`
}

// KlinesContractRequestParams represents Klines request data.
type KlinesContractRequestParams struct {
	Pair         string // Required field; example LTCBTC, BTCUSDT
	contractType asset.Item
	Interval     string // Time interval period
	Limit        int    // Default 500; max 500.
	StartTime    int64
	EndTime      int64
}

// PreminuIndexResponse represents Klines request data.
type PreminuIndexResponse struct {
	Synbol          string    `json:"symbol"`          // Required field; example LTCBTC, BTCUSDT
	MarkPrice       float64   `json:"markPrice"`       // 标记价格
	IndexPrice      float64   `json:"indexPrice"`      // 指数价格
	LastFundingRate float64   `json:"lastFundingRate"` // 最近更新的资金费率
	NextFundingTime time.Time `json:"nextFundingTime"` // 下次资金费时间
	InterestRate    float64   `json:"interestRate"`    // 标的资产基础利率
	Time            time.Time `json:"time"`            // 更新时间
}

// AccountSnapshotRequest 查询每日资产快照 (USER_DATA)
type AccountSnapshotRequest struct {
	Type      asset.Item `json:"type"`
	Price     float64    `json:"price"`
	Limit     int64      `json:"limit"`
	StartTime int64      `json:"startTime"`
	EndTime   int64      `json:"endTime"`
}

// AccountSnapshotResponse 查询每日资产快照 (USER_DATA) - 返回信息
type AccountSnapshotResponse struct {
	TotalAssetOfBtc float64    `json:"totalAssetOfBtc"`
	Asset           asset.Item `json:"asset"`
	Symbol          string     `json:"symbol"`
	Free            float64    `json:"free"`
	Locked          float64    `json:"locked"`
	UpdateTime      time.Time  `json:"updateTime"`
}
