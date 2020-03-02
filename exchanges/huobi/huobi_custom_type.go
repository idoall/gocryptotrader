package huobi

//----------合约用用户帐号信息相关

// SymbolBaseType 基础信息
type SymbolBaseType struct {
	Symbol string `json:"symbol"`
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
	LeverRate         float64                   `json:"lever_rate"`         // 杠杆倍数
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
	LeverRate      int               `json:"lever_rate"`      // 杠杆倍数
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
	Data        []ContractOrdersItem `json:"orders"`
	TotalPage   int                  `json:"total_page"`   // 总页数
	CurrentPage int                  `json:"current_page"` // 当前页
	TotalSize   int                  `json:"total_size"`   // 总条数
}

// ContractOrdersItem 合约订单信息
type ContractOrdersItem struct {
	Symbol         string       `json:"symbol"`           //品种代码
	ContractType   ContractType `json:"contract_type"`    // 合约类型		当周:"this_week", 次周:"next_week", 季度:"quarter"
	Volume         float64      `json:"volume"`           //委托数量
	Price          float64      `json:"price"`            //委托价格
	OrderPriceType string       `json:"order_price_type"` //订单报价类型 1限价单，3对手价，4闪电平仓，5计划委托，6post_only
	Direction      string       `json:"direction"`        // 买卖方向
	Offset         string       `json:"offset"`           // 开平方向
	LeverRate      string       `json:"lever_rate"`       // 杠杆倍数
	OrderID        string       `json:"order_id"`         // 订单ID
	OrderIDStr     string       `json:"order_id_str"`     // String类型订单ID
	OrderSource    string       `json:"order_source"`     // 订单来源
	CreateDate     int64        `json:"create_date"`      // 创建时间
	TradeVolume    float64      `json:"trade_volume"`     // 成交数量
	TradeTurnover  float64      `json:"trade_turnover"`   // 成交总金额
	Fee            float64      `json:"fee"`              //手续费
	TradeAvgPrice  float64      `json:"trade_avg_price"`  // 成交均价
	MarginFrozen   float64      `json:"margin_frozen"`    // 冻结保证金
	Profit         float64      `json:"profit"`           // 收益
	Status         int          `json:"status"`           // 订单状态
	OrderType      int          `json:"order_type"`       // 订单类型
	FeeAsset       string       `json:"fee_asset"`        // 手续费币种
}

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
