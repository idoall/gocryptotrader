package huobi

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
