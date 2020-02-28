package okex

import (
	"fmt"
	"github.com/idoall/gocryptotrader/exchanges/okgroup"
	"net/http"
)

// GetFuturesHistoricCandles 返回合约 K 线
func (o *OKEX) GetFuturesHistoricCandles(request okgroup.GetSpotMarketDataRequest) (resp okgroup.GetSpotMarketDataResponse, _ error) {
	requestURL := fmt.Sprintf("%v/%v/%v%v", okgroup.OKGroupInstruments, request.InstrumentID, okgroup.OKGroupGetSpotMarketData, okgroup.FormatParameters(request))
	return resp, o.SendHTTPRequest(http.MethodGet, okGroupFuturesSubsection, requestURL, nil, &resp, false)
}