package simulator

import (
	"testing"

	"github.com/idoall/gocryptotrader/currency"
	"github.com/idoall/gocryptotrader/exchanges/asset"
	"github.com/idoall/gocryptotrader/exchanges/bitstamp"
)

func TestSimulate(t *testing.T) {
	b := bitstamp.Bitstamp{}
	b.SetDefaults()
	b.Verbose = false
	o, err := b.FetchOrderbook(currency.NewPair(currency.BTC, currency.USD), asset.Spot)
	if err != nil {
		t.Error(err)
	}

	r := o.SimulateOrder(10000000, true)
	t.Log(r.Status)
	r = o.SimulateOrder(2171, false)
	t.Log(r.Status)
}
