package binance

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/idoall/gocryptotrader/common/convert"
	"github.com/idoall/gocryptotrader/exchanges/asset"
)

// GetAccountSnapshot 查询每日资产快照 (USER_DATA)
func (b *Binance) GetAccountSnapshot(arg AccountSnapshotRequest) (snapshot []AccountSnapshotResponse, err error) {

	type respObj struct {
		Code        int    `json:"code"`
		Msg         string `json:"msg"`
		SnapshotVos []struct {
			Type       string `json:"type"`
			UpdateTime int    `json:"updateTime"`
			Data       struct {
				TotalAssetOfBtc string `json:"totalAssetOfBtc"`
				Balances        []struct {
					Asset  string `json:"asset"`
					Free   string `json:"free"`
					Locked string `json:"locked"`
				}
			} `json:"data"`
		} `json:"snapshotVos"`
	}
	var resp respObj

	params := url.Values{}
	params.Set("type", strings.ToUpper(arg.Type.String()))
	if arg.Limit == 0 {
		params.Set("limit", "10")
	} else {
		params.Set("limit", strconv.FormatInt(arg.Limit, 10))
	}
	if arg.StartTime != 0 {
		params.Set("startTime", strconv.FormatInt(arg.StartTime, 10))
	}
	if arg.EndTime != 0 {
		params.Set("endTime", strconv.FormatInt(arg.EndTime, 10))
	}

	path := b.API.Endpoints.URL + accountSnapshot

	if err = b.SendAuthHTTPRequest(http.MethodGet, path, params, limitOrder, &resp); err != nil {
		return
	}

	for _, v := range resp.SnapshotVos {
		var totalAssetOfBtc float64
		if totalAssetOfBtc, err = strconv.ParseFloat(v.Data.TotalAssetOfBtc, 64); err != nil {
			return
		}
		var free float64
		if free, err = strconv.ParseFloat(v.Data.Balances[0].Free, 64); err != nil {
			return
		}
		var locked int64
		if locked, err = strconv.ParseInt(v.Data.Balances[0].Locked, 10, 64); err != nil {
			return
		}
		var updateTime time.Time
		if updateTime, err = convert.TimeFromUnixTimestampFloat(float64(v.UpdateTime)); err != nil {
			return
		}
		snapshot = append(snapshot, AccountSnapshotResponse{
			Asset:           asset.Item(v.Type),
			Symbol:          v.Data.Balances[0].Asset,
			TotalAssetOfBtc: totalAssetOfBtc,
			Free:            free,
			Locked:          locked,
			UpdateTime:      updateTime,
		})
	}

	return
}
