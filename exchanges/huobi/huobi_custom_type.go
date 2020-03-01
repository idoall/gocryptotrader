package huobi

//----------合约用用户帐号信息相关
// ContractAccountInfoRequest stores a kline item
type ContractAccountInfoRequest struct {
	Symbol string `json:"symbol"`
}

// ContractAccountInfoResponseDataItem 帐号信息返回参数
type ContractAccountInfoResponseDataItem struct {
	Symbol            string  `json:"symbol"`
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
