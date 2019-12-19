package exchange

import (
	"sync"

	"github.com/idoall/gocryptotrader/config"
	"github.com/idoall/gocryptotrader/currency"
	"github.com/idoall/gocryptotrader/exchanges/asset"
	"github.com/idoall/gocryptotrader/exchanges/kline"
	"github.com/idoall/gocryptotrader/exchanges/order"
	"github.com/idoall/gocryptotrader/exchanges/orderbook"
	"github.com/idoall/gocryptotrader/exchanges/ticker"
	"github.com/idoall/gocryptotrader/exchanges/websocket/wshandler"
	"github.com/idoall/gocryptotrader/exchanges/withdraw"
)

// IBotExchange enforces standard functions for all exchanges supported in
// GoCryptoTrader
type IBotExchange interface {
	Setup(exch *config.ExchangeConfig) error
	Start(wg *sync.WaitGroup)
	SetDefaults()
	GetName() string
	IsEnabled() bool
	SetEnabled(bool)
	FetchTicker(currency currency.Pair, assetType asset.Item) (*ticker.Price, error)
	UpdateTicker(currency currency.Pair, assetType asset.Item) (*ticker.Price, error)
	FetchOrderbook(currency currency.Pair, assetType asset.Item) (*orderbook.Base, error)
	UpdateOrderbook(currency currency.Pair, assetType asset.Item) (*orderbook.Base, error)
	FetchTradablePairs(assetType asset.Item) ([]string, error)
	UpdateTradablePairs(forceUpdate bool) error
	GetEnabledPairs(assetType asset.Item) currency.Pairs
	GetAvailablePairs(assetType asset.Item) currency.Pairs
	GetAccountInfo() (AccountInfo, error)
	GetAuthenticatedAPISupport(endpoint uint8) bool
	SetPairs(pairs currency.Pairs, assetType asset.Item, enabled bool) error
	GetAssetTypes() asset.Items
	GetExchangeHistory(currencyPair currency.Pair, assetType asset.Item) ([]TradeHistory, error)
	SupportsAutoPairUpdates() bool
	SupportsRESTTickerBatchUpdates() bool
	GetFeeByType(feeBuilder *FeeBuilder) (float64, error)
	GetLastPairsUpdateTime() int64
	GetWithdrawPermissions() uint32
	FormatWithdrawPermissions() string
	SupportsWithdrawPermissions(permissions uint32) bool
	GetFundingHistory() ([]FundHistory, error)
	SubmitOrder(s *order.Submit) (order.SubmitResponse, error)
	ModifyOrder(action *order.Modify) (string, error)
	CancelOrder(order *order.Cancel) error
	CancelAllOrders(orders *order.Cancel) (order.CancelAllResponse, error)
	GetOrderInfo(orderID string) (order.Detail, error)
	GetDepositAddress(cryptocurrency currency.Code, accountID string) (string, error)
	GetOrderHistory(getOrdersRequest *order.GetOrdersRequest) ([]order.Detail, error)
	GetActiveOrders(getOrdersRequest *order.GetOrdersRequest) ([]order.Detail, error)
	WithdrawCryptocurrencyFunds(withdrawRequest *withdraw.CryptoRequest) (string, error)
	WithdrawFiatFunds(withdrawRequest *withdraw.FiatRequest) (string, error)
	WithdrawFiatFundsToInternationalBank(withdrawRequest *withdraw.FiatRequest) (string, error)
	SetHTTPClientUserAgent(ua string)
	GetHTTPClientUserAgent() string
	SetClientProxyAddress(addr string) error
	SupportsWebsocket() bool
	SupportsREST() bool
	IsWebsocketEnabled() bool
	GetWebsocket() (*wshandler.Websocket, error)
	SubscribeToWebsocketChannels(channels []wshandler.WebsocketChannelSubscription) error
	UnsubscribeToWebsocketChannels(channels []wshandler.WebsocketChannelSubscription) error
	AuthenticateWebsocket() error
	GetSubscriptions() ([]wshandler.WebsocketChannelSubscription, error)
	GetDefaultConfig() (*config.ExchangeConfig, error)
	GetBase() *Base
	SupportsAsset(assetType asset.Item) bool
	// GetKlines 自定义获取 K 线
	GetKlines(arg interface{}) ([]*kline.Kline, error)
}
