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

	huobiContractInfo         = "contract_contract_info"
	huobiContractAccountInfo  = "contract_account_info"
	huobiGetContractKlineList = "/market/history/kline"
)

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
	urlPath := h.API.Endpoints.URL + common.EncodeURLValues(endpoint, values)

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
