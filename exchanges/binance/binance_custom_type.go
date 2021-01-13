package binance

import (
	"time"

	"github.com/idoall/gocryptotrader/currency"
	"github.com/idoall/gocryptotrader/exchanges/asset"
	"github.com/idoall/gocryptotrader/exchanges/order"
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
	// 撤销订单 (TRADE)
	binanceFutureCancelOrder = "/fapi/v1/order"
	// 查看当前全部挂单
	binanceFutureOpenOrders = "/fapi/v1/openOrders"
	// 调整开仓杠杆
	binanceFutureLeverage = "/fapi/v1/leverage"
	// 获取账户损益资金流水(USER_DATA)
	binanceFutureIncome = "/fapi/v1/income"
	// 用户持仓风险V2 (USER_DATA)
	binanceFuturePositionRisk = "/fapi/v2/positionRisk"
	// 变换逐全仓模式 (USER_DATA)
	binanceFutureMarginType = "/fapi/v1/marginType"

	// 用户万向划转
	binanceTransfer = "/sapi/v1/asset/transfer"

	// 交易手续费率查询
	binanceSpotTradeFee   = "/wapi/v3/tradeFee.html"
	binanceFutureTradeFee = "/fapi/v1/commissionRate"
)

// CommissionRateResponse 交易手续费率
type CommissionRateResponse struct {
	Symbol string
	Maker  float64
	Taker  float64
}

// PositionRiskResponse 用户持仓风险
type PositionRiskResponse struct {
	EntryPrice       float64      `json:"entryPrice,string"` //开仓均价
	MarginType       MarginType   `json:"marginType"`        //逐仓模式或全仓模式
	IsAutoAddMargin  bool         `json:"isAutoAddMargin, string"`
	IsolatedMargin   float64      `json:"isolatedMargin, string"`   //  逐仓保证金
	Leverage         int64        `json:"leverage, string"`         // 当前杠杆倍数
	LiquidationPrice float64      `json:"liquidationPrice, string"` // 参考强平价格
	MarkPrice        float64      `json:"markPrice, string"`        // 当前标记价格
	MaxNotionalValue int64        `json:"maxNotionalValue, string"` // 当前杠杆倍数允许的名义价值上限
	PositionAmt      float64      `json:"positionAmt, string"`      // 头寸数量，符号代表多空方向, 正数为多，负数为空
	Symbol           string       `json:"symbol"`                   // 交易对
	UnRealizedProfit float64      `json:"unRealizedProfit, string"` // 持仓未实现盈亏
	PositionSide     PositionSide `json:"positionSide"`             //  持仓方向
}

// MarginType 保证金模式
type MarginType string

const (
	// MarginType_ISOLATED 逐仓
	MarginType_ISOLATED = MarginType("ISOLATED")
	// MarginType_CROSSED 全仓
	MarginType_CROSSED = MarginType("CROSSED")
)

type FutureFundingRateRequest struct {
	Symbol    currency.Pair `json:"symbol"` //交易对
	StartTime int64         `json:"startTime"`
	EndTime   int64         `json:"endTime"`
	Limit     int64         `json:"limit"`
}

type FutureFundingRateResponeItem struct {
	Symbol      string    `json:"symbol"`
	FundingRate float64   `json:"fundingRate"`
	FundingTime time.Time `json:"fundingTime"`
}

type FutureIncomeResponse struct {
	Symbol     string     `json:"symbol"`         //交易对
	Income     float64    `json:"income, string"` //资金流数量，正数代表流入，负数代表流出
	IncomeType IncomeType `json:"incomeType"`     // 收益类型
	Asset      string     `json:"asset"`
	Info       string     `json:"info,string"`
	Time       time.Time  `json:"time"`
	TranId     int64      `json:"tranId,string"`
	TradeId    string     `json:"tradeId"`
}

type FutureIncomeRequest struct {
	Symbol     string     `json:"symbol"`     //交易对
	IncomeType IncomeType `json:"incomeType"` // 收益类型
	StartTime  int64      `json:"startTime"`
	EndTime    int64      `json:"endTime"`
	Limit      int64      `json:"limit"`
}

// WorkingType 条件价格触发类型 (workingType)
type WorkingType string

const (
	WorkingType_MARK_PRICE     = WorkingType("MARK_PRICE")
	WorkingType_CONTRACT_PRICE = WorkingType("CONTRACT_PRICE")
)

// IncomeType收益类型
type IncomeType string

const (
	IncomeType_TRANSFER        = IncomeType("TRANSFER")
	IncomeType_WELCOME_BONUS   = IncomeType("WELCOME_BONUS")
	IncomeType_REALIZED_PNL    = IncomeType("REALIZED_PNL")
	IncomeType_FUNDING_FEE     = IncomeType("FUNDING_FEE")
	IncomeType_COMMISSION      = IncomeType("COMMISSION")
	IncomeType_INSURANCE_CLEAR = IncomeType("INSURANCE_CLEAR")
	IncomeType_ALL             = IncomeType("")
)

// TransferType 用户万向划转 类型
type TransferType string

const (
	//MAIN_C2C 现货钱包转向C2C钱包
	TransferType_MAIN_C2C = TransferType("MAIN_C2C")
	//MAIN_UMFUTURE 现货钱包转向U本位合约钱包
	TransferType_MAIN_UMFUTURE = TransferType("MAIN_UMFUTURE")
	//MAIN_CMFUTURE 现货钱包转向币本位合约钱包
	TransferType_MAIN_CMFUTURE = TransferType("MAIN_CMFUTURE")
	//MAIN_MARGIN 现货钱包转向杠杆全仓钱包
	TransferType_MAIN_MARGIN = TransferType("MAIN_MARGIN")
	//MAIN_MINING 现货钱包转向矿池钱包
	TransferType_MAIN_MINING = TransferType("MAIN_MINING")
	//C2C_MAIN C2C钱包转向现货钱包
	TransferType_C2C_MAIN = TransferType("C2C_MAIN")
	//C2C_UMFUTURE C2C钱包转向U本位合约钱包
	TransferType_C2C_UMFUTURE = TransferType("C2C_UMFUTURE")
	//C2C_MINING C2C钱包转向矿池钱包
	TransferType_C2C_MINING = TransferType("C2C_MINING")
	//UMFUTURE_MAIN U本位合约钱包转向现货钱包
	TransferType_UMFUTURE_MAIN = TransferType("UMFUTURE_MAIN")
	//UMFUTURE_C2C U本位合约钱包转向C2C钱包
	TransferType_UMFUTURE_C2C = TransferType("UMFUTURE_C2C")
	//UMFUTURE_MARGIN U本位合约钱包转向杠杆全仓钱包
	TransferType_UMFUTURE_MARGIN = TransferType("UMFUTURE_MARGIN")
	//CMFUTURE_MAIN 币本位合约钱包转向现货钱包
	TransferType_CMFUTURE_MAIN = TransferType("CMFUTURE_MAIN")
	//MARGIN_MAIN 杠杆全仓钱包转向现货钱包
	TransferType_MARGIN_MAIN = TransferType("MARGIN_MAIN")
	//MARGIN_UMFUTURE 杠杆全仓钱包转向U本位合约钱包
	TransferType_MARGIN_UMFUTURE = TransferType("MARGIN_UMFUTURE")
	//MINING_MAIN 矿池钱包转向现货钱包
	TransferType_MINING_MAIN = TransferType("MINING_MAIN")
	//TransferType_MINING_UMFUTURE MINING_UMFUTURE 矿池钱包转向U本位合约钱包
	TransferType_MINING_UMFUTURE = TransferType("MINING_UMFUTURE")
	// TransferType_MINING_C2CMINING_C2C 矿池钱包转向C2C钱包
	TransferType_MINING_C2C = TransferType("MINING_C2C")
)

type FutureLeverageResponse struct {
	Leverage         int    `json:"leverage,string"`          // 平均成交价
	MaxNotionalValue int64  `json:"maxNotionalValue, string"` // 用户自定义的订单号
	Symbol           string `json:"symbol"`                   //成交金额
}

// FutureQueryOrderData holds query order data
type FutureQueryOrderData struct {
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
	Time          time.Time                  `json:"time"`                 // 订单时间
	TimeInForce   RequestParamsTimeForceType `json:"timeInForce"`          // 有效方法
	Type          string                     `json:"type"`                 //订单类型
	ActivatePrice float64                    `json:"activatePrice,string"` // 跟踪止损激活价格, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	PriceRate     float64                    `json:"priceRate,string"`     // 跟踪止损回调比例, 仅`TRAILING_STOP_MARKET` 订单返回此字段
	UpdateTime    time.Time                  `json:"updateTime"`
	WorkingType   WorkingType                `json:"workingType"`  // 条件价格触发类型
	PriceProtect  bool                       `json:"priceProtect"` // 是否开启条件单触发保护
}

// PositionSide 持仓方向
type PositionSide string

const (
	// PositionSideBOTH 单一持仓方向
	PositionSideBOTH = PositionSide("BOTH")
	// PositionSideLONG 多头(双向持仓下)
	PositionSideLONG = PositionSide("LONG")
	// PositionSideSHORT 空头(双向持仓下)
	PositionSideSHORT = PositionSide("SHORT")
)

// ContractType 合约类型
type ContractType string

const (
	// ContractTypePERPETUAL 永续合约
	ContractTypePERPETUAL = ContractType("PERPETUAL")
	// ContractTypeCURRENT_MONTH 当月交割合约
	ContractTypeCURRENT_MONTH = ContractType("CURRENT_MONTH")
	// PositionSideSHORT 次月交割合约
	ContractTypeNEXT_MONTH = ContractType("NEXT_MONTH")
)

// FutureNewOrderRequest request type
type FutureNewOrderRequest struct {
	// Symbol (currency pair to trade)
	Symbol string
	// Side Buy or Sell
	Side order.Side
	// 持仓方向，单向持仓模式下非必填，默认且仅可填BOTH;在双向持仓模式下必填,且仅可选择 LONG 或 SHORT
	PositionSide PositionSide
	// Type 订单类型 LIMIT, MARKET, STOP, TAKE_PROFIT, STOP_MARKET, TAKE_PROFIT_MARKET, TRAILING_STOP_MARKET
	Type RequestParamsOrderType
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
	Symbol        string  `json:"symbol"` //交易对
	OrderID       int64   `json:"orderId"`
	ClientOrderID string  `json:"clientOrderId"`
	AvgPrice      float64 `json:"avgPrice, string"` //平均成交价
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

	UpdateTime   time.Time   `json:"updateTime"`   // 更新时间
	WorkingType  WorkingType `json:"workingType"`  // 条件价格触发类型
	PriceProtect bool        `json:"priceProtect"` // 是否开启条件单触发保护
}

// KlinesContractRequestParams represents Klines request data.
type KlinesContractRequestParams struct {
	Pair         string // Required field; example LTCBTC, BTCUSDT
	contractType ContractType
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
