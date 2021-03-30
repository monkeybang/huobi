package huobi

import (
	"encoding/json"
	"errors"
	"github.com/spf13/cast"
	"github.com/tidwall/gjson"
	"log"
	"math"
	"time"
)

type Exchange struct {
	name      string
	accountId string
	accessKey string
	secretKey string
	symbols   map[string]*SymbolsData
}

func NewExchange(ak, sk string) *Exchange {
	ex := &Exchange{
		accessKey: ak,
		secretKey: sk,
	}
	ex.name = "huobi"
	ex.symbols = ex.GetSymbols()
	accounts := ex.GetAccounts()
	for _, data := range accounts.Data {
		if data.Type == `spot` {
			ex.accountId = cast.ToString(data.ID)
		}
	}
	return ex
}

func (huobi *Exchange) GetSymbols() map[string]*SymbolsData {
	symbolMap := make(map[string]*SymbolsData)
	symbolsReturn := GetSymbols()
	if symbolsReturn.Status == "ok" {
		for i := range symbolsReturn.Data {
			symbolMap[symbolsReturn.Data[i].BaseCurrency+symbolsReturn.Data[i].QuoteCurrency] = symbolsReturn.Data[i]
		}
	}
	return symbolMap
}

func (huobi *Exchange) TruncPrice(symbol string, price float64) (float64, bool) {
	if data, ok := huobi.symbols[symbol]; ok == true {
		pre := math.Pow10(data.PricePrecision)
		tPrice := math.Round(price*pre) / pre
		return tPrice, true
	}
	return 0, false
}

func (huobi *Exchange) TruncAmount(symbol string, amount float64) (float64, bool) {
	if data, ok := huobi.symbols[symbol]; ok == true {
		pre := math.Pow10(data.AmountPrecision)
		amount := math.Trunc(amount*pre) / pre
		return amount, true
	}
	return 0, false
}

func (huobi *Exchange) Trunc(symbol string, price float64, amount float64) (float64, float64) {
	truncPrice, err1 := huobi.TruncPrice(symbol, price)
	truncAmount, err2 := huobi.TruncAmount(symbol, amount)
	if err1 == false || err2 == false {
		log.Fatal(err1, err2)
	}
	return truncPrice, truncAmount
}

func (huobi *Exchange) BuyLimitEver(symbol string, amount float64, price float64) (string, error) {
	placeParams := &PlaceRequestParams{}
	placeParams.AccountID = huobi.accountId
	placeParams.Amount = cast.ToString(amount)
	placeParams.Price = cast.ToString(price)
	placeParams.Source = "api"
	placeParams.Symbol = symbol
	placeParams.Type = "buy-limit"

	times := 20
	for times > 0 {
		times--
		placeReturn := huobi.Place(placeParams)
		if placeReturn.Status == "ok" {
			//log.Println("Place return:", placeReturn.Data)
			return placeReturn.Data, nil
		} else {
			log.Println("place error:", placeReturn.ErrMsg, amount, price)
			time.Sleep(time.Millisecond * 100)
		}
	}
	return "", errors.New("buy failed")
}

func (huobi *Exchange) SellLimitEver(symbol string, amount float64, price float64) (string, error) {
	placeParams := &PlaceRequestParams{}
	placeParams.AccountID = huobi.accountId
	placeParams.Amount = cast.ToString(amount)
	placeParams.Price = cast.ToString(price)
	placeParams.Source = "api"
	placeParams.Symbol = symbol
	placeParams.Type = "sell-limit"

	times := 20
	for times > 0 {
		times--
		placeReturn := huobi.Place(placeParams)
		if placeReturn.Status == "ok" {
			//log.Println("Place return:", placeReturn.Data)
			return placeReturn.Data, nil
		} else {
			log.Println("place error:", placeReturn.ErrMsg, amount, price, symbol)
			time.Sleep(time.Millisecond * 100)
		}
	}
	return "", errors.New("buy failed")
}

func (huobi *Exchange) PlaceOrder(symbol string, orderType string, amount float64, price float64) {
	placeParams := &PlaceRequestParams{}
	placeParams.AccountID = huobi.accountId
	placeParams.Amount = cast.ToString(amount)
	placeParams.Price = cast.ToString(price)
	placeParams.Source = "api"
	placeParams.Symbol = symbol
	placeParams.Type = orderType

	placeReturn := huobi.Place(placeParams)
	if placeReturn.Status == "ok" {
		log.Print("Place return :", placeReturn.Data)
	} else {
		log.Println("place error:", placeReturn.ErrMsg, placeReturn, symbol, orderType, amount, price)
	}
}

func (ex *Exchange) BatchCancelOrders(symbol string) {
	params := make(map[string]string)
	params["account-id"] = ex.accountId
	params["symbol"] = symbol

	strRequest := "/v1/order/orders/batchCancelOpenOrders"
	jsonPlaceReturn := ex.ApiKeyPost(make(map[string]string), strRequest)
	log.Print(jsonPlaceReturn)
}

func (huobi *Exchange) GetAccountId() string {
	return huobi.accountId
}

func (ex *Exchange) OpenOrders(symbol string) *OrderReturn {
	params := make(map[string]string)
	params["account-id"] = ex.accountId
	params["symbol"] = symbol
	params["size"] = "500"

	strRequest := "/v1/order/openOrders"
	str := ex.ApiKeyGet(make(map[string]string), strRequest)

	orderReturn := &OrderReturn{}

	err := json.Unmarshal([]byte(str), orderReturn)
	if err != nil {
		log.Println(str, err)
	}
	return orderReturn
}

func (ex *Exchange) GetOrder(orderId string) *Order {
	params := make(map[string]string)

	strRequest := "/v1/order/orders/" + cast.ToString(orderId)
	str := ex.ApiKeyGet(params, strRequest)
	//log.Println(str)

	orderReturnSingle := &OrderReturnSingle{}
	err := json.Unmarshal([]byte(str), orderReturnSingle)
	if err != nil {
		log.Println(str, err)
		return nil
	}
	return &orderReturnSingle.Data
}

func (ex *Exchange) CancelOrder(orderId string) string {
	params := make(map[string]string)

	strRequest := "/v1/order/orders/" + orderId + "/submitcancel"
	str := ex.ApiKeyPost(params, strRequest)
	id := gjson.Get(str, "data").String()
	return id
}

func (ex *Exchange) EtpRedemption(symbol, usdt string, amount float64) string {
	mapParams := make(map[string]string)
	mapParams["etpName"] = symbol
	mapParams["currency"] = usdt
	mapParams["amount"] = cast.ToString(amount)
	strRequestUrl := "/v2/etp/redemption"
	jsonMarketDetailReturn := ex.ApiKeyPost(mapParams, strRequestUrl)
	return jsonMarketDetailReturn
}

func (ex *Exchange) getEtpTransaction(id string) string {
	mapParams := make(map[string]string)
	mapParams["transactId"] = id
	strRequestUrl := "/v2/etp/transaction"
	jsonMarketDetailReturn := ex.ApiKeyGet(mapParams, strRequestUrl)
	return jsonMarketDetailReturn
}

func (ex *Exchange) GetAggregateBalance() *Aggregate {
	mapParams := make(map[string]string)
	strRequestUrl := "/v1/subuser/aggregate-balance"
	resp := ex.ApiKeyGet(mapParams, strRequestUrl)
	//log.Println(resp)
	aggregate := &Aggregate{}
	err := json.Unmarshal([]byte(resp), aggregate)
	if err != nil {
		log.Println(err, resp)
		return nil
	}
	if gjson.Get(resp, "status").String() == "error" {
		log.Println(resp)
	}
	return aggregate
}

func (ex *Exchange) GetContractPositionInfo() *ContractAggregate {
	mapParam := make(map[string]string)
	strUrl := "/api/v1/contract_sub_account_list"
	resp := ex.ContractKeyPost(mapParam, strUrl)
	agg := &ContractAggregate{}
	err := json.Unmarshal([]byte(resp), agg)
	if err != nil {
		log.Println(err)
		return nil
	}
	return agg
}

func (ex *Exchange) GetSwapPositionInfo() *ContractAggregate {
	mapParam := make(map[string]string)
	strUrl := "/swap-api/v1/swap_sub_account_list"
	resp := ex.ContractKeyPost(mapParam, strUrl)
	agg := &ContractAggregate{}
	err := json.Unmarshal([]byte(resp), agg)
	if err != nil {
		log.Println(err)
		return nil
	}
	return agg
}

func (ex *Exchange) GetLinearSwapPositionInfo() *ContractAggregate {
	mapParam := make(map[string]string)
	strUrl := "/linear-swap-api/v1/swap_sub_account_list"
	resp := ex.ContractKeyPost(mapParam, strUrl)
	agg := &ContractAggregate{}
	err := json.Unmarshal([]byte(resp), agg)
	if err != nil {
		log.Println(err)
		return nil
	}
	return agg
}

func (ex *Exchange) GetLinearSwapCrossPositionInfo() *ContractAggregate {
	mapParam := make(map[string]string)
	strUrl := "/linear-swap-api/v1/swap_cross_sub_account_list"
	resp := ex.ContractKeyPost(mapParam, strUrl)
	agg := &ContractAggregate{}
	err := json.Unmarshal([]byte(resp), agg)
	if err != nil {
		log.Println(err, resp)
		return nil
	}
	return agg
}

func (ex *Exchange) GetOptionPositionInfo() *ContractAggregate {
	mapParam := make(map[string]string)
	strUrl := "/option-api/v1/option_sub_account_list"
	resp := ex.ContractKeyPost(mapParam, strUrl)
	log.Println(resp)
	agg := &ContractAggregate{}
	err := json.Unmarshal([]byte(resp), agg)
	if err != nil {
		log.Println(err)
		return nil
	}
	return agg
}

func (ex *Exchange) GetAssetValuation(accountType string, valuationCurrency string) (string, error) {
	mapParams := make(map[string]string)
	mapParams["accountType"] = accountType
	mapParams["valuationCurrency"] = valuationCurrency
	strRequestUrl := "/v2/account/asset-valuation"
	resp := ex.ApiKeyGet(mapParams, strRequestUrl)
	if gjson.Get(resp, "code").Int() == 200 {
		return resp, nil
	}
	return resp, errors.New(resp)
}

func (ex *Exchange) GetContractBalanceValuation(valuation_asset string) (string, error) {
	mapParams := make(map[string]string)
	mapParams["valuation_asset"] = valuation_asset
	strUrl := "/api/v1/contract_balance_valuation"
	resp := ex.ContractKeyPost(mapParams, strUrl)
	if gjson.Get(resp, "status").String() == "ok" {
		return resp, nil
	}
	return resp, errors.New(resp)
}

func (ex *Exchange) GetSwapBalanceValuation(valuation_asset string) (string, error) {
	mapParams := make(map[string]string)
	mapParams["valuation_asset"] = valuation_asset
	strUrl := "/swap-api/v1/swap_balance_valuation"
	resp := ex.ContractKeyPost(mapParams, strUrl)
	if gjson.Get(resp, "status").String() == "ok" {
		return resp, nil
	}
	return resp, errors.New(resp)
}

func (ex *Exchange) GetLinearSwapBalanceValuation(valuation_asset string) (string, error) {
	mapParams := make(map[string]string)
	mapParams["valuation_asset"] = valuation_asset
	strUrl := "/linear-swap-api/v1/swap_balance_valuation"
	resp := ex.ContractKeyPost(mapParams, strUrl)
	if gjson.Get(resp, "status").String() == "ok" {
		return resp, nil
	}
	return resp, errors.New(resp)
}