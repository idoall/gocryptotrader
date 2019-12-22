package huobi

import (
	"errors"
	"reflect"

	"github.com/idoall/gocryptotrader/common/convert"
	"github.com/idoall/gocryptotrader/exchanges/kline"
)

// GetKlines  checks and returns a requested kline if it exists
func (b *HUOBI) GetKlines(arg interface{}) ([]*kline.Kline, error) {

	var klines []*kline.Kline

	// 判断是否是 struct
	if reflect.TypeOf(arg).Kind() != reflect.Struct {
		return klines, errors.New("arg argument must be a struct address")
	}

	// 判断类型是否是 KlinesRequestParams
	klineParams, ok := arg.(KlinesRequestParams)
	if !ok {
		return klines, errors.New("arg argument must be a KlinesRequestParams struct")
	}

	// 获取数据
	candleStickList, err := b.GetSpotKline(klineParams)
	if err != nil {
		return klines, err
	}

	// 解析数据
	for _, v := range candleStickList {
		klines = append(klines,
			&kline.Kline{
				Open:      v.Open,
				Close:     v.Close,
				High:      v.High,
				Low:       v.Low,
				Vol:       v.Volume,
				Amount:    v.Amount,
				CloseTime: convert.UnixTimestampToTime(v.ID),
			},
		)
	}

	return klines, nil
}
