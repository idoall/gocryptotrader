package huobi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/idoall/gocryptotrader/common"
	"github.com/idoall/gocryptotrader/common/crypto"
	"github.com/idoall/gocryptotrader/currency"
	exchange "github.com/idoall/gocryptotrader/exchanges"
	"github.com/idoall/gocryptotrader/exchanges/asset"
	"github.com/idoall/gocryptotrader/exchanges/order"
	"github.com/idoall/gocryptotrader/exchanges/request"
)

const (
	// HuobiContractAPIURL = "https://api.btcgateway.pro"

	// huobiContractInfo                = "contract_contract_info"
	// huobiContractAccountInfo         = "contract_account_info"
	// huobiGetContractKlineList        = "/market/history/kline"
	// huobiContractHisorders           = "contract_hisorders"             //获取历史委托
	// huobiContractMatchResults        = "contract_matchresults"          //获取历史成交记录
	// huobiContractOpenOrders          = "contract_openorders"            //获取当前未成效委托
	// huobiContractTriggerOpenOrders   = "contract_trigger_openorders"    //获取计划委托当前委托
	// huobiContractNewOrder            = "contract_order"                 //合约下单
	// huobiContractNewTriggerOrder     = "contract_trigger_orders"        //合约计划委托下单
	// huobiContractAccountPositionInfo = "contract_account_position_info" // 查询用户帐号和持仓信息

	huobiAccountAssetValuation = "account/asset-valuation"
	huobiAccountGetUID         = "user/uid"
)

// GetUID 查询单个订单详情
func (h *HUOBI) GetUID() (userID int64, err error) {
	resp := struct {
		Code int64 `json:"code"`
		Data int64 `json:"data"`
	}{}

	err = h.SendAuthenticatedHTTPRequest(http.MethodGet, huobiAccountGetUID, url.Values{}, nil, &resp, true)
	if err != nil {
		return
	}
	userID = resp.Data
	return
}

// SearchOrder 查询单个订单详情
func (h *HUOBI) SearchOrder(orderID int64) (order.Detail, error) {
	var orderDetail order.Detail

	resp := struct {
		Order OrderInfo `json:"data"`
	}{}
	if err := h.SendAuthenticatedHTTPRequest(http.MethodGet,
		huobiGetOrders+"/"+strconv.FormatInt(orderID, 10),
		url.Values{},
		nil,
		&resp,
		false); err != nil {
		return orderDetail, err
	}

	typeDetails := strings.Split(resp.Order.Type, "-")
	orderSide, err := order.StringToOrderSide(typeDetails[0])
	if err != nil {
		return orderDetail, err
	}

	orderType, err := order.StringToOrderType(typeDetails[1])
	if err != nil {
		return orderDetail, err
	}

	var orderStatus order.Status
	if strings.EqualFold(resp.Order.State, "submitted") {
		orderStatus = order.New
	} else if strings.EqualFold(resp.Order.State, "partial-filled") {
		orderStatus = order.PartiallyFilled
	} else {
		orderStatus, err = order.StringToOrderStatus(resp.Order.State)
		if err != nil {
			return orderDetail, err
		}
	}

	var p currency.Pair
	var a asset.Item
	p, a, err = h.GetRequestFormattedPairAndAssetType(resp.Order.Symbol)
	if err != nil {
		return orderDetail, err
	}

	orderDetail = order.Detail{
		Exchange:    h.Name,
		ID:          strconv.FormatInt(resp.Order.ID, 10),
		AccountID:   strconv.FormatInt(resp.Order.AccountID, 10),
		Pair:        p,
		Type:        orderType,
		Side:        orderSide,
		Date:        time.Unix(0, resp.Order.CreatedAt*int64(time.Millisecond)),
		Status:      orderStatus,
		Price:       resp.Order.Price,
		Amount:      resp.Order.Amount,
		Cost:        resp.Order.FieldCashAmount, //已成交总金额
		Fee:         resp.Order.FieldFees,       //已成交手续费（买入为币，卖出为钱）
		LastUpdated: time.Unix(0, resp.Order.FinishedAt*int64(time.Millisecond)),
		AssetType:   a,
	}
	return orderDetail, err
}

// GetAccountAssetValuation 获取账户资产估值
// @accountType spot：现货账户， margin：逐仓杠杆账户，otc：OTC 账户，super-margin：全仓杠杆账户
// @valuationCurrency 可选法币有：BTC、CNY、USD、JPY、KRW、GBP、TRY、EUR、RUB、VND、HKD、TWD、MYR、SGD、AED、SAR （大小写敏感）
func (h *HUOBI) GetAccountAssetValuation(accountType asset.Item, valuationCurrency string) (*AccountAssetValuationResponse, error) {
	resp := struct {
		Code int64 `json:"code"`
		OK   bool  `json:"ok"`
		Data struct {
			Balance   string `json:"balance"`
			TimeStamp int64  `json:"timestamp"`
		} `json:"data"`
	}{}
	vals := url.Values{}
	vals.Set("accountType", string(accountType))
	if valuationCurrency != "" {
		vals.Set("valuationCurrency", valuationCurrency)
	}

	err := h.SendAuthenticatedHTTPRequest(http.MethodGet, huobiAccountAssetValuation, vals, nil, &resp, true)
	if err != nil {
		return nil, err
	}

	var balance float64
	if balance, err = strconv.ParseFloat(resp.Data.Balance, 64); err != nil {
		return nil, err
	}

	return &AccountAssetValuationResponse{
		Balance: balance,
		Date:    time.Unix(0, resp.Data.TimeStamp*int64(time.Millisecond)),
	}, err

}

// // GetContractAccountPositionInfo 查询用户帐号和持仓信息
// func (h *HUOBI) GetContractAccountPositionInfo(req SymbolBaseType) ([]ContractAccountPositionInfoResponse, error) {
// 	type response struct {
// 		Response
// 		Data []ContractAccountPositionInfoResponse `json:"data"`
// 	}
// 	var result response
// 	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractAccountPositionInfo, nil, req, &result, false)
// 	return result.Data, err
// }

// // GetContractHisorders 获取历史委托
// func (h *HUOBI) GetContractHisorders(req ContractHisordersRequest) (ContractHisordersData, error) {
// 	type response struct {
// 		Response
// 		Data ContractHisordersData `json:"data"`
// 	}
// 	var result response
// 	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractHisorders, nil, req, &result, false)
// 	return result.Data, err
// }

// // ContractNewOrder 合约下单
// func (h *HUOBI) ContractNewOrder(req ContractNewOrderRequest) (ContractNewOrderResponse, error) {
// 	var result ContractNewOrderResponse
// 	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractNewOrder, nil, req, &result, false)
// 	return result, err
// }

// // ContractNewTriggerOrder 合约计划委托下单
// func (h *HUOBI) ContractNewTriggerOrder(req ContractNewTriggerOrderRequest) (ContractNewTriggerOrderResponse, error) {
// 	var result ContractNewTriggerOrderResponse
// 	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractNewTriggerOrder, nil, req, &result, false)
// 	return result, err
// }

// // GetContractMatchResults 获取历史成交记录
// func (h *HUOBI) GetContractMatchResults(req ContractMatchResultsRequest) (ContractMatchResultData, error) {
// 	type response struct {
// 		Response
// 		Data ContractMatchResultData `json:"data"`
// 	}
// 	var result response
// 	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractMatchResults, nil, req, &result, false)
// 	return result.Data, err
// }

// // GetContractOpenOrders 获取火币合约 当前未成交委托
// func (h *HUOBI) GetContractOpenOrders(symbol string, pageIndex, pageSize int) (ContractOpenOrderData, error) {
// 	vals := url.Values{}
// 	vals.Set("symbol", symbol)
// 	if pageIndex != 0 {
// 		vals.Set("pageIndex", strconv.Itoa(pageIndex))
// 	}

// 	if pageSize != 0 {
// 		vals.Set("size", strconv.Itoa(pageSize))
// 	}

// 	type response struct {
// 		Response
// 		Data ContractOpenOrderData `json:"data"`
// 	}

// 	var result response
// 	fmt.Println(vals)
// 	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractOpenOrders, nil, vals, &result, false)
// 	return result.Data, err
// }

// // GetContractTriggerOpenOrders 获取火币合约 获取计划委托当前委托
// func (h *HUOBI) GetContractTriggerOpenOrders(symbol string, pageIndex, pageSize int) (ContractOpenOrderData, error) {
// 	vals := url.Values{}
// 	vals.Set("symbol", symbol)

// 	if pageIndex != 0 {
// 		vals.Set("pageIndex", strconv.Itoa(pageIndex))
// 	}

// 	if pageSize != 0 {
// 		vals.Set("size", strconv.Itoa(pageSize))
// 	}

// 	type response struct {
// 		Response
// 		Data ContractOpenOrderData `json:"data"`
// 	}

// 	var result response
// 	fmt.Println(vals)
// 	err := h.SendContractAuthenticatedHTTPRequest(http.MethodPost, huobiContractTriggerOpenOrders, nil, vals, &result, false)
// 	return result.Data, err
// }

// GetAccountInfoContract 获取用户帐户信息
func (h *HUOBI) GetAccountInfoContract(assetType asset.Item, contractCode string) ([]AccountInfoResponseDataItem, error) {
	type response struct {
		Response
		Data []AccountInfoResponseDataItem `json:"data"`
	}
	params := url.Values{}
	if contractCode != "" {
		params.Set("contract_code", contractCode)
	}

	var result response
	err := h.SendAuthenticatedHTTPRequestContract(http.MethodPost, assetType, huobiAccountInfoContract, params, nil, &result, false)
	return result.Data, err
}

// // GetContractInfo 获取合约信息
// func (h *HUOBI) GetContractInfo(req ContractInfoRequest) (ContractInfoResponse, error) {
// 	vals := url.Values{}
// 	if req.Symbol != "" {
// 		vals.Set("symbol", req.Symbol)
// 	}
// 	if req.ContractCode != "" {
// 		vals.Set("contract_code", req.ContractCode)
// 	}
// 	if req.ContractType != "" {
// 		vals.Set("contract_type", req.ContractType)
// 	}

// 	type response struct {
// 		Response
// 		TradeHistory []TradeHistory `json:"data"`
// 	}

// 	var result ContractInfoResponse
// 	urlPath := fmt.Sprintf("%s/%s", HuobiContractAPIURL, huobiContractInfo)
// 	err := h.SendHTTPRequest(common.EncodeURLValues(urlPath, vals), &result)
// 	if result.ErrorMessage != "" || err != nil {
// 		return result, errors.New(result.ErrorMessage)
// 	}
// 	return result, nil
// }

// // GetContractKline 获取火币合约 Kline
// func (h *HUOBI) GetContractKline(arg KlinesRequestParams) ([]KlineItem, error) {
// 	vals := url.Values{}
// 	vals.Set("symbol", arg.Symbol)
// 	vals.Set("period", string(arg.Period))

// 	if arg.Size != 0 {
// 		vals.Set("size", strconv.Itoa(arg.Size))
// 	}

// 	type response struct {
// 		Response
// 		Data []KlineItem `json:"data"`
// 	}

// 	var result response
// 	urlPath := fmt.Sprintf("%s/%s", HuobiContractAPIURL, huobiGetContractKlineList)

// 	err := h.SendHTTPRequest(common.EncodeURLValues(urlPath, vals), &result)
// 	if result.ErrorMessage != "" {
// 		return nil, errors.New(result.ErrorMessage)
// 	}
// 	return result.Data, err
// }

func (h *HUOBI) SendAuthenticatedHTTPRequestContract(method string, assetType asset.Item, endpoint string, values url.Values, data, result interface{}, isVersion2API bool) error {
	if !h.AllowAuthenticatedRequest() {
		return fmt.Errorf(exchange.WarningAuthenticatedRequestWithoutCredentialsSet, h.Name)
	}

	if values == nil {
		values = url.Values{}
	}

	now := time.Now()
	values.Set("AccessKeyId", h.API.Credentials.Key)
	values.Set("SignatureMethod", "HmacSHA256")
	values.Set("SignatureVersion", "2")
	values.Set("Timestamp", now.UTC().Format("2006-01-02T15:04:05"))

	if isVersion2API {
		if assetType == asset.Future {
			endpoint = fmt.Sprintf("/%s/v%s/%s", futureRESTBasePath, huobiAPIVersion2, endpoint)
		} else if assetType == asset.PerpetualContract {
			endpoint = fmt.Sprintf("/%s/v%s/%s", perpetualRESTBasePath, huobiAPIVersion2, endpoint)
		}
	} else {
		if assetType == asset.Future {
			endpoint = fmt.Sprintf("/%s/v%s/%s", futureRESTBasePath, huobiAPIVersion, endpoint)
		} else if assetType == asset.PerpetualContract {
			endpoint = fmt.Sprintf("/%s/v%s/%s", perpetualRESTBasePath, huobiAPIVersion, endpoint)
		}
	}

	var payload string
	if assetType == asset.Future {
		payload = fmt.Sprintf("%s\n%s\n%s\n%s",
			method, futureApiURL, endpoint, values.Encode())
	} else if assetType == asset.PerpetualContract {
		payload = fmt.Sprintf("%s\n%s\n%s\n%s",
			method, perpetualApiURL, endpoint, values.Encode())
	}

	headers := make(map[string]string)

	if method == http.MethodGet {
		headers["Content-Type"] = "application/x-www-form-urlencoded"
	} else {
		headers["Content-Type"] = "application/json"
	}

	hmac := crypto.GetHMAC(crypto.HashSHA256, []byte(payload), []byte(h.API.Credentials.Secret))
	values.Set("Signature", crypto.Base64Encode(hmac))

	var urlPath string
	if assetType == asset.Future {
		urlPath = "https://" + futureApiURL + common.EncodeURLValues(endpoint, values)
	} else if assetType == asset.PerpetualContract {
		urlPath = "https://" + perpetualApiURL + common.EncodeURLValues(endpoint, values)
	}

	var body []byte
	if data != nil {
		encoded, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = encoded
	}

	// Time difference between your timestamp and standard should be less than 1 minute.
	ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Minute))
	defer cancel()
	interim := json.RawMessage{}
	err := h.SendPayload(ctx, &request.Item{
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
