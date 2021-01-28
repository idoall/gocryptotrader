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
	"github.com/idoall/gocryptotrader/exchanges/request"
	"github.com/idoall/gocryptotrader/exchanges/trade"
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
	flag.IntVar(&settings.RequestMaxRetryAttempts, "httpmaxretryattempts", request.DefaultMaxRetryAttempts, "sets the number of retry attempts after a retryable HTTP failure")
	flag.DurationVar(&settings.HTTPTimeout, "httptimeout", time.Duration(0), "sets the HTTP timeout value for HTTP requests")
	flag.StringVar(&settings.HTTPUserAgent, "httpuseragent", "", "sets the HTTP user agent")
	flag.StringVar(&settings.HTTPProxy, "httpproxy", "", "sets the HTTP proxy server")
	flag.BoolVar(&settings.EnableExchangeHTTPDebugging, "exchangehttpdebugging", false, "sets the exchanges HTTP debugging")
	flag.DurationVar(&settings.TradeBufferProcessingInterval, "tradeprocessinginterval", trade.DefaultProcessorIntervalTime, "sets the interval to save trade buffer data to the database")

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

	// collect flags
	flagSet := make(map[string]bool)
	// Stores the set flags
	flag.Visit(func(f *flag.Flag) { flagSet[f.Name] = true })
	if !flagSet["config"] {
		// If config file is not explicitly set, fall back to default path resolution
		settings.ConfigFile = ""
	}

	engine.Bot, err = engine.NewFromSettings(&settings, flagSet)
	if engine.Bot == nil || err != nil {
		log.Fatalf("Unable to initialise bot engine. Error: %s\n", err)
	}

	// {

	// var exch binance.Binance
	// exch.SetDefaults()
	// //获取交易所 -- 测试不需要使用 engine，直接使用 实例 ，也可以访问
	// exchCfg, _ := engine.Bot.Config.GetExchangeConfig("Binance")
	// exchCfg.Verbose = true
	// exchCfg.Features.Enabled.Websocket = true
	// exchCfg.AuthenticatedWebsocketAPISupport = &exchCfg.Features.Enabled.Websocket
	// exch.API.AuthenticatedSupport = true
	// exch.API.AuthenticatedWebsocketSupport = true

	// exch.SkipAuthCheck = true
	// exch.Verbose = true
	// logCfg := gctlog.GenDefaultSettings()
	// gctlog.GlobalLogConfig = &logCfg
	// exch.Setup(exchCfg)
	// exch.WebsocketFuture.SetCanUseAuthenticatedEndpoints(true)

	// exch.CurrencyPairs.Pairs[asset.Future] = &currency.PairStore{
	// 	RequestFormat: &currency.PairFormat{
	// 		Uppercase: true,
	// 	},
	// 	ConfigFormat: &currency.PairFormat{
	// 		Uppercase: true,
	// 	},
	// }
	// if err = exch.CurrencyPairs.SetAssetEnabled(asset.Future, true); err != nil {
	// 	panic("exch.CurrencyPairs.SetAssetEnabled Error")
	// }
	// symbolPair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolPair.Delimiter = ""
	// exch.CurrencyPairs.Pairs[asset.Future].Available = exch.CurrencyPairs.Pairs[asset.Future].Available.Add(symbolPair)
	// exch.CurrencyPairs.Pairs[asset.Future].Enabled = exch.CurrencyPairs.Pairs[asset.Future].Enabled.Add(symbolPair)
	// // err = exch.CurrencyPairs.EnablePair(asset.Future, symbolPair)
	// // if err != nil {
	// // 	panic(err)
	// // }
	// go func() {

	// 	err = exch.WebsocketFuture.Connect()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }()
	// interruptx := signaler.WaitForInterrupt()
	// gctlog.Infof(gctlog.Global, "Captured %v, shutdown requested.\n", interruptx)
	// return
	// 持仓ADL队列估算
	// symbolFuturePair := currency.NewPair(currency.NewCode("XLM"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// list, err := exch.ADLQuantile(symbolFuturePair)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("%+v\n", list)
	// }
	// return

	// 调整逐仓保证金
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// list, err := exch.PositionMargin(binance.PositionMarginRequest{
	// 	Symbol:       symbolFuturePair,
	// 	PositionSide: binance.PositionSideSHORT,
	// 	Type:         binance.PositionMarginTypeSub,
	// 	Amount:       10,
	// })
	// if err != nil {
	// 	panic(err)
	// } else {
	// fmt.Printf("%+v\n", list)
	// }

	// 逐仓保证金变动历史 (TRADE)
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// list, err := exch.PositionMarginHistory(binance.PositionMarginHistoryRequest{
	// 	Symbol: symbolFuturePair,
	// 	// PositionSide: binance.PositionSideSHORT,
	// 	// Type:         binance.PositionMarginTypeAdd,
	// 	// Amount:       10,
	// })
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list {
	// 		fmt.Printf("%+v\n", v)
	// 	}
	// }
	// return

	// 获取资金费率表
	// list, err := exch.GetFutureFundingRate(binance.FutureFundingRateRequest{
	// 	Limit: 1,
	// })
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list {
	// 		fmt.Printf("%+v\n", v)
	// 	}

	// }
	// return

	// 获取用户手续费率
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""

	// symbolPair := currency.NewPair(currency.NewCode("SXP"), currency.NewCode("USDT"))
	// // symbolPair.Delimiter = "-"
	// list, err := exch.MarginTypeFuture(symbolPair, binance.MarginType_CROSSED)
	// if err != nil {
	// 	panic(err)
	// } else {

	// 	fmt.Printf("%+v\n", list)

	// }
	// return

	//获取交易规则和交易对
	// list, err := exch.GetExchangeInfo(asset.Spot)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	minPrice := 0
	// 	minBase := 0
	// 	minQuoue := 0
	// 	for _, v := range list.Symbols {
	// 		if v.BaseAsset == "ETH" && v.QuoteAsset == "USDT" {
	// 			fmt.Printf("%+v\n", v)
	// 			for _, vv := range v.Filters {
	// 				if vv.FilterType == "PRICE_FILTER" {
	// 					minPrice = len(strings.Split(decimal.NewFromFloat(vv.MinPrice).String(), ".")[1])
	// 				}
	// 				if vv.FilterType == "LOT_SIZE" {
	// 					minBase = len(strings.Split(decimal.NewFromFloat(vv.MinQty).String(), ".")[1])
	// 				}

	// 			}
	// 			minQuoue = v.QuotePrecision
	// 		}
	// 	}
	// 	fmt.Printf("minPrice:%d\n", minPrice)
	// 	fmt.Printf("minBase:%d\n", minBase)
	// 	fmt.Printf("minQuoue:%d\n", minQuoue)
	// }
	// return

	// 设置杠杆倍数
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// futureLeverageResponse, err := exch.FutureLeverage(symbolFuturePair.String(), 10)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("v:%+v\n", futureLeverageResponse)
	// }
	// return

	// 查看持仓风险
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// list, err := exch.PositionRiskFuture(symbolFuturePair.String())
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list {
	// 		fmt.Printf("v:%+v\n", v)
	// 	}
	// }
	// return

	// 下合约订单
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// oresp, err := exch.NewOrderFuture(&binance.FutureNewOrderRequest{
	// 	Symbol:       symbolFuturePair.String(),
	// 	Side:         order.Sell,
	// 	Type:         binance.BinanceRequestParamsOrderLimit,
	// 	PositionSide: binance.PositionSideSHORT,
	// 	TimeInForce:  binance.BinanceRequestParamsTimeGTC,
	// 	Quantity:     0.01,
	// 	Price:        1500.0,
	// })
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("oresp:%+v\n", oresp)
	// }
	// return

	// 查询合约订单
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// oresp, err := exch.QueryOrderFuture(symbolFuturePair.String(), 8389765490780171484, "")
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("oresp:%+v\n", oresp)
	// }
	// return

	// 查询打开的订单
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// list, err := exch.OpenOrdersFuture(symbolFuturePair.String())
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list {
	// 		fmt.Printf("OpenOrdersFuture:%+v", v)
	// 	}
	// }
	// return

	// 取消合约订单
	// symbolFuturePair := currency.NewPair(currency.NewCode("ETH"), currency.NewCode("USDT"))
	// symbolFuturePair.Delimiter = ""
	// oresp, err := exch.CancelExistingOrderFuture(symbolFuturePair.String(), 8389765490780171484, "")
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("oresp:%+v\n", oresp)
	// }
	// return

	// 获取合约K线
	// symbol, _ := currency.NewPairFromStrings("BTC", "USDT")
	// symbol.Delimiter = ""
	// startTime := time.Now().Add(-time.Minute * 20)
	// list, err := exch.GetHistoricCandlesFuture(symbol, binance.ContractTypePERPETUAL, startTime, time.Now(), kline.FiveMin)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list.Candles {
	// 		fmt.Printf("%+v\n", v)
	// 	}
	// }

	// 万向划转
	// tranid, err := exch.Transfer(binance.TransferType_MAIN_UMFUTURE, "USDT", 10)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("%+v\n", tranid)

	// }
	// 获取账户损益资金流水
	// todayTimeStat := "2021-01-11 20:40:10"
	// loc, _ := time.LoadLocation("Local") //重要：获取时区
	// timeStat, _ := time.ParseInLocation("2006-01-02 15:04:05", todayTimeStat, loc)
	// timeStatID := timeStat.UnixNano() / int64(time.Millisecond)
	// fmt.Println("")
	// // 1610381770000
	// // 1610380800000
	// list, err := exch.IncomeFuture(binance.FutureIncomeRequest{})
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list {
	// 		fmt.Printf("%+v\n", v)
	// 	}
	// }
	// return

	// 获取永续合约当前价格
	// symbol, _ := currency.NewPairFromStrings("BTC", "USDT")
	// symbol.Delimiter = ""
	// startTime := time.Now().Add(-time.Minute * 20)
	// list, err := exch.GetHistoricCandlesFuture(symbol, asset.PERPETUAL, startTime, time.Now(), kline.FiveMin)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list.Candles {
	// 		fmt.Printf("%+v\n", v)
	// 	}
	// }
	// return
	// 最新标记价格和资金费率
	// symbol, _ := currency.NewPairFromStrings("BTC", "USDT")
	// symbol.Delimiter = ""
	// list, err := exch.GetFuturePremiumIndex(symbol)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	fmt.Printf("%+v\n", list)

	// }

	// 最新标记价格和资金费率
	// symbol, _ := currency.NewPairFromStrings("BTC", "USDT")
	// symbol.Delimiter = ""
	// list, err := exch.GetFutureFundingRate(symbol, 0, 0, 0)
	// if err != nil {
	// 	panic(err)
	// } else {
	// 	for _, v := range list {
	// 		fmt.Printf("%+v\n", v)
	// 	}

	// }
	// return
	// 	//--------历史委托信息
	// 	// req := huobi.ContractHisordersRequest{
	// 	// 	Symbol:     "BTC",
	// 	// 	TradeType:  huobi.TradeType0,
	// 	// 	Type:       huobi.ContractHisOrderType1,
	// 	// 	Status:     huobi.ContractOrderStatus0,
	// 	// 	CreateDate: 90,
	// 	// }
	// 	// if res, err := huobiExch.GetContractHisorders(req); err != nil {
	// 	// 	panic(err)
	// 	// } else {
	// 	// 	fmt.Println("res", res)
	// 	// }

	// }

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
