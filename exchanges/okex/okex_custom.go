package okex

import (
	"fmt"
	"github.com/idoall/gocryptotrader/exchanges/okgroup"
	"net/http"
)

// GetFuturesHistoricCandles 返回交割合约 K 线
func (o *OKEX) GetFuturesHistoricCandles(request okgroup.GetSpotMarketDataRequest) (resp okgroup.GetSpotMarketDataResponse, _ error) {
	requestURL := fmt.Sprintf("%v/%v/%v%v", okgroup.OKGroupInstruments, request.InstrumentID, okgroup.OKGroupGetSpotMarketData, okgroup.FormatParameters(request))
	return resp, o.SendHTTPRequest(http.MethodGet, okGroupFuturesSubsection, requestURL, nil, &resp, false)
}


// GetFuturesHistoricCandles 返回永续合约 K 线
func (o *OKEX) GetSwapHistoricCandles(request okgroup.GetSpotMarketDataRequest) (resp okgroup.GetSpotMarketDataResponse, _ error) {
	requestURL := fmt.Sprintf("%v/%v/%v%v", okgroup.OKGroupInstruments, request.InstrumentID, okgroup.OKGroupGetSpotMarketData, okgroup.FormatParameters(request))
	return resp, o.SendHTTPRequest(http.MethodGet, okGroupSwapSubsection, requestURL, nil, &resp, false)
}

//GetHistoricCandles returns rangesize number of candles for the given granularity and pair starting from the latest available
func (o *OKEX) GetSpotHistoricCandles(request okgroup.GetSpotMarketDataRequest) (resp okgroup.GetSpotMarketDataResponse, _ error) {
	requestURL := fmt.Sprintf("%v/%v/%v%v", okgroup.OKGroupInstruments, request.InstrumentID, okgroup.OKGroupGetSpotMarketData, okgroup.FormatParameters(request))
	return resp, o.SendHTTPRequest(http.MethodGet, okGroupSwapSubsection, requestURL, nil, &resp, false)
}