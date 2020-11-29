package kline

import (
	"time"

	"github.com/idoall/gocryptotrader/currency"
	"github.com/idoall/gocryptotrader/exchanges/asset"
)

// Consts here define basic time intervals
const (
	FifteenSecond = Interval(15 * time.Second)
	OneMin        = Interval(time.Minute)
	ThreeMin      = 3 * OneMin
	FiveMin       = 5 * OneMin
	TenMin        = 10 * OneMin
	FifteenMin    = 15 * OneMin
	ThirtyMin     = 30 * OneMin
	OneHour       = Interval(time.Hour)
	TwoHour       = 2 * OneHour
	FourHour      = 4 * OneHour
	SixHour       = 6 * OneHour
	EightHour     = 8 * OneHour
	TwelveHour    = 12 * OneHour
	OneDay        = 24 * OneHour
	ThreeDay      = 3 * OneDay
	SevenDay      = 7 * OneDay
	FifteenDay    = 15 * OneDay
	OneWeek       = 7 * OneDay
	TwoWeek       = 2 * OneWeek
	OneMonth      = 31 * OneDay
	OneYear       = 365 * OneDay
)

const (
	// ErrRequestExceedsExchangeLimits locale for exceeding rate limits message
	ErrRequestExceedsExchangeLimits = "requested data would exceed exchange limits please lower range or use GetHistoricCandlesEx"
)

// Kline K线的映射
type Kline struct {
	Amount    float64   `json:"amount" description:"成交量"`
	Count     int       `json:"count" description:"成交笔数"`
	Open      float64   `json:"open" description:"开盘价"`
	Close     float64   `json:"close" description:"收盘价"`
	Low       float64   `json:"low" description:"最低价"`
	High      float64   `json:"high" description:"最高价"`
	Vol       float64   `json:"vol" description:"成交额,即SUM(每一笔成交价 * 该笔的成交数量)"`
	OpenTime  time.Time `json:"opentime" description:"开盘时间"`
	CloseTime time.Time `json:"closetime" description:"收盘时间"`
}

// Item holds all the relevant information for internal kline elements
type Item struct {
	Exchange string
	Pair     currency.Pair
	Asset    asset.Item
	Interval Interval
	Candles  []Candle
}

// Candle holds historic rate information.
type Candle struct {
	Time   time.Time
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// By Date allows for sorting candle entries by date
type ByDate []Candle

func (b ByDate) Len() int {
	return len(b)
}

func (b ByDate) Less(i, j int) bool {
	return b[i].Time.Before(b[j].Time)
}

func (b ByDate) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// ExchangeCapabilitiesSupported all kline related exchange supported options
type ExchangeCapabilitiesSupported struct {
	Intervals  bool
	DateRanges bool
}

// ExchangeCapabilitiesEnabled all kline related exchange enabled options
type ExchangeCapabilitiesEnabled struct {
	Intervals   map[string]bool `json:"intervals,omitempty"`
	ResultLimit uint32
}

// Interval type for kline Interval usage
type Interval time.Duration

// ErrorKline struct to hold kline interval errors
type ErrorKline struct {
	Asset    asset.Item
	Pair     currency.Pair
	Interval Interval
	Err      error
}

// Error returns short interval unsupported message
func (k *ErrorKline) Error() string {
	return k.Err.Error()
}

// Unwrap returns interval unsupported message
func (k *ErrorKline) Unwrap() error {
	return k.Err
}

// DateRange holds a start and end date for kline usage
type DateRange struct {
	Start time.Time
	End   time.Time
}
