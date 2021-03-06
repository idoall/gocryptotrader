package poloniex

import (
	"github.com/idoall/gocryptotrader/common"
	"github.com/idoall/gocryptotrader/exchanges/kline"
)

// GetKlines  checks and returns a requested kline if it exists
func (b *Poloniex) GetKlines(arg interface{}) ([]*kline.Kline, error) {

	var klines []*kline.Kline

	return klines, common.ErrFunctionNotSupported
}
