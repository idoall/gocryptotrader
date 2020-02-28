package huobi

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/idoall/gocryptotrader/common"
)

const (
	huobiContractAPIURL = "https://api.btcgateway.pro"

	huobiGetContractInfo = "contract_contract_info"
	huobiGetContractList = "/market/history/kline"
)

// GetContractInfo 获取合约信息
func (h *HUOBI) GetContractInfo(req ContractInfoRequest) (ContractInfoResponse, error) {
	vals := url.Values{}
	if req.Symbol != "" {
		vals.Set("symbol", req.Symbol)
	}
	if req.ContractCode != "" {
		vals.Set("contract_code", req.ContractCode)
	}
	if req.ContractType != "" {
		vals.Set("contract_type", req.ContractType)
	}

	type response struct {
		Response
		TradeHistory []TradeHistory `json:"data"`
	}

	var result ContractInfoResponse
	urlPath := fmt.Sprintf("%s/%s", huobiContractAPIURL, huobiGetContractList)
	err := h.SendHTTPRequest(common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" || err != nil {
		return result, errors.New(result.ErrorMessage)
	}
	return result, nil
}

// GetContractKline 获取火币合约 Kline
func (h *HUOBI) GetContractKline(arg KlinesRequestParams) ([]KlineItem, error) {
	vals := url.Values{}
	vals.Set("symbol", arg.Symbol)
	vals.Set("period", string(arg.Period))

	if arg.Size != 0 {
		vals.Set("size", strconv.Itoa(arg.Size))
	}

	type response struct {
		Response
		Data []KlineItem `json:"data"`
	}

	var result response
	urlPath := fmt.Sprintf("%s/%s", huobiContractAPIURL, huobiGetContractList)

	err := h.SendHTTPRequest(common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Data, err
}
