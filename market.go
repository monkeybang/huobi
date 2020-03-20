package huobi

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

var MARKET_URL = `https://api.huobi.pro`
var TRADE_URL = `https://api.huobi.pro`
var HOST_NAME = `api.huobi.pro`

// 获取聚合行情
// strSymbol: 交易对, btcusdt, bccbtc......
// return: TickReturn对象
func GetTicker(strSymbol string) *TickerReturn {
	tickerReturn := &TickerReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail/merged"
	strUrl := MARKET_URL + strRequestUrl

	jsonTickReturn := HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonTickReturn), tickerReturn)
	if err != nil {
		log.Print(err)
		return nil
	}
	return tickerReturn
}

// 获取交易深度信息
// strSymbol: 交易对, btcusdt, bccbtc......
// strType: Depth类型, step0、step1......stpe5 (合并深度0-5, 0时不合并)
// return: MarketDepthReturn对象
func GetMarketDepth(strSymbol, strType string) *MarketDepthReturn {
	marketDepthReturn := &MarketDepthReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["type"] = strType

	strRequestUrl := "/market/depth"
	strUrl := MARKET_URL + strRequestUrl

	jsonMarketDepthReturn := HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonMarketDepthReturn), &marketDepthReturn)
	if err != nil {
		log.Println(err)
	}
	return marketDepthReturn
}

// 获取交易细节信息
// strSymbol: 交易对, btcusdt, bccbtc......
// return: TradeDetailReturn对象
func GetTradeDetail(strSymbol string) *TradeDetailReturn {
	tradeDetailReturn := &TradeDetailReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/trade"
	strUrl := MARKET_URL + strRequestUrl

	jsonTradeDetailReturn := HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonTradeDetailReturn), &tradeDetailReturn)
	if err != nil {
		log.Println(err)
	}
	return tradeDetailReturn
}

// 批量获取最近的交易记录
// strSymbol: 交易对, btcusdt, bccbtc......
// nSize: 获取交易记录的数量, 范围1-2000
// return: TradeReturn对象
func GetTrade(strSymbol string, nSize int) *TradeReturn {
	tradeReturn := &TradeReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/trade"
	strUrl := MARKET_URL + strRequestUrl

	jsonTradeReturn := HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonTradeReturn), &tradeReturn)
	if err != nil {
		log.Println(err)
	}
	return tradeReturn
}

// 获取Market Detail 24小时成交量数据
// strSymbol: 交易对, btcusdt, bccbtc......
// return: MarketDetailReturn对象
func GetMarketDetail(strSymbol string) *MarketDetailReturn {
	marketDetailReturn := &MarketDetailReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail"
	strUrl := MARKET_URL + strRequestUrl

	jsonMarketDetailReturn := HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonMarketDetailReturn), &marketDetailReturn)
	if err != nil {
		log.Println(err)
	}
	return marketDetailReturn
}

func GetKline(strSymbol string) *KLineReturn {
	kLineReturn := &KLineReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/history/kline"
	strUrl := MARKET_URL + strRequestUrl

	jsonMarketKlineReturn := HttpGetRequest(strUrl, mapParams)
	err := json.Unmarshal([]byte(jsonMarketKlineReturn), &kLineReturn)
	if err != nil {
		log.Println(err)
	}
	return kLineReturn
}

//----------------------------------------
// 公共API

// 查询系统支持的所有交易及精度
// return: SymbolsReturn对象
func GetSymbols() *SymbolsReturn {
	symbolsReturn := &SymbolsReturn{}

	strRequestUrl := "/v1/common/symbols"
	strUrl := TRADE_URL + strRequestUrl

	jsonSymbolsReturn := HttpGetRequest(strUrl, nil)
	err := json.Unmarshal([]byte(jsonSymbolsReturn), &symbolsReturn)
	if err != nil {
		log.Println(err)
	}
	return symbolsReturn
}

// 查询系统支持的所有币种
// return: CurrencysReturn对象
func GetCurrencys() *CurrencysReturn {
	currencysReturn := &CurrencysReturn{}

	strRequestUrl := "/v1/common/currencys"
	strUrl := TRADE_URL + strRequestUrl

	jsonCurrencysReturn := HttpGetRequest(strUrl, nil)
	err := json.Unmarshal([]byte(jsonCurrencysReturn), &currencysReturn)
	if err != nil {
		log.Println(err)
	}
	return currencysReturn
}

// 查询系统当前时间戳
// return: TimestampReturn对象
func GetTimestamp() *TimestampReturn {
	timestampReturn := &TimestampReturn{}

	strRequest := "/v1/common/timestamp"
	strUrl := TRADE_URL + strRequest

	jsonTimestampReturn := HttpGetRequest(strUrl, nil)
	err := json.Unmarshal([]byte(jsonTimestampReturn), &timestampReturn)
	if err != nil {
		log.Println(err)
	}
	return timestampReturn
}

//------------------------------------------------------------------------------------------
// 用户资产API

// 查询当前用户的所有账户, 根据包含的私钥查询
// return: AccountsReturn对象
func GetAccounts() *AccountsReturn {
	accountsReturn := &AccountsReturn{}

	strRequest := "/v1/account/accounts"
	jsonAccountsReturn := ApiKeyGet(make(map[string]string), strRequest)
	err := json.Unmarshal([]byte(jsonAccountsReturn), &accountsReturn)
	if err != nil {
		log.Println(err)
	}
	return accountsReturn
}

// 根据账户ID查询账户余额
// nAccountID: 账户ID, 不知道的话可以通过GetAccounts()获取, 可以只现货账户, C2C账户, 期货账户
// return: BalanceReturn对象
func GetAccountBalance(strAccountID string) *BalanceReturn {
	balanceReturn := &BalanceReturn{}

	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", strAccountID)
	jsonBanlanceReturn := ApiKeyGet(make(map[string]string), strRequest)
	err := json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)
	if err != nil {
		log.Println(err)
	}
	return balanceReturn
}

//------------------------------------------------------------------------------------------
// 交易API

// 下单
// placeRequestParams: 下单信息
// return: PlaceReturn对象
func Place(placeRequestParams *PlaceRequestParams) *PlaceReturn {
	placeReturn := &PlaceReturn{}

	mapParams := make(map[string]string)
	mapParams["account-id"] = placeRequestParams.AccountID
	mapParams["amount"] = placeRequestParams.Amount
	if 0 < len(placeRequestParams.Price) {
		mapParams["price"] = placeRequestParams.Price
	}
	if 0 < len(placeRequestParams.Source) {
		mapParams["source"] = placeRequestParams.Source
	}
	mapParams["symbol"] = placeRequestParams.Symbol
	mapParams["type"] = placeRequestParams.Type

	strRequest := "/v1/order/orders/place"
	jsonPlaceReturn := ApiKeyPost(mapParams, strRequest)
	err := json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)
	if err != nil {
		log.Println(err)
	}
	return placeReturn
}

// 申请撤销一个订单请求
// strOrderID: 订单ID
// return: PlaceReturn对象
func SubmitCancel(strOrderID string) *PlaceReturn {
	placeReturn := &PlaceReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", strOrderID)
	jsonPlaceReturn := ApiKeyPost(make(map[string]string), strRequest)
	err := json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)
	if err != nil {
		log.Println(err)
	}
	return placeReturn
}
