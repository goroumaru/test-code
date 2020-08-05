package FTX

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/goroumaru/test-code/utils"
	"github.com/thrasher-corp/gocryptotrader/currency"
	"github.com/thrasher-corp/gocryptotrader/exchanges/asset"
	"github.com/thrasher-corp/gocryptotrader/exchanges/ftx"
	"github.com/thrasher-corp/gocryptotrader/exchanges/order"
	"gopkg.in/ini.v1"
)

var ApiKey, ApiSecret string

func init() {
	config, err := ini.Load(".env")
	if err != nil {
		fmt.Printf("Error: configLoad: %v\n", err)
	}
	ApiKey = config.Section("api").Key("API_KEY").MustString("")
	ApiSecret = config.Section("api").Key("API_SECRET").MustString("")
}

func TestGetFuture(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret
	//
	// Get Future
	//
	tickerNew, err := f.GetFuture("BTC-PERP")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", tickerNew)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			tickerNew, err := f.GetFuture("BTC-PERP")
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%+v\n", tickerNew)
		case <-ctx.Done():
			fmt.Printf("context done: %v\n", ctx.Err())
			return
		}
	}
}

func TestGetOrderbook(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret
	//
	// Get Orderbook
	//
	var depth int64 = 20
	ob, err := f.GetOrderbook("BTC-PERP", depth)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ob)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	tick := time.NewTicker(10 * time.Second)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			obNew, err := f.GetOrderbook("BTC-PERP", depth)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%+v\n", obNew)
		case <-ctx.Done():
			fmt.Printf("context done: %v\n", ctx.Err())
			return
		}
	}
}

func TestOrder(t *testing.T) {
	//
	// ※wrapperにバグ(SideにBid,Ask文字列が入ってしまう）ため、wrap前のモノを利用する
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// 現在価格取得
	pair := currency.NewPairWithDelimiter("BTC", "PERP", "-") // ペアがない場合、これで作成する
	assetType := asset.Futures
	ticker, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ticker)

	//
	// オープンオーダ
	//
	// NG!
	// openOrder := &order.Submit{
	// 	AssetType: assetType,
	// 	Pair:      pair,
	// 	Side:      order.Buy,
	// 	Type:      order.Limit,
	// 	Price:     10000, // $
	// 	Amount:    0.01,
	// }
	// resp, err := f.SubmitOrder(openOrder)

	// OK!
	type createOrder struct {
		MarketName string
		Side       string
		OrderType  string
		Price      float64
		Size       float64
		ReduceOnly string // option
		Ioc        string // option
		PostOnly   string // option
		ClientID   string // option
	}
	order := createOrder{
		MarketName: "BTC-PERP",
		Side:       "buy",
		OrderType:  "market", // "market"
		ReduceOnly: "true",   // can not use "false", use "". ポジションがないと、Reduceできない。
		// Price:      ticker.Last,
		Size: 0.001,
	}
	resp, err := f.Order(order.MarketName, order.Side, order.OrderType, order.ReduceOnly, "", "", "", order.Price, order.Size)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestGetOpenOrders(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	orders, err := f.GetOpenOrders("BTC-PERP") // {CreatedAt:2020-08-04 10:08:21.326848 +0000 +0000 FilledSize:0 Future:BTC-PERP ID:7251165376 Market:BTC-PERP Price:10285 AvgFillPrice:0 RemainingSize:0.001 Side:buy Size:0.001 Status:open OrderType:limit ReduceOnly:false IOC:false PostOnly:false ClientID:}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", orders)
}

func TestDeleteOrder(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// 現在価格取得
	pair := currency.NewPairWithDelimiter("BTC", "PERP", "-") // ペアがない場合、これで作成する
	assetType := asset.Futures
	ticker, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ticker)

	//
	// オープンオーダ
	//
	type createOrder struct {
		MarketName string
		Side       string
		OrderType  string
		Price      float64
		Size       float64
		ReduceOnly string // option
		Ioc        string // option
		PostOnly   string // option
		ClientID   string // option
	}
	order := createOrder{
		MarketName: "BTC-PERP",
		Side:       "buy",
		OrderType:  "limit", // "market"
		ReduceOnly: "",      // can not use "false", use "". ポジションがないと、Reduceできない。
		Price:      ticker.Last - 1000,
		Size:       0.001,
	}
	resp, err := f.Order(order.MarketName, order.Side, order.OrderType, order.ReduceOnly, "", "", "", order.Price, order.Size)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)

	// open order
	openOrder, err := f.GetOpenOrders("BTC-PERP") // {CreatedAt:2020-08-04 10:08:21.326848 +0000 +0000 FilledSize:0 Future:BTC-PERP ID:7251165376 Market:BTC-PERP Price:10285 AvgFillPrice:0 RemainingSize:0.001 Side:buy Size:0.001 Status:open OrderType:limit ReduceOnly:false IOC:false PostOnly:false ClientID:}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", openOrder)

	orderID := strconv.FormatInt(openOrder[0].ID, 10)
	delseteOrder, err := f.DeleteOrder(orderID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", delseteOrder)
}

func TestDeleteTriggerOrder(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// 現在価格取得
	pair := currency.NewPairWithDelimiter("BTC", "PERP", "-") // ペアがない場合、これで作成する
	assetType := asset.Futures
	ticker, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ticker)

	//
	// オープントリガーオーダ
	//
	type createTriggerOrder struct {
		MarketName   string
		Side         string
		OrderType    string  // "stop" or "trailingStop" or "takeProfit"
		ReduceOnly   string  // option
		OrderPrice   float64 // open order price
		TriggerPrice float64 // オーダが発動する価格
		TrailValue   float64 // option , トレイリングストップを使用するとき
		Size         float64
	}
	order := createTriggerOrder{
		MarketName:   "BTC-PERP",
		Side:         "buy",
		OrderType:    "trailingStop",
		ReduceOnly:   "",
		OrderPrice:   ticker.Last - 100,
		TriggerPrice: ticker.Last - 100,
		TrailValue:   10, // plus : buy , minus : sell
		Size:         0.001,
	}
	resp, err := f.TriggerOrder(order.MarketName, order.Side, order.OrderType, order.ReduceOnly, "",
		order.Size, order.TriggerPrice, order.OrderPrice, order.TrailValue)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)

	// オープントリガーオーダ確認
	openOrder, err := f.GetOpenTriggerOrders("BTC-PERP", "trailingStop") //
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", openOrder)

	// オープントリガーオーダ削除
	orderID := strconv.FormatInt(openOrder[0].ID, 10)
	deleteOrder, err := f.DeleteTriggerOrder(orderID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", deleteOrder)
}

func TestTriggerOrder(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// 現在価格取得
	pair := currency.NewPairWithDelimiter("BTC", "PERP", "-") // ペアがない場合、これで作成する
	assetType := asset.Futures
	ticker, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ticker)

	//
	// オープンオーダ
	//
	type createTriggerOrder struct {
		MarketName   string
		Side         string
		OrderType    string  // "stop" or "trailingStop" or "takeProfit"
		ReduceOnly   string  // option
		OrderPrice   float64 // open order price
		TriggerPrice float64 // オーダが発動する価格
		TrailValue   float64 // option , トレイリングストップを使用するとき
		Size         float64
	}
	order := createTriggerOrder{
		MarketName:   "BTC-PERP",
		Side:         "buy",
		OrderType:    "trailingStop",
		ReduceOnly:   "",
		OrderPrice:   ticker.Last - 100,
		TriggerPrice: ticker.Last - 100,
		TrailValue:   10, // plus : buy , minus : sell
		Size:         0.001,
	}
	resp, err := f.TriggerOrder(order.MarketName, order.Side, order.OrderType, order.ReduceOnly, "",
		order.Size, order.TriggerPrice, order.OrderPrice, order.TrailValue)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestGetPositions(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// ポジション取得
	positions, err := f.GetPositions() // 実際にオーダしてみたが、正しいメンバーに値が入ってない
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", positions)
}

func TestBalances(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// ポジション取得
	balances, err := f.GetBalances() //[{Coin:USD Free:66.70684172 Total:96.11855811} {Coin:BTC Free:0 Total:0}]
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", balances)
}

//
// FTX wrapper
//
func TestFetchTicker(t *testing.T) {
	//
	// FTX wrapper
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx_wrapper.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret
	// pair
	pair := currency.NewPair(currency.BTC, currency.USD)
	// assetType
	assetType := asset.Spot
	//
	// Fetch ticker
	//
	tickerNew, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", tickerNew)
}

func TestFetchOrderBook(t *testing.T) {
	//
	// FTX wrapper
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx_wrapper.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret
	// pair
	pair := currency.NewPair(currency.BTC, currency.USD)
	// assetType
	assetType := asset.Spot
	//
	// Fetch Orderbook(depth)
	//
	ob, err := f.FetchOrderbook(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ob)
}

func TestFetchAccountInfo(t *testing.T) {
	//
	// FTX wrapper
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx_wrapper.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret
	//
	// アカウント情報取得（メインアカウントのみ）
	//
	holdings, err := f.FetchAccountInfo()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", holdings.Accounts) // [{ID: Currencies:[{CurrencyName:USD TotalValue:46.53845886 Hold:29.51270195} {CurrencyName:BTC TotalValue:0 Hold:0}]}]
	fmt.Printf("%+v\n", holdings.Exchange) // FTX
}

func TestGetAllWalletBalances(t *testing.T) {
	//
	// FTX wrapper
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx_wrapper.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret
	//
	// 資産合計取得（メインアカウントのみ）
	//
	allWallet, err := f.GetAllWalletBalances()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", allWallet) // {Main:[{Coin:USD Free:17.02575691 Total:46.53845886} {Coin:BTC Free:0 Total:0}] BattleRoyale:[]}
}

func TestGetOpenOrdersWrpper(t *testing.T) {
	//
	// FTX wrapper
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx_wrapper.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret
	//
	// オープンオーダ取得(spot = 現物)
	//
	marketName := "BTC-PERP"
	fmt.Println(marketName)
	openOrders, err := f.GetOpenOrders(marketName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", openOrders)
	//
	// オープンオーダ取得(futures = 先物、デリバティブ)
	//
	now := time.Now().UTC() // UTCでmarketNameの日付が変わることに注意する
	_, month, day := utils.ConvertFormatDate(now)

	marketName = "BTC-MOVE-" + month + day // "BTC-MOVE-0705"
	fmt.Println(marketName)
	openOrders, err = f.GetOpenOrders(marketName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", openOrders)
}

func TestCancelOrder(t *testing.T) {
	//
	// FTX wrapper
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx_wrapper.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// FTXのペア確認
	markets, _ := f.GetMarkets()
	fmt.Println(markets)

	// 現在価格取得
	pair := currency.NewPairWithDelimiter("BTC", "PERP", "-") // ペアがない場合、これで作成する
	assetType := asset.Futures
	ticker, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ticker)

	// キャンセルするオーダ生成
	type createOrder struct {
		MarketName string
		Side       string
		OrderType  string
		Price      float64
		Size       float64
		ReduceOnly string // option
		Ioc        string // option
		PostOnly   string // option
		ClientID   string // option
	}
	openOrder := createOrder{
		MarketName: "BTC-PERP",
		Side:       "buy",
		OrderType:  "limit",
		Price:      ticker.Last - 1000, // 現在価格とかけ離れていると注文が通らないため
		Size:       0.01,
	}
	resp, err := f.Order(openOrder.MarketName, openOrder.Side, openOrder.OrderType, "", "", "", "", openOrder.Price, openOrder.Size)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)

	//
	// オーダキャンセル(orderID指定)
	//
	cancelOrder := &order.Cancel{
		ID: strconv.FormatInt(resp.ID, 10),
	}
	if err := f.CancelOrder(cancelOrder); err != nil { // RESTリクエストがDELETEではなく、GETとなっているバグあり
		fmt.Println(err)
	}
}

func TestCancelAllOrders(t *testing.T) {
	//
	// FTX wrapper
	// https://github.com/thrasher-corp/gocryptotrader/blob/master/exchanges/ftx/ftx_wrapper.go
	//
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// FTXのペア確認
	markets, _ := f.GetMarkets()
	fmt.Println(markets)

	// 現在価格取得
	pair := currency.NewPairWithDelimiter("BTC", "PERP", "-") // ペアがない場合、これで作成する
	assetType := asset.Futures
	ticker, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ticker)

	// キャンセルするオーダ生成
	type createOrder struct {
		MarketName string
		Side       string
		OrderType  string
		Price      float64
		Size       float64
		ReduceOnly string // option
		Ioc        string // option
		PostOnly   string // option
		ClientID   string // option
	}
	openOrder := createOrder{
		MarketName: "BTC-PERP",
		Side:       "buy",
		OrderType:  "limit",
		Price:      ticker.Last - 1000, // 現在価格とかけ離れていると注文が通らないため
		Size:       0.01,
	}
	respOpenOrder, err := f.Order(openOrder.MarketName, openOrder.Side, openOrder.OrderType, "", "", "", "", openOrder.Price, openOrder.Size)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", respOpenOrder)

	//
	// オーダキャンセル(すべて)
	//
	cancelOrder := &order.Cancel{
		Pair:      pair,
		AssetType: assetType,
	}
	respCancelOrder, err := f.CancelAllOrders(cancelOrder)
	if err != nil { // RESTリクエストがDELETEではなく、GETとなっているバグあり
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", respCancelOrder)
}

func TestDeleteTriggerOrderWrapper(t *testing.T) {
	var f ftx.FTX
	f.SetDefaults()
	// API Key
	f.API.Credentials.Key = ApiKey
	f.API.Credentials.Secret = ApiSecret

	// FTXのペア確認
	markets, _ := f.GetMarkets()
	fmt.Println(markets)

	// 現在価格取得
	pair := currency.NewPairWithDelimiter("BTC", "PERP", "-") // ペアがない場合、これで作成する
	assetType := asset.Futures
	ticker, err := f.FetchTicker(pair, assetType)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", ticker)

	// キャンセルするオーダ生成
	type createOrder struct {
		MarketName string
		Side       string
		OrderType  string
		Price      float64
		Size       float64
		ReduceOnly string // option
		Ioc        string // option
		PostOnly   string // option
		ClientID   string // option
	}
	openOrder := createOrder{
		MarketName: "BTC-PERP",
		Side:       "buy",
		OrderType:  "limit",
		Price:      ticker.Last - 1000, // 現在価格とかけ離れていると注文が通らないため
		Size:       0.01,
	}
	resp, err := f.Order(openOrder.MarketName, openOrder.Side, openOrder.OrderType, "", "", "", "", openOrder.Price, openOrder.Size)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", resp)

	//
	// オーダキャンセル(orderID指定)
	//
	cancelOrder := &order.Cancel{
		ID: strconv.FormatInt(resp.ID, 10),
	}
	respDeleteOrder, err := f.DeleteTriggerOrder(cancelOrder.ID)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v\n", respDeleteOrder)
}
