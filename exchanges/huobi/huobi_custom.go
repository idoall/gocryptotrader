package huobi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/idoall/gocryptotrader/common"
	"github.com/idoall/gocryptotrader/common/crypto"
	exchange "github.com/idoall/gocryptotrader/exchanges"
	"github.com/idoall/gocryptotrader/exchanges/request"
)

const (
	HuobiContractAPIURL = "https://api.btcgateway.pro"

	huobiContractInfo                = "contract_contract_info"
	huobiContractAccountInfo         = "contract_account_info"
	huobiGetContractKlineList        = "/market/history/kline"
	huobiContractHisorders           = "contract_hisorders"             //获取历史委托
	huobiContractMatchResults        = "contract_matchresults"          //获取历史成交记录
	huobiContractOpenOrders          = "contract_openorders"            //获取当前未成效委托
	huobiContractTriggerOpenOrders   = "contract_trigger_openorders"    //获取计划委托当前委托
	huobiContractNewOrder            = "contract_order"                 //合约下单
	huobiContractNewTriggerOrder     = "contract_trigger_orders"        //合约计划委托下单
	huobiContractAccountPositionInfo = "contract_account_position_info" // 查询用户帐号和持仓信息
)

// GetContractAccountPositionInfo 查询用户帐号和持仓信息
func (h *HUOBI) GetContractAccountPositionInfo(req SymbolBaseType) ([]ContractAccountPositionInfoResponse, error) {
	type response struct {
		Response
		Data []ContractAccountPositionInfoResponse `json:"data"`
	}
	var result response
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractAccountPositionInfo, nil, req, &result, false)
	return result.Data, err
}

// GetContractHisorders 获取历史委托
func (h *HUOBI) GetContractHisorders(req ContractHisordersRequest) (ContractHisordersData, error) {
	type response struct {
		Response
		Data ContractHisordersData `json:"data"`
	}
	var result response
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractHisorders, nil, req, &result, false)
	return result.Data, err
}

// ContractNewOrder 合约下单
func (h *HUOBI) ContractNewOrder(req ContractNewOrderRequest) (ContractNewOrderResponse, error) {
	var result ContractNewOrderResponse
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractNewOrder, nil, req, &result, false)
	return result, err
}

// ContractNewTriggerOrder 合约计划委托下单
func (h *HUOBI) ContractNewTriggerOrder(req ContractNewTriggerOrderRequest) (ContractNewTriggerOrderResponse, error) {
	var result ContractNewTriggerOrderResponse
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractNewTriggerOrder, nil, req, &result, false)
	return result, err
}

// GetContractMatchResults 获取历史成交记录
func (h *HUOBI) GetContractMatchResults(req ContractMatchResultsRequest) (ContractMatchResultData, error) {
	type response struct {
		Response
		Data ContractMatchResultData `json:"data"`
	}
	var result response
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractMatchResults, nil, req, &result, false)
	return result.Data, err
}

// GetContractOpenOrders 获取火币合约 当前未成交委托
func (h *HUOBI) GetContractOpenOrders(symbol string, pageIndex, pageSize int) (ContractOpenOrderData, error) {
	vals := url.Values{}
	vals.Set("symbol", symbol)
	if pageIndex != 0 {
		vals.Set("pageIndex", strconv.Itoa(pageIndex))
	}

	if pageSize != 0 {
		vals.Set("size", strconv.Itoa(pageSize))
	}

	type response struct {
		Response
		Data ContractOpenOrderData `json:"data"`
	}

	var result response
	fmt.Println(vals)
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractOpenOrders, nil, vals, &result, false)
	return result.Data, err
}

// GetContractTriggerOpenOrders 获取火币合约 获取计划委托当前委托
func (h *HUOBI) GetContractTriggerOpenOrders(symbol string, pageIndex, pageSize int) (ContractOpenOrderData, error) {
	vals := url.Values{}
	vals.Set("symbol", symbol)

	if pageIndex != 0 {
		vals.Set("pageIndex", strconv.Itoa(pageIndex))
	}

	if pageSize != 0 {
		vals.Set("size", strconv.Itoa(pageSize))
	}

	type response struct {
		Response
		Data ContractOpenOrderData `json:"data"`
	}

	var result response
	fmt.Println(vals)
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractTriggerOpenOrders, nil, vals, &result, false)
	return result.Data, err
}

// GetContractAccountInfo 获取用户帐户信息
func (h *HUOBI) GetContractAccountInfo(req ContractAccountInfoRequest) ([]ContractAccountInfoResponseDataItem, error) {
	type response struct {
		Response
		Data []ContractAccountInfoResponseDataItem `json:"data"`
	}
	var result response
	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractAccountInfo, nil, req, &result, false)
	return result.Data, err
}

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
	urlPath := fmt.Sprintf("%s/%s", HuobiContractAPIURL, huobiContractInfo)
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
	urlPath := fmt.Sprintf("%s/%s", HuobiContractAPIURL, huobiGetContractKlineList)

	err := h.SendHTTPRequest(common.EncodeURLValues(urlPath, vals), &result)
	if result.ErrorMessage != "" {
		return nil, errors.New(result.ErrorMessage)
	}
	return result.Data, err
}

// SendContractAuthenticatedHTTPRequest sends authenticated requests to the HUOBI API
func (h *HUOBI) SendContractAuthenticatedHTTPRequest(method, endpoint string, values url.Values, data, result interface{}, isVersion2API bool) error {
	if !h.AllowAuthenticatedRequest() {
		return fmt.Errorf(exchange.WarningAuthenticatedRequestWithoutCredentialsSet, h.Name)
	}

	if values == nil {
		values = url.Values{}
	}

	values.Set("AccessKeyId", h.API.Credentials.Key)
	values.Set("SignatureMethod", "HmacSHA256")
	values.Set("SignatureVersion", "2")
	values.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))

	if isVersion2API {
		endpoint = fmt.Sprintf("/v%s/%s", huobiAPIVersion2, endpoint)
	} else {
		endpoint = fmt.Sprintf("/api/v%s/%s", huobiAPIVersion, endpoint)
	}

	payload := fmt.Sprintf("%s\napi.btcgateway.pro\n%s\n%s",
		method, endpoint, values.Encode())

	headers := make(map[string]string)

	if method == http.MethodGet {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else {
		headers["Content-Type"] = "application/json"
	}

	hmac := crypto.GetHMAC(crypto.HashSHA256, []byte(payload), []byte(h.API.Credentials.Secret))
	values.Set("Signature", crypto.Base64Encode(hmac))
	urlPath := HuobiContractAPIURL + common.EncodeURLValues(endpoint, values)

	var body []byte
	if data != nil {
		encoded, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = encoded
	}
	interim := json.RawMessage{}
	err := h.SendPayload(&request.Item{
		Method:        method,
		Path:          urlPath,
		Headers:       headers,
		Body:          bytes.NewReader(body),
		Result:        &interim,
		AuthRequest:   true,
		Verbose:       h.Verbose,
		HTTPDebugging: h.HTTPDebugging,
		HTTPRecording: h.HTTPRecording,
	})
	if err != nil {
		return err
	}

	if isVersion2API {
		var errCap ResponseV2
		if err = json.Unmarshal(interim, &errCap); err == nil {
			if errCap.Code != 200 && errCap.Message != "" {
				return errors.New(errCap.Message)
			}
		}
	} else {
		var errCap Response
		if err = json.Unmarshal(interim, &errCap); err == nil {
			if errCap.Status == huobiStatusError && errCap.ErrorMessage != "" {
				return errors.New(errCap.ErrorMessage)
			}
		}
	}
	return json.Unmarshal(interim, result)
}
