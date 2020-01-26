package huobi

import (
	"github.com/spf13/cast"
	"log"
	"math"
	"time"
)

var accessKey string
var secretKey string
var accountId string

type Exchange struct {
	name    string
	symbols map[string]*SymbolsData
}

func NewExchange(ak, sk string) *Exchange {
	accessKey = ak
	secretKey = sk
	accountId = accessKey
	huobi := &Exchange{}
	huobi.name = "huobi"
	huobi.symbols = huobi.GetSymbols()
	return huobi
}

func (huobi *Exchange) GetSymbols() map[string]*SymbolsData {
	symbolMap := make(map[string]*SymbolsData)
	symbolsReturn := GetSymbols()
	if symbolsReturn.Status != "ok" {
		for i := range symbolsReturn.Data {
			symbolMap[symbolsReturn.Data[i].SymbolPartition] = symbolsReturn.Data[i]
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

func (huobi *Exchange) BuyLimitEver(symbol string, amount float64, price float64) {
	placeParams := &PlaceRequestParams{}
	placeParams.AccountID = accountId
	placeParams.Amount = cast.ToString(amount)
	placeParams.Price = cast.ToString(price)
	placeParams.Source = "api"
	placeParams.Symbol = symbol
	placeParams.Type = "buy-limit"

	for {
		placeReturn := Place(placeParams)
		if placeReturn.Status == "ok" {
			log.Println("Place return:", placeReturn.Data)
			break
		} else {
			log.Println("place error:", placeReturn.ErrMsg)
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (huobi *Exchange) SellLimitEver(symbol string, amount float64, price float64) {
	placeParams := &PlaceRequestParams{}
	placeParams.AccountID = accountId
	placeParams.Amount = cast.ToString(amount)
	placeParams.Price = cast.ToString(price)
	placeParams.Source = "api"
	placeParams.Symbol = symbol
	placeParams.Type = "sell-limit"

	for {
		placeReturn := Place(placeParams)
		if placeReturn.Status == "ok" {
			log.Println("Place return:", placeReturn.Data)
			break
		} else {
			log.Println("place error:", placeReturn.ErrMsg)
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (huobi *Exchange) PlaceOrder(symbol string, orderType string, amount float64, price float64) {
	placeParams := &PlaceRequestParams{}
	placeParams.AccountID = accountId
	placeParams.Amount = cast.ToString(amount)
	placeParams.Price = cast.ToString(price)
	placeParams.Source = "api"
	placeParams.Symbol = symbol
	placeParams.Type = orderType

	placeReturn := Place(placeParams)
	if placeReturn.Status == "ok" {
		log.Print("Place return :", placeReturn.Data)
	} else {
		log.Println("place error:", placeReturn.ErrMsg, placeReturn, symbol, orderType, amount, price)
	}
}

func (huobi *Exchange) BatchCancelOrders(symbol string) {
	params := make(map[string]string)
	params["account-id"] = accountId
	params["symbol"] = symbol

	strRequest := "/v1/order/orders/batchCancelOpenOrders"
	jsonPlaceReturn := ApiKeyPost(make(map[string]string), strRequest)
	log.Print(jsonPlaceReturn)
}
