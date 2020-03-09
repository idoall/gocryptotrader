package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/idoall/gocryptotrader/common"
	"github.com/idoall/gocryptotrader/config"
	"github.com/idoall/gocryptotrader/core"
	"github.com/idoall/gocryptotrader/dispatch"
	"github.com/idoall/gocryptotrader/engine"
	"github.com/idoall/gocryptotrader/exchanges/huobi"
	"github.com/idoall/gocryptotrader/exchanges/request"
	"github.com/idoall/gocryptotrader/gctscript"
	gctscriptVM "github.com/idoall/gocryptotrader/gctscript/vm"
	gctlog "github.com/idoall/gocryptotrader/log"
	"github.com/idoall/gocryptotrader/portfolio/withdraw"
	"github.com/idoall/gocryptotrader/signaler"
)

func main() {

	// Handle flags
	var settings engine.Settings
	versionFlag := flag.Bool("version", false, "retrieves current GoCryptoTrader version")

	// Core settings
	flag.StringVar(&settings.ConfigFile, "config", config.DefaultFilePath(), "config file to load")
	flag.StringVar(&settings.DataDir, "datadir", common.GetDefaultDataDir(runtime.GOOS), "default data directory for GoCryptoTrader files")
	flag.IntVar(&settings.GoMaxProcs, "gomaxprocs", runtime.GOMAXPROCS(-1), "sets the runtime GOMAXPROCS value")
	flag.BoolVar(&settings.EnableDryRun, "dryrun", false, "dry runs bot, doesn't save config file")
	flag.BoolVar(&settings.EnableAllExchanges, "enableallexchanges", false, "enables all exchanges")
	flag.BoolVar(&settings.EnableAllPairs, "enableallpairs", false, "enables all pairs for enabled exchanges")
	flag.BoolVar(&settings.EnablePortfolioManager, "portfoliomanager", true, "enables the portfolio manager")
	flag.DurationVar(&settings.PortfolioManagerDelay, "portfoliomanagerdelay", time.Duration(0), "sets the portfolio managers sleep delay between updates")
	flag.BoolVar(&settings.EnableGRPC, "grpc", true, "enables the grpc server")
	flag.BoolVar(&settings.EnableGRPCProxy, "grpcproxy", false, "enables the grpc proxy server")
	flag.BoolVar(&settings.EnableWebsocketRPC, "websocketrpc", true, "enables the websocket RPC server")
	flag.BoolVar(&settings.EnableDeprecatedRPC, "deprecatedrpc", true, "enables the deprecated RPC server")
	flag.BoolVar(&settings.EnableCommsRelayer, "enablecommsrelayer", true, "enables available communications relayer")
	flag.BoolVar(&settings.Verbose, "verbose", false, "increases logging verbosity for GoCryptoTrader")
	flag.BoolVar(&settings.EnableExchangeSyncManager, "syncmanager", true, "enables to exchange sync manager")
	flag.BoolVar(&settings.EnableWebsocketRoutine, "websocketroutine", true, "enables the websocket routine for all loaded exchanges")
	flag.BoolVar(&settings.EnableCoinmarketcapAnalysis, "coinmarketcap", false, "overrides config and runs currency analysis")
	flag.BoolVar(&settings.EnableEventManager, "eventmanager", true, "enables the event manager")
	flag.BoolVar(&settings.EnableOrderManager, "ordermanager", true, "enables the order manager")
	flag.BoolVar(&settings.EnableDepositAddressManager, "depositaddressmanager", true, "enables the deposit address manager")
	flag.BoolVar(&settings.EnableConnectivityMonitor, "connectivitymonitor", true, "enables the connectivity monitor")
	flag.BoolVar(&settings.EnableDatabaseManager, "databasemanager", true, "enables database manager")
	flag.BoolVar(&settings.EnableGCTScriptManager, "gctscriptmanager", true, "enables gctscript manager")
	flag.DurationVar(&settings.EventManagerDelay, "eventmanagerdelay", time.Duration(0), "sets the event managers sleep delay between event checking")
	flag.BoolVar(&settings.EnableNTPClient, "ntpclient", true, "enables the NTP client to check system clock drift")
	flag.BoolVar(&settings.EnableDispatcher, "dispatch", true, "enables the dispatch system")
	flag.IntVar(&settings.DispatchMaxWorkerAmount, "dispatchworkers", dispatch.DefaultMaxWorkers, "sets the dispatch package max worker generation limit")
	flag.IntVar(&settings.DispatchJobsLimit, "dispatchjobslimit", dispatch.DefaultJobsLimit, "sets the dispatch package max jobs limit")

	// Exchange syncer settings
	flag.BoolVar(&settings.EnableTickerSyncing, "tickersync", true, "enables ticker syncing for all enabled exchanges")
	flag.BoolVar(&settings.EnableOrderbookSyncing, "orderbooksync", true, "enables orderbook syncing for all enabled exchanges")
	flag.BoolVar(&settings.EnableTradeSyncing, "tradesync", false, "enables trade syncing for all enabled exchanges")
	flag.IntVar(&settings.SyncWorkers, "syncworkers", engine.DefaultSyncerWorkers, "the amount of workers (goroutines) to use for syncing exchange data")
	flag.BoolVar(&settings.SyncContinuously, "synccontinuously", true, "whether to sync exchange data continuously (ticker, orderbook and trade history info")
	flag.DurationVar(&settings.SyncTimeout, "synctimeout", engine.DefaultSyncerTimeout,
		"the amount of time before the syncer will switch from one protocol to the other (e.g. from REST to websocket)")

	// Forex provider settings
	flag.BoolVar(&settings.EnableCurrencyConverter, "currencyconverter", false, "overrides config and sets up foreign exchange Currency Converter")
	flag.BoolVar(&settings.EnableCurrencyLayer, "currencylayer", false, "overrides config and sets up foreign exchange Currency Layer")
	flag.BoolVar(&settings.EnableFixer, "fixer", false, "overrides config and sets up foreign exchange Fixer.io")
	flag.BoolVar(&settings.EnableOpenExchangeRates, "openexchangerates", false, "overrides config and sets up foreign exchange Open Exchange Rates")

	// Exchange tuning settings
	flag.BoolVar(&settings.EnableExchangeAutoPairUpdates, "exchangeautopairupdates", false, "enables automatic available currency pair updates for supported exchanges")
	flag.BoolVar(&settings.DisableExchangeAutoPairUpdates, "exchangedisableautopairupdates", false, "disables exchange auto pair updates")
	flag.BoolVar(&settings.EnableExchangeWebsocketSupport, "exchangewebsocketsupport", false, "enables Websocket support for exchanges")
	flag.BoolVar(&settings.EnableExchangeRESTSupport, "exchangerestsupport", true, "enables REST support for exchanges")
	flag.BoolVar(&settings.EnableExchangeVerbose, "exchangeverbose", false, "increases exchange logging verbosity")
	flag.BoolVar(&settings.ExchangePurgeCredentials, "exchangepurgecredentials", false, "purges the stored exchange API credentials")
	flag.BoolVar(&settings.EnableExchangeHTTPRateLimiter, "ratelimiter", true, "enables the rate limiter for HTTP requests")
	flag.IntVar(&settings.MaxHTTPRequestJobsLimit, "requestjobslimit", int(request.DefaultMaxRequestJobs), "sets the max amount of jobs the HTTP request package stores")
	flag.IntVar(&settings.RequestTimeoutRetryAttempts, "exchangehttptimeoutretryattempts", request.DefaultTimeoutRetryAttempts, "sets the amount of retry attempts after a HTTP request times out")
	flag.DurationVar(&settings.ExchangeHTTPTimeout, "exchangehttptimeout", time.Duration(0), "sets the exchangs HTTP timeout value for HTTP requests")
	flag.StringVar(&settings.ExchangeHTTPUserAgent, "exchangehttpuseragent", "", "sets the exchanges HTTP user agent")
	flag.StringVar(&settings.ExchangeHTTPProxy, "exchangehttpproxy", "", "sets the exchanges HTTP proxy server")
	flag.BoolVar(&settings.EnableExchangeHTTPDebugging, "exchangehttpdebugging", false, "sets the exchanges HTTP debugging")

	// Common tuning settings
	flag.DurationVar(&settings.GlobalHTTPTimeout, "globalhttptimeout", time.Duration(0), "sets common HTTP timeout value for HTTP requests")
	flag.StringVar(&settings.GlobalHTTPUserAgent, "globalhttpuseragent", "", "sets the common HTTP client's user agent")
	flag.StringVar(&settings.GlobalHTTPProxy, "globalhttpproxy", "", "sets the common HTTP client's proxy server")

	// GCTScript tuning settings
	flag.UintVar(&settings.MaxVirtualMachines, "maxvirtualmachines", uint(gctscriptVM.DefaultMaxVirtualMachines), "set max virtual machines that can load")

	// Withdraw Cache tuning settings
	flag.Uint64Var(&settings.WithdrawCacheSize, "withdrawcachesize", withdraw.CacheSize, "set cache size for withdrawal requests")

	flag.Parse()

	if *versionFlag {
		fmt.Print(core.Version(true))
		os.Exit(0)
	}

	fmt.Println(core.Banner)
	fmt.Println(core.Version(false))

	var err error
	settings.CheckParamInteraction = true
	engine.Bot, err = engine.NewFromSettings(&settings)
	if engine.Bot == nil || err != nil {
		log.Fatalf("Unable to initialise bot engine. Error: %s\n", err)
	}

	{

		var huobiExch huobi.HUOBI
		huobiExch.SetDefaults()
		//获取交易所 -- 测试不需要使用 engine，直接使用 实例 ，也可以访问
		// exchCfg, _ := engine.Bot.Config.GetExchangeConfig("Huobi")
		huobiExch.API.AuthenticatedSupport = true
		huobiExch.API.AuthenticatedWebsocketSupport = true
		huobiExch.API.Credentials.Key = ""
		huobiExch.API.Credentials.Secret = ""
		huobiExch.SkipAuthCheck = true
		huobiExch.Verbose = true
		//--------历史委托信息
		// req := huobi.ContractHisordersRequest{
		// 	Symbol:     "BTC",
		// 	TradeType:  huobi.TradeType0,
		// 	Type:       huobi.ContractHisOrderType1,
		// 	Status:     huobi.ContractOrderStatus0,
		// 	CreateDate: 90,
		// }
		// if res, err := huobiExch.GetContractHisorders(req); err != nil {
		// 	panic(err)
		// } else {
		// 	fmt.Println("res", res)
		// }

		//-----下新订单
		req := huobi.ContractNewOrderRequest{}
		req.Symbol = "BTC"
		req.ContractType = huobi.ContractTypeQuarter // 季度合约
		req.OrderPriceType = huobi.ContractOrderPriceTypePostOnly
		req.Direction = huobi.ContractOrderDirectionBuy
		req.Offset = huobi.ContractOrderOffsetOpen
		req.ClientOrderID = time.Now().Unix()
		req.LeverRae = huobi.LeverRae1
		req.Volume = 1
		req.Price = 7500.00

		if res, err := huobiExch.ContractNewOrder(req); err != nil {
			panic(err)
		} else {
			fmt.Println(res)
		}

		//------历史成交记录
		// _timeFormat_local := "2006-01-02 15:04:05"
		// req := huobi.ContractMatchResultsRequest{
		// 	TradeType:  huobi.TradeType0,
		// 	CreateDate: 90,
		// 	PageSize:   50,
		// 	// Symbol:     "BTC",
		// }
		// req.Symbol = "BTC"
		// if res, err := huobiExch.GetContractMatchResults(req); err != nil {
		// 	panic(err)
		// } else {
		// 	fmt.Println(" - ", "历史成交记录", "当前页", res.CurrentPage, "总页数", res.TotalPage, "total_siz", res.TotalSize)
		// 	tradeList := make(map[int64]*huobi.ContractMatchResultDataItem)
		// 	for _, v := range res.Trades {
		// 		var tradeItem *huobi.ContractMatchResultDataItem
		// 		idArray := strings.Split(v.ID, "-")

		// 		tradeID, _ := strconv.ParseInt(idArray[1], 10, 64)
		// 		if tradeList[tradeID] == nil {
		// 			tradeItem = &huobi.ContractMatchResultDataItem{}
		// 		} else {
		// 			tradeItem = tradeList[tradeID]
		// 		}
		// 		tradeItem.CreateDate = v.CreateDate
		// 		tradeItem.Symbol = v.Symbol
		// 		tradeItem.Direction = v.Direction
		// 		tradeItem.Offset = v.Offset
		// 		tradeItem.TradePrice = v.TradePrice
		// 		tradeItem.TradeVolume += v.TradeVolume
		// 		tradeItem.TradeTurnover += v.TradeTurnover
		// 		tradeItem.OffsetProfitloss += v.OffsetProfitloss
		// 		tradeItem.TradeFee += v.TradeFee
		// 		tradeList[tradeID] = tradeItem
		// 	}

		// 	for k, v := range tradeList {
		// 		fmt.Printf("\t ID :%d\n", k)
		// 		fmt.Printf("\t\t 订单时间 :%s\n", time.Unix(0, int64(v.CreateDate)*int64(time.Millisecond)).Format(_timeFormat_local))
		// 		fmt.Printf("\t\t 累计成交数量: %.2f\n", v.TradeVolume)
		// 		fmt.Printf("\t\t 品种代码: %s\n", v.Symbol)

		// 		// 开平方向
		// 		if v.Offset == huobi.ContractOrderOffsetOpen && v.Direction == huobi.ContractOrderDirectionBuy {
		// 			fmt.Printf("\t\t 交易类型: %s %s\n", "开多", "买入开多")
		// 		} else if v.Offset == huobi.ContractOrderOffsetOpen && v.Direction == huobi.ContractOrderDirectionSell {
		// 			fmt.Printf("\t\t 交易类型: %s %s\n", "开空", "卖出开空")
		// 		} else if v.Offset == huobi.ContractOrderOffsetClose && v.Direction == huobi.ContractOrderDirectionBuy {
		// 			fmt.Printf("\t\t 交易类型: %s %s\n", "平空", "买入平空")
		// 		} else if v.Offset == huobi.ContractOrderOffsetClose && v.Direction == huobi.ContractOrderDirectionSell {
		// 			fmt.Printf("\t\t 交易类型: %s %s\n", "平多", "卖出平多")
		// 		}

		// 		fmt.Printf("\t\t 累计成交数量: %.2f\n", v.TradeVolume)
		// 		fmt.Printf("\t\t 成交价格: %.2f\n", v.TradePrice)
		// 		fmt.Printf("\t\t 本笔成交金额: %.2f\n", v.TradeTurnover)
		// 		fmt.Printf("\t\t 平仓盈亏: %.8f\n", v.OffsetProfitloss)
		// 		fmt.Printf("\t\t 成交手续费: %.8f %s\n", v.TradeFee, v.FeeAsset)
		// 	}

		// }
		return
	}

	gctscript.Setup()

	engine.PrintSettings(&engine.Bot.Settings)
	if err = engine.Bot.Start(); err != nil {
		gctlog.Errorf(gctlog.Global, "Unable to start bot engine. Error: %s\n", err)
		os.Exit(1)
	}

	interrupt := signaler.WaitForInterrupt()
	gctlog.Infof(gctlog.Global, "Captured %v, shutdown requested.\n", interrupt)
	engine.Bot.Stop()
	gctlog.Infoln(gctlog.Global, "Exiting.")
}
