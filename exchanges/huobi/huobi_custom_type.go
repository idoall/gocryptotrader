package huobi

import (
	"time"
)

//----------合约用用户帐号信息相关

// SymbolBaseType 基础信息
type SymbolBaseType struct {
	Symbol string `json:"symbol"`
}

// ContractBaseType 基础信息
type ContractBaseType struct {
	ContractType ContractType `json:"contract_type"` // 合约类型 当周:"this_week", 次周:"next_week", 季度:"quarter"
	ContractCode string       `json:"contract_code"` // 合约代码 "BTC180914" ...
}

// ContractAccountInfoRequest stores a kline item
type ContractAccountInfoRequest struct {
	SymbolBaseType
}

// ContractAccountInfoResponseDataItem 帐号信息返回参数
type ContractAccountInfoResponseDataItem struct {
	SymbolBaseType
	MarginBalance     float64 `json:"margin_balance"`     //账户权益
	MarginPosition    float64 `json:"margin_position"`    //持仓保证金（当前持有仓位所占用的保证金）
	MarginFrozen      float64 `json:"margin_frozen"`      //冻结保证金
	MarginAvailable   float64 `json:"margin_available"`   //可用保证金
	ProfitReal        float64 `json:"profit_real"`        //已实现盈亏
	ProfitUnreal      float64 `json:"profit_unreal"`      //未实现盈亏
	WithdrawAvailable float64 `json:"withdraw_available"` //可划转数量
	RiskRate          float64 `json:"risk_rate"`          //保证金率
	LiquidationPrice  float64 `json:"liquidation_price"`  //预估强平价
	AdjustFactor      float64 `json:"adjust_factor"`      //调整系数
	LeverRate         float64 `json:"lever_rate"`         //杠杠倍数
	MarginStatic      float64 `json:"margin_static"`      //静态权益
}

//-----------合约新订单

// ContractNewOrderRequest 新订单 - 请求信息
type ContractNewOrderRequest struct {
	SymbolBaseType
	ContractBaseType
	ClientOrderID  int64                  `json:"client_order_id"`  //客户订单ID
	Price          float64                `json:"price"`            // 价格
	Volume         int64                  `json:"volume"`           //委托数量(张)
	LeverRae       LeverRae               `json:"lever_rate"`       // 杠杆倍数[“开仓”若有10倍多单，就不能再下20倍多单]
	Direction      ContractOrderDirection `json:"direction"`        //交易方向	"buy":买 "sell":卖
	Offset         ContractOrderOffset    `json:"offset"`           //"open":开 "close":平
	OrderPriceType ContractOrderPriceType `json:"order_price_type"` //订单报价类型 "limit":限价 "opponent":对手价 "post_only":只做maker单,post only下单只受用户持仓数量限制,optimal_5：最优5档、optimal_10：最优10档、optimal_20：最优20档，ioc:IOC订单，fok：FOK订单, "opponent_ioc"： 对手价-IOC下单，"optimal_5_ioc"：最优5档-IOC下单，"optimal_10_ioc"：最优10档-IOC下单，"optimal_20_ioc"：最优20档-IOC下单,"opponent_fok"： 对手价-FOK下单，"optimal_5_fok"：最优5档-FOK下单，"optimal_10_fok"：最优10档-FOK下单，"optimal_20_fok"：最优20档-FOK下单
}

// ContractNewOrderResponse 新订单请求信息 - 返回信息
type ContractNewOrderResponse struct {
	Response
	OrderID       int64  `json:"order_id"`        //订单ID
	OrderIDStr    string `json:"order_id_str"`    //String订单ID
	ClientOrderID int64  `json:"client_order_id"` //客户订单ID
}

// ContractNewTriggerOrderRequest 合约计划委托下单 - 请求信息
type ContractNewTriggerOrderRequest struct {
	SymbolBaseType
	ContractBaseType
	TriggerType    ContractTriggerType    `json:"trigger_type"`     //触发类型： ge大于等于(触发价比最新价大)；le小于(触发价比最新价小)
	TriggerPrice   float64                `json:"trigger_price"`    // 触发价，精度超过最小变动单位会报错
	OrderPrice     float64                `json:"order_price"`      //委托价，精度超过最小变动单位会报错
	OrderPriceType ContractOrderPriceType `json:"order_price_type"` // 委托类型： 不填默认为limit; 限价：limit ，最优5档：optimal_5，最优10档：optimal_10，最优20档：optimal_20
	Direction      ContractOrderDirection `json:"direction"`        //交易方向	"buy":买 "sell":卖
	Offset         ContractOrderOffset    `json:"offset"`           //"open":开 "close":平
	Volume         int                    `json:"volume"`           //委托数量(张)
	LeverRae       LeverRae               `json:"lever_rate"`       // 杠杆倍数[开仓若有10倍多单，就不能再下20倍多单]
}

// ContractTriggerType 合约下单 触发类型
type ContractTriggerType string

var (
	// ContractTriggerTypeGE ge大于等于(触发价比最新价大)
	ContractTriggerTypeGE = ContractTriggerType("ge")
	// ContractTriggerTypeLE 小于(触发价比最新价小)
	ContractTriggerTypeLE = ContractTriggerType("le")
)

// ContractNewTriggerOrderResponse 合约计划委托下单 - 返回信息
type ContractNewTriggerOrderResponse struct {
	Response
	Data struct {
		OrderID    int64  `json:"order_id"`     //订单ID
		OrderIDStr string `json:"order_id_str"` //String订单ID
	} `json:"data"` //客户订单ID
}

//-----------合约订单信息

// ContractOpenOrderData 合约订单信息
type ContractOpenOrderData struct {
	Orders []ContractOrderDataItem `json:"orders"`
	ContractOrderPages
}

// ContractOrderDataItem 合约订单信息
type ContractOrderDataItem struct {
	SymbolBaseType
	ContractType   ContractType               `json:"contract_type"`    // 合约类型 当周:"this_week", 次周:"next_week", 季度:"quarter"
	ContractCode   string                     `json:"contract_code"`    // 合约代码 "BTC180914" ...
	Volume         float64                    `json:"volume"`           //委托数量
	Price          float64                    `json:"price"`            //委托价格
	OrderPriceType ContractOpenOrderPriceType `json:"order_price_type"` //订单报价类型
	OrderType      ContractOpenOrderType      `json:"order_type"`       //订单类型
	Direction      ContractOrderDirection     `json:"direction"`        //交易方向	"buy":买 "sell":卖
	Offset         ContractOrderOffset        `json:"offset"`           //"open":开 "close":平
	LeverRate      LeverRae                   `json:"lever_rate"`       //杠杆倍数
	OrderID        int64                      `json:"order_id"`         //订单ID
	OrderIDStr     string                     `json:"order_id_str"`     //String订单ID
	ClientOrderID  int64                      `json:"client_order_id"`  //客户订单ID
	OrderSource    string                     `json:"order_source"`     //订单来源
	CreateDate     int64                      `json:"created_at"`       //订单创建时间
	TradeVolume    float64                    `json:"trade_volume"`     //成交数量
	TradeTurnover  float64                    `json:"trade_turnover"`   //成交总金额
	Fee            float64                    `json:"fee"`              //手续费
	TradeAvgPrice  float64                    `json:"trade_avg_price"`  //成交均价
	MarginFrozen   float64                    `json:"margin_frozen"`    //冻结保证金
	Profit         float64                    `json:"profit"`           //收益
	Status         ContractOpenOrderStauts    `json:"status"`           //订单状态
	FeeAsset       string                     `json:"fee_asset"`        //手续费币种
}

// ContractMatchResultsRequest 历史成交记录 - 请求信息
type ContractMatchResultsRequest struct {
	SymbolBaseType
	TradeType    TradeType `json:"trade_type"`
	CreateDate   int64     `json:"create_date"`   //可随意输入正整数，如果参数超过90则默认查询90天的数据
	ContractCode string    `json:"contract_code"` // 合约代码 "BTC180914" ...
	PageIndex    int       `json:"page_index"`    // 页码，不填默认第1页
	PageSize     int       `json:"page_size"`     // 不填默认20，不得多于50
}

// ContractMatchResultData 历史成交记录
type ContractMatchResultData struct {
	ContractOrderPages
	Trades []ContractMatchResultDataItem `json:"trades"`
}

// ContractMatchResultDataItem 历史成交记录
type ContractMatchResultDataItem struct {
	ID               string                 `json:"id"`                // 全局唯一的交易标识
	MatchID          int64                  `json:"match_id"`          // 撮合结果id, 与订单ws推送orders.$symbol以及撮合订单ws推送matchOrders.$symbol推送结果中的trade_id是相同的，非唯一，可重复，注意：一个撮合结果代表一个taker单和N个maker单的成交记录的集合，如果一个taker单吃了N个maker单，那这N笔trade都是一样的撮合结果id
	ContractType     ContractType           `json:"contract_type"`     // 合约类型 当周:"this_week", 次周:"next_week", 季度:"quarter"
	ContractCode     string                 `json:"contract_code"`     // 合约代码 "BTC180914" ...
	CreateDate       int64                  `json:"create_date"`       //订单创建时间
	Direction        ContractOrderDirection `json:"direction"`         //交易方向	"buy":买 "sell":卖
	Offset           ContractOrderOffset    `json:"offset"`            //"open":开 "close":平
	OffsetProfitloss float64                `json:"offset_profitloss"` //平仓盈亏
	OrderID          int64                  `json:"order_id"`          //订单ID
	OrderIDStr       string                 `json:"order_id_str"`      //String订单ID
	SymbolBaseType
	OrderSource   string  `json:"order_source"`   //订单来源
	TradeFee      float64 `json:"trade_fee"`      //成交手续费
	TradePrice    float64 `json:"trade_price"`    //成交价格
	TradeTurnover float64 `json:"trade_turnover"` //本笔成交金额
	TradeVolume   float64 `json:"trade_volume"`   //成交数量
	Role          string  `json:"role"`           //taker或maker
	FeeAsset      string  `json:"fee_asset"`      //手续费币种
}

// ContractOpenOrderStauts 订单状态
type ContractOpenOrderStauts int

var (
	// ContractOpenOrderStauts3 未成交
	ContractOpenOrderStauts3 = ContractOpenOrderStauts(3)
	// ContractOpenOrderStauts4 部分成交
	ContractOpenOrderStauts4 = ContractOpenOrderStauts(4)
	// ContractOpenOrderStauts5 部分成交已撤单
	ContractOpenOrderStauts5 = ContractOpenOrderStauts(5)
	// ContractOpenOrderStauts6 全部成交
	ContractOpenOrderStauts6 = ContractOpenOrderStauts(6)
	// ContractOpenOrderStauts7 已撤单
	ContractOpenOrderStauts7 = ContractOpenOrderStauts(7)
)

// ContractOrderOffset 开仓或平仓
type ContractOrderOffset string

var (
	// ContractOrderOffsetOpen "open":开
	ContractOrderOffsetOpen = ContractOrderOffset("open")
	// ContractOrderOffsetClose "close":平
	ContractOrderOffsetClose = ContractOrderOffset("close")
)

// ContractOrderDirection 订单交易方向
type ContractOrderDirection string

var (
	// ContractOrderDirectionBuy 买
	ContractOrderDirectionBuy = ContractOrderDirection("buy")
	// ContractOrderDirectionSell 卖
	ContractOrderDirectionSell = ContractOrderDirection("sell")
)

// ContractOpenOrderPriceType 订单报价类型
type ContractOpenOrderPriceType string

var (
	// ContractOpenOrderPriceTypeLimit 限价
	ContractOpenOrderPriceTypeLimit = ContractOpenOrderPriceType("limit")
	// ContractOpenOrderPriceTypeOpponent 对手价
	ContractOpenOrderPriceTypeOpponent = ContractOpenOrderPriceType("opponent")
	// ContractOpenOrderPriceTypePostOnly 只做maker单,post only下单只受用户持仓数量限制
	ContractOpenOrderPriceTypePostOnly = ContractOpenOrderPriceType("post_only")
)

// ContractOpenOrderType 订单类型
type ContractOpenOrderType string

var (
	// ContractOpenOrderType1   1:报单
	ContractOpenOrderType1 = ContractOpenOrderType("1") //
	// ContractOpenOrderType2	2:撤单
	ContractOpenOrderType2 = ContractOpenOrderType("2") //
	// ContractOpenOrderType3 3:强平
	ContractOpenOrderType3 = ContractOpenOrderType("3") //
	// ContractOpenOrderType4 3:强平
	ContractOpenOrderType4 = ContractOpenOrderType("4") //
)

//--------帐号和持仓信息

// ContractAccountPositionInfoResponse 用户帐户信息
type ContractAccountPositionInfoResponse struct {
	SymbolBaseType
	MarginBalance     float64                   `json:"margin_balance"`     // 账户权益
	MarginPosition    float64                   `json:"margin_position"`    // 持仓保证金
	MarginFrozen      float64                   `json:"margin_frozen"`      // 冻结保证金
	MarginAvailable   float64                   `json:"margin_available"`   // 可用保证金
	ProfitReal        float64                   `json:"profit_real"`        // 已实现盈亏
	ProfitUnReal      float64                   `json:"profit_unreal"`      // 未实现盈亏
	RiskRate          float64                   `json:"risk_rate"`          // 保证金率
	WithdrawAvailable float64                   `json:"withdraw_available"` // 可划转数量
	LiquidationPrice  float64                   `json:"liquidation_price"`  // 预估爆仓价
	LeverRate         LeverRae                  `json:"lever_rate"`         // 杠杆倍数
	AdjustFactor      float64                   `json:"adjust_factor"`      // 调整系数
	MarginStatic      float64                   `json:"margin_static"`      // 静态权益
	Positions         []ContractAccountPosition `json:"positions"`          // 持仓信息
}

// ContractAccountPosition 用户持仓信息
type ContractAccountPosition struct {
	SymbolBaseType
	ContractType   ContractType      `json:"contract_type"`   // 合约类型 当周:"this_week", 次周:"next_week", 季度:"quarter"
	ContractCode   string            `json:"contract_code"`   // 合约代码 "BTC180914" ...
	Volume         float64           `json:"volume"`          // 持仓量
	Available      float64           `json:"available"`       // 可平仓数量
	Frozen         float64           `json:"frozen"`          // 冻结数量
	CostOpen       float64           `json:"cost_open"`       // 开仓均价
	CostHold       float64           `json:"cost_hold"`       // 持仓均价
	ProfitUnreal   float64           `json:"profit_unreal"`   // 未实现盈亏
	PofitRate      float64           `json:"profit_rate"`     // 收益率
	Pofit          float64           `json:"profit"`          // 收益
	PositionMargin float64           `json:"position_margin"` // 持仓保证金
	LeverRate      LeverRae          `json:"lever_rate"`      // 杠杆倍数
	Direction      ContractDirection `json:"direction"`       // "buy":买 "sell":卖
	LastPrice      float64           `json:"last_price"`      // 最新价
}

// ContractDirection 交易方向
type ContractDirection string

var (
	// ContractDirectionBuy 买入开多
	ContractDirectionBuy = ContractDirection("buy")
	// ContractDirectionSell 卖出开空
	ContractDirectionSell = ContractDirection("sell")
)

//--------历史订单相关

// ContractHisordersData 历史订单信息
type ContractHisordersData struct {
	Data []ContractOrderDataItem `json:"orders"`
	ContractOrderPages
}

// ContractOrderPages 分页相关
type ContractOrderPages struct {
	TotalPage   int `json:"total_page"`   // 总页数
	CurrentPage int `json:"current_page"` // 当前页
	TotalSize   int `json:"total_size"`   // 总条数

}

// ContractOrdersItem 合约订单信息
// type ContractOrdersItem struct {
// 	Symbol         string       `json:"symbol"`           //品种代码
// 	ContractType   ContractType `json:"contract_type"`    // 合约类型		当周:"this_week", 次周:"next_week", 季度:"quarter"
// 	Volume         float64      `json:"volume"`           //委托数量
// 	Price          float64      `json:"price"`            //委托价格
// 	OrderPriceType string       `json:"order_price_type"` //订单报价类型 1限价单，3对手价，4闪电平仓，5计划委托，6post_only
// 	Direction      string       `json:"direction"`        // 买卖方向
// 	Offset         string       `json:"offset"`           // 开平方向
// 	LeverRate      string       `json:"lever_rate"`       // 杠杆倍数
// 	OrderID        string       `json:"order_id"`         // 订单ID
// 	OrderIDStr     string       `json:"order_id_str"`     // String类型订单ID
// 	OrderSource    string       `json:"order_source"`     // 订单来源
// 	CreateDate     int64        `json:"create_date"`      // 创建时间
// 	TradeVolume    float64      `json:"trade_volume"`     // 成交数量
// 	TradeTurnover  float64      `json:"trade_turnover"`   // 成交总金额
// 	Fee            float64      `json:"fee"`              //手续费
// 	TradeAvgPrice  float64      `json:"trade_avg_price"`  // 成交均价
// 	MarginFrozen   float64      `json:"margin_frozen"`    // 冻结保证金
// 	Profit         float64      `json:"profit"`           // 收益
// 	Status         int          `json:"status"`           // 订单状态
// 	OrderType      int          `json:"order_type"`       // 订单类型
// 	FeeAsset       string       `json:"fee_asset"`        // 手续费币种
// }

// ContractType 合约类型
type ContractType string

var (
	// ContractTypeThisWeek i当周
	ContractTypeThisWeek = ContractType("this_week")
	// ContractTypeNextWeek 次周
	ContractTypeNextWeek = ContractType("next_week")
	// ContractTypeQuarter 季度
	ContractTypeQuarter = ContractType("quarter")
)

// ContractHisordersRequest 合约历史委托请求参数
type ContractHisordersRequest struct {
	Symbol       string               `json:"symbol"`
	TradeType    TradeType            `json:"trade_type"` //交易类型
	Type         ContractHisOrderType `json:"type"`
	Status       ContractOrderStatus  `json:"status"`
	CreateDate   int                  `json:"create_date"` //可随意输入正整数, ，如果参数超过90则默认查询90天的数据
	PageIndex    int                  `json:"page_index"`
	PageSize     int                  `json:"page_size"`
	ContractCode string               `json:"contract_code"` // 合约代码
	OrderType    ContractOrderType    `json:"order_type"`
}

// LeverRae 杠杆倍数
type LeverRae int

var (
	// LeverRae1 1
	LeverRae1 = LeverRae(1)
	// LeverRae5 5
	LeverRae5 = LeverRae(5)
	// LeverRae10 10
	LeverRae10 = LeverRae(10)
	// LeverRae20 20
	LeverRae20 = LeverRae(20)
)

// ContractOrderPriceType 订单报价类型
type ContractOrderPriceType string

var (
	// ContractOrderPriceTypeLimit 限价
	ContractOrderPriceTypeLimit = ContractOrderPriceType("limit")
	// ContractOrderPriceTypeOpponent 对手价
	ContractOrderPriceTypeOpponent = ContractOrderPriceType("opponent")
	// ContractOrderPriceTypePostOnly 只做maker单,post only下单只受用户持仓数量限制 默认是“只做Maker（Post only）”，不会立刻在市场成交，保证用户始终为Maker；如果委托会立即与已有委托成交，那么该委托会被取消。
	ContractOrderPriceTypePostOnly = ContractOrderPriceType("post_only")
	// ContractOrderPriceTypeOptimal5 最优5档
	ContractOrderPriceTypeOptimal5 = ContractOrderPriceType("optimal_5")
	// ContractOrderPriceTypeOptimal10 最优10档
	ContractOrderPriceTypeOptimal10 = ContractOrderPriceType("optimal_10")
	// ContractOrderPriceTypeOptimal20 最优20档
	ContractOrderPriceTypeOptimal20 = ContractOrderPriceType("optimal_20")
	// ContractOrderPriceTypeIOC IOC下单 市价单会在一个最佳可成交价执行尽量多的交易量,此单据可能被部份执行,剩余的部份将会自动删除。
	ContractOrderPriceTypeIOC = ContractOrderPriceType("ioc")
	// ContractOrderPriceTypeFOK FOK订单 市价单要么在一个最佳可成交价上全部成交，要么就会直接删除，即不会分多次成交，也不会部份成交。
	ContractOrderPriceTypeFOK = ContractOrderPriceType("fok")
	// ContractOrderPriceTypeOpponentIOC  对手价-IOC下
	ContractOrderPriceTypeOpponentIOC = ContractOrderPriceType("opponent_ioc")
	// ContractOrderPriceTypeOptimal5IOC 最优5档-IOC下单
	ContractOrderPriceTypeOptimal5IOC = ContractOrderPriceType("optimal_5_ioc")
	// ContractOrderPriceTypeOptimal10IOC 最优10档-IOC下单
	ContractOrderPriceTypeOptimal10IOC = ContractOrderPriceType("optimal_10_ioc")
	// ContractOrderPriceTypeOptimal20IOC 最优20档-IOC下单
	ContractOrderPriceTypeOptimal20IOC = ContractOrderPriceType("optimal_20_ioc")
	// ContractOrderPriceTypeOpponentFOK 对手价-FOK下单
	ContractOrderPriceTypeOpponentFOK = ContractOrderPriceType("opponent_fok")
	// ContractOrderPriceTypeOptimal5FOK 最优5档-FOK下单
	ContractOrderPriceTypeOptimal5FOK = ContractOrderPriceType("optimal_5_fok")
	// ContractOrderPriceTypeOptimal10FOK 最优10档-FOK下单
	ContractOrderPriceTypeOptimal10FOK = ContractOrderPriceType("optimal_10_fok")
	// ContractOrderPriceTypeOptimal20FOK 最优20档-FOK下单
	ContractOrderPriceTypeOptimal20FOK = ContractOrderPriceType("optimal_20_fok")
)

// ContractOrderType 订单类型
type ContractOrderType int

var (
	// ContractOrderType1 限价单
	ContractOrderType1 = ContractOrderType(1)
	// ContractOrderType2 对手价
	ContractOrderType2 = ContractOrderType(2)
	// ContractOrderType3 对手价
	ContractOrderType3 = ContractOrderType(3)
	// ContractOrderType4 闪电平仓
	ContractOrderType4 = ContractOrderType(4)
	// ContractOrderType5 计划委托
	ContractOrderType5 = ContractOrderType(5)
	// ContractOrderType6 post_only
	ContractOrderType6 = ContractOrderType(6)
	// ContractOrderType7 最优5档
	ContractOrderType7 = ContractOrderType(7)
	// ContractOrderType8 最优10档
	ContractOrderType8 = ContractOrderType(8)
	// ContractOrderType9 最优20档
	ContractOrderType9 = ContractOrderType(9)
	// ContractOrderType10 fok
	ContractOrderType10 = ContractOrderType(10)
	// ContractOrderType11 ioc
	ContractOrderType11 = ContractOrderType(11)
)

// ContractOrderStatus 订单状态
type ContractOrderStatus int

var (
	// ContractOrderStatus0 全部
	ContractOrderStatus0 = ContractOrderStatus(0)
	// ContractOrderStatus3 未成交
	ContractOrderStatus3 = ContractOrderStatus(3)
	// ContractOrderStatus4 部分成交
	ContractOrderStatus4 = ContractOrderStatus(4)
	// ContractOrderStatus5 部分成交已撤单
	ContractOrderStatus5 = ContractOrderStatus(5)
	// ContractOrderStatus6 全部成交
	ContractOrderStatus6 = ContractOrderStatus(6)
	// ContractOrderStatus7 已撤单
	ContractOrderStatus7 = ContractOrderStatus(7)
)

// ContractHisOrderType 历史订单查询类型
type ContractHisOrderType int

var (
	// ContractHisOrderType1 所有订单
	ContractHisOrderType1 = ContractHisOrderType(1)
	// ContractHisOrderType2 结束状态的订单
	ContractHisOrderType2 = ContractHisOrderType(2)
)

// TradeType 交易类型
type TradeType int

var (
	// TradeType0 : 全部
	TradeType0 = TradeType(0)
	// TradeType1 : 买入开多
	TradeType1 = TradeType(1)
	// TradeType2 : 卖出开空,
	TradeType2 = TradeType(2)
	// TradeType3 : 买入平空
	TradeType3 = TradeType(3)
	// TradeType4 : 卖出平多,
	TradeType4 = TradeType(4)
	// TradeType5 : 卖出强平
	TradeType5 = TradeType(5)
	// TradeType6 : 买入强平
	TradeType6 = TradeType(6)
	// TradeType7 :交割平多
	TradeType7 = TradeType(7)
	// TradeType8 : 交割平空
	TradeType8 = TradeType(8)
)

//----------合约信息相关

// ContractInfoRequest stores a kline item
type ContractInfoRequest struct {
	Symbol       string `json:"symbol"`
	ContractType string `json:"contract_type"`
	ContractCode string `json:"contract_code"`
}

// ContractInfoResponse stores a kline item
type ContractInfoResponse struct {
	Data         []ContractInfoResponseDataItem `json:"data"`
	Status       string                         `json:"status"`
	ErrorMessage string                         `json:"err-msg"`
}

// ContractInfoResponseDataItem
type ContractInfoResponseDataItem struct {
	Symbol         string  `json:"symbol"`
	ContractType   string  `json:"contract_type"`
	ContractCode   string  `json:"contract_code"`
	ContractSize   float64 `json:"contract_size"`
	PriceTick      float64 `json:"price_tick"`
	DeliveryDate   string  `json:"delivery_date"`
	CreateDate     string  `json:"create_date"`
	ContractStatus int     `json:"contract_status"`
}

// AccountAssetValuationResponse 获取账户资产估值返回数据.
type AccountAssetValuationResponse struct {
	Balance float64   `json:"balance"`
	Date    time.Time `json:"date"`
}
