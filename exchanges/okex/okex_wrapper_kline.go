package okex

// GetKlines  checks and returns a requested kline if it exists
// func (o *OKEX) GetKlines(arg interface{}) ([]*kline.Kline, error) {

// 	var klines []*kline.Kline

// 	// 判断是否是 struct
// 	if reflect.TypeOf(arg).Kind() != reflect.Struct {
// 		return klines, errors.New("arg argument must be a struct address")
// 	}

// 	// 判断类型是否是 KlinesRequestParams
// 	r, ok := arg.(okgroup.GetSpotMarketDataRequest)
// 	if !ok {
// 		return klines, errors.New("arg argument must be a okgroup.GetSpotMarketDataRequest struct")
// 	}

// 	_timeFormat := "2006-01-02T15:04:05.999Z"
// 	// 解析数据
// 	if candleList, err := o.GetFuturesHistoricCandles(r); err != nil {
// 		return nil, err
// 	} else {
// 		for _, v := range candleList {
// 			item := v.([]interface{})
// 			// ot := v.(string)

// 			_ot, err := time.ParseInLocation(_timeFormat, item[0].(string), time.Local)
// 			if err != nil {
// 				fmt.Println("err", err)
// 			} else {
// 				//fmt.Println("_ot", _ot, item[0].(string))
// 			}

// 			_open, err := convert.FloatFromString(item[1].(string))
// 			if err != nil {
// 				return nil, fmt.Errorf("cannot parse Kline.Open. Err: %s", err)
// 			}
// 			_high, err := convert.FloatFromString(item[2].(string))
// 			if err != nil {
// 				return nil, fmt.Errorf("cannot parse Kline.High. Err: %s", err)
// 			}
// 			_low, err := convert.FloatFromString(item[3].(string))
// 			if err != nil {
// 				return nil, fmt.Errorf("cannot parse Kline.Low. Err: %s", err)
// 			}
// 			_close, err := convert.FloatFromString(item[4].(string))
// 			if err != nil {
// 				return nil, fmt.Errorf("cannot parse Kline.Close. Err: %s", err)
// 			}
// 			_vol, err := convert.FloatFromString(item[5].(string))
// 			if err != nil {
// 				return nil, fmt.Errorf("cannot parse Kline.Volume. Err: %s", err)
// 			}

// 			klines = append(klines, &kline.Kline{
// 				OpenTime: _ot.Add(time.Hour * 8),
// 				Open:     _open,
// 				High:     _high,
// 				Low:      _low,
// 				Close:    _close,
// 				Vol:      _vol,
// 			})
// 		}
// 	}

// 	return klines, nil
// }
