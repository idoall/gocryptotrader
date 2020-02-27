package okex

import (
	"fmt"
	"time"

	"github.com/idoall/gocryptotrader/common/convert"
	"github.com/idoall/gocryptotrader/currency"
	exchange "github.com/idoall/gocryptotrader/exchanges"
	"github.com/idoall/gocryptotrader/exchanges/okgroup"
)

// GetHistoricCandles returns rangesize number of candles for the given granularity and pair starting from the latest available
func (o *OKEX) GetHistoricCandles(pair currency.Pair, rangesize, granularity int64) ([]exchange.Candle, error) {
	// o.GetAssetTypeFromTableName()
	var result []exchange.Candle
	r := okgroup.GetSpotMarketDataRequest{
		InstrumentID: pair.Upper().String(),
		Granularity:  granularity,
	}
	// var list []exchange.Candle
	_timeFormat := "2006-01-02T15:04:05.000Z"
	_loc, _ := time.LoadLocation("Local")
	if list, err := o.GetSpotMarketData(r); err != nil {
		return nil, err
	} else {
		for _, v := range list {
			item := v.([]interface{})
			// ot := v.(string)

			_ot, err := time.ParseInLocation(_timeFormat, item[0].(string), _loc)
			if err != nil {
				fmt.Println("err", err)
			} else {
				// fmt.Println("_ot", _ot)
			}

			_open, err := convert.FloatFromString(item[1].(string))
			if err != nil {
				return nil, fmt.Errorf("cannot parse Kline.Open. Err: %s", err)
			}
			_high, err := convert.FloatFromString(item[2].(string))
			if err != nil {
				return nil, fmt.Errorf("cannot parse Kline.High. Err: %s", err)
			}
			_low, err := convert.FloatFromString(item[3].(string))
			if err != nil {
				return nil, fmt.Errorf("cannot parse Kline.Low. Err: %s", err)
			}
			_close, err := convert.FloatFromString(item[4].(string))
			if err != nil {
				return nil, fmt.Errorf("cannot parse Kline.Close. Err: %s", err)
			}
			_vol, err := convert.FloatFromString(item[5].(string))
			if err != nil {
				return nil, fmt.Errorf("cannot parse Kline.Volume. Err: %s", err)
			}

			result = append(result, exchange.Candle{
				Time:   _ot.Unix(),
				Open:   _open,
				High:   _high,
				Low:    _low,
				Close:  _close,
				Volume: _vol,
			})
		}
	}
	return result, nil
	// return nil, common.ErrFunctionNotSupported
}
