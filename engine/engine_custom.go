package engine

import (
	"time"

	"github.com/idoall/gocryptotrader/exchanges/kline"
	"github.com/idoall/gocryptotrader/exchanges/okex"
	"github.com/idoall/gocryptotrader/exchanges/okgroup"
)

// getOKEXKline 拼接多条 kline 数据返回
func getOKEXKline(request okgroup.GetSpotMarketDataRequest, okExch okex.OKEX, listLength int) ([]*kline.Kline, error) {

	var klineList []*kline.Kline

	if listLength == 0 {
		listLength = 200
	}
	_timeFormat_ok := "2006-01-02T15:04:05.999Z"
	for {
		list, err := okExch.GetKlines(request)
		if err != nil {
			return nil, err
		}
		for k, v := range list {
			//第2次累计加的时候，会多读取1条，为了避免读取多个TimeInternal判断，第2次直接去掉第1条
			if len(klineList) > 0 && k == 0 {
				continue
			}
			klineList = append(klineList, v)
		}

		if len(klineList) < listLength {
			request.End = list[len(list)-1].OpenTime.Add(-time.Hour * 8).Format(_timeFormat_ok)
		} else {
			return klineList, nil
		}
	}

}
