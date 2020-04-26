package okgroup

import (
	"github.com/idoall/gocryptotrader/currency"
	"github.com/idoall/gocryptotrader/exchanges/asset"
	"time"
)

// WebsocketResponsePosition defines
type WebsocketResponsePosition struct {
	Timestamp    time.Time
	Pair         currency.Pair
	AssetType    asset.Item
	ExchangeName string
	Holding []WebsocketResponsePositionHoldingData `json:"holding,omitempty"`
}

// WebsocketResponsePositionHoldingData contains formatted data for user position holding data
type WebsocketResponsePositionHoldingData struct {
	InstrumentID string    `json:"instrument_id"`
	AvailablePosition float64   `json:"avail_position,string,omitempty"`
	AverageCost       float64   `json:"avg_cost,string,omitempty"`
	Leverage          float64   `json:"leverage,string,omitempty"`
	LiquidationPrice  float64   `json:"liquidation_price,string,omitempty"`
	Margin            float64   `json:"margin,string,omitempty"`
	Position          float64   `json:"position,string,omitempty"`
	RealizedPnl       float64   `json:"realized_pnl,string,omitempty"`
	SettlementPrice   float64   `json:"settlement_price,string,omitempty"`
	Side              string    `json:"side,omitempty"`
	Timestamp         time.Time `json:"timestamp,omitempty"`
}

// WebsocketResponseMarkPrice defines
type WebsocketResponseMarkPrice struct {
	Timestamp    time.Time
	Pair         currency.Pair
	AssetType    asset.Item
	ExchangeName string
	Price        float64
}


type WebsocketResponseOrders struct {
	Timestamp    time.Time
	Pair         currency.Pair
	AssetType    asset.Item
	ExchangeName string
	OrderInfo []WebsocketResponseOrdersData `json:"order_info"`
}

// GetSwapOrderListResponseData individual order data from GetSwapOrderList
type WebsocketResponseOrdersData struct {
	ContractVal  float64   `json:"contract_val,string"`
	Fee          float64   `json:"fee,string"`
	FilledQty    float64   `json:"filled_qty,string"`
	InstrumentID string    `json:"instrument_id"`
	Leverage     float64   `json:"leverage,string"` //  	Leverage value:10\20 default:10
	OrderID      int64     `json:"order_id,string"`
	Price        float64   `json:"price,string"`
	PriceAvg     float64   `json:"price_avg,string"`
	Size         float64   `json:"size,string"`
	Status       int64     `json:"status,string"` // Order Status （-1 canceled; 0: pending, 1: partially filled, 2: fully filled)
	Timestamp    time.Time `json:"timestamp"`
	Type         int64     `json:"type,string"` //  	Type (1: open long 2: open short 3: close long 4: close short)
}