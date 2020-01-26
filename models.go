package huobi

import (
	"encoding/json"
	"log"
)

// 子账户结构
type SubAccount struct {
	Currency string `json:"currency"` // 币种
	Balance  string `json:"balance"`  // 结余
	Type     string `json:"type"`     // 类型, trade: 交易余额, frozen: 冻结余额
}

type Balance struct {
	ID     int64        `json:"id"`    // 账户ID
	State  string       `json:"state"` // 账户状态, working: 正常, lock: 账户被锁定
	Type   string       `json:"type"`  // 账户类型, spot: 现货账户
	List   []SubAccount `json:"list"`  // 子账户数组
	UserID int64        `json:"user-id"`
}

type BalanceReturn struct {
	Status  string  `json:"status"` // 请求状态
	Data    Balance `json:"data"`   // 账户余额
	ErrCode string  `json:"err-code"`
	ErrMsg  string  `json:"err-msg"`
}

type AccountsData struct {
	ID     int64  `json:"id"`      // Account ID
	Type   string `json:"type"`    // 账户类型, spot: 现货账户
	State  string `json:"state"`   // 账户状态, working: 正常, lock: 账户被锁定
	UserID int64  `json:"user-id"` // 用户ID
}

type AccountsReturn struct {
	Status  string         `json:"status"` // 请求状态
	Data    []AccountsData `json:"data"`   // 用户数据
	ErrCode string         `json:"err-code"`
	ErrMsg  string         `json:"err-msg"`
}

type CurrencysReturn struct {
	Status  string   `json:"status"` // 请求状态
	Data    []string `json:"data"`   // 系统支持的所有币种
	ErrCode string   `json:"err-code"`
	ErrMsg  string   `json:"err-msg"`
}

type KLineData struct {
	ID     int64   `json:"id"`     // K线ID
	Amount float64 `json:"amount"` // 成交量
	Count  int64   `json:"count"`  // 成交笔数
	Open   float64 `json:"open"`   // 开盘价
	Close  float64 `json:"close"`  // 收盘价, 当K线为最晚的一根时, 时最新成交价
	Low    float64 `json:"low"`    // 最低价
	High   float64 `json:"high"`   // 最高价
	Vol    float64 `json:"vol"`    // 成交额, 即SUM(每一笔成交价 * 该笔的成交数量)
}

type KLineReturn struct {
	Status  string      `json:"status"`   // 请求处理结果, "ok"、"error"
	Ts      int64       `json:"ts"`       // 响应生成时间点, 单位毫秒
	Data    []KLineData `json:"data"`     // KLine数据
	Ch      string      `json:"ch"`       // 数据所属的Channel, 格式: market.$symbol.kline.$period
	ErrCode string      `json:"err-code"` // 错误代码
	ErrMsg  string      `json:"err-msg"`  // 错误提示
}

type MarketDepth struct {
	ID   int64       `json:"id,omitempty"` // 消息ID
	Ts   int64       `json:"ts,omitempty"` // 消息声称事件, 单位: 毫秒
	Bids [][]float64 `json:"bids"`         // 买盘, [price(成交价), amount(成交量)], 按price降序排列
	Asks [][]float64 `json:"asks"`         // 卖盘, [price(成交价), amount(成交量)], 按price升序排列
}

type MarketDepthReturn struct {
	Status  string      `json:"status"` // 请求状态, ok或者error
	Ts      int64       `json:"ts"`     // 响应生成时间点, 单位: 毫秒
	Tick    MarketDepth `json:"tick"`   // Depth数据
	Ch      string      `json:"ch"`     //  数据所属的Channel, 格式: market.$symbol.depth.$type
	ErrCode string      `json:"err-code,omitempty"`
	ErrMsg  string      `json:"err-msg,omitempty"`
}

type MarketDetail struct {
	ID     int64   `json:"id"`     // 消息ID
	Ts     int64   `json:"ts"`     // 24小时统计时间
	Amount float64 `json:"amount"` // 24小时成交量
	Open   float64 `json:"open"`   // 前24小时成交价
	Close  float64 `json:"close"`  // 当前成交价
	High   float64 `json:"high"`   // 近24小时最高价
	Low    float64 `json:"low"`    // 近24小时最低价
	Count  int64   `json:"count"`  // 近24小时累计成交数
	Vol    float64 `json:"vol"`    // 近24小时累计成交额, 即SUM(每一笔成交价 * 该笔的成交量)
}

type MarketDetailReturn struct {
	Status  string       `json:"status"` // 请求状态
	Ts      int64        `json:"ts"`     // 响应生成时间点
	Tick    MarketDetail `json:"tick"`   // Market Detail 24小时成交量数据
	Ch      string       `json:"ch"`     // 数据所属的Channel, 格式: market.$symbol.depth.$type
	ErrCode string       `json:"err-code"`
	ErrMsg  string       `json:"err-msg"`
}

type PlaceRequestParams struct {
	AccountID string `json:"account-id"` // 账户ID
	Amount    string `json:"amount"`     // 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
	Price     string `json:"price"`      // 下单价格, 市价单不传该参数
	Source    string `json:"source"`     // 订单来源, api: API调用, margin-api: 借贷资产交易
	Symbol    string `json:"symbol"`     // 交易对, btcusdt, bccbtc......
	Type      string `json:"type"`       // 订单类型, buy-market: 市价买, sell-market: 市价卖, buy-limit: 限价买, sell-limit: 限价卖
}

func (place *PlaceRequestParams) String() string {
	bytes, err := json.Marshal(place)
	if err != nil {
		log.Fatalln(place)
	}
	return string(bytes)
}

type PlaceReturn struct {
	Status  string `json:"status"`
	Data    string `json:"data"`
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type SymbolsData struct {
	BaseCurrency    string `json:"base-currency"`    // 基础币种
	QuoteCurrency   string `json:"quote-currency"`   // 计价币种
	PricePrecision  int    `json:"price-precision"`  // 价格精度位数(0为个位)
	AmountPrecision int    `json:"amount-precision"` // 数量精度位数(0为个位)
	SymbolPartition string `json:"symbol-partition"` // 交易区, main: 主区, innovation: 创新区, bifurcation: 分叉区
}

type SymbolsReturn struct {
	Status  string         `json:"status"` // 请求状态
	Data    []*SymbolsData `json:"data"`   // 交易及精度数据
	ErrCode string         `json:"err-code"`
	ErrMsg  string         `json:"err-msg"`
}

type Ticker struct {
	ID     int64     `json:"id"`     // K线ID
	Amount float64   `json:"amount"` // 成交量
	Count  int64     `json:"count"`  // 成交笔数
	Open   float64   `json:"open"`   // 开盘价
	Close  float64   `json:"close"`  // 收盘价
	Low    float64   `json:"low"`    // 最低价
	High   float64   `json:"high"`   // 最高价
	Vol    float64   `json:"vol"`    // 成交额
	Bid    []float64 `json:"bid"`    // [买1价, 买1量]
	Ask    []float64 `json:"ask"`    // [卖1价, 卖1量]
}

type TickerReturn struct {
	Status  string `json:"status"` // 请求处理结果
	Ts      int64  `json:"ts"`     // 响应生成时间点
	Tick    Ticker `json:"tick"`   // K线聚合数据
	Ch      string `json:"ch"`     // 数据所属的Channel, 格式: market.$symbol.detail.merged
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type TimestampReturn struct {
	Status  string `json:"status"` // 请求状态
	Data    int64  `json:"data"`   // 时间戳
	ErrCode string `json:"err-code"`
	ErrMsg  string `json:"err-msg"`
}

type TradeData struct {
	ID        int64   `json:"id"`        //成交ID
	Price     float64 `json:"price"`     // 成交价
	Amount    float64 `json:"amount"`    // 成交量
	Direction string  `json:"direction"` // 主动成交方向
	Ts        int64   `json:"ts"`        // 成交时间
}

type TradeTick struct {
	ID   int64       `json:"id"`   // 消息ID
	Ts   int64       `json:"ts"`   // 最新成交时间
	Data []TradeData `json:"data"` // Trade数据
}

type TradeReturn struct {
	Status  string      `json:"status"` // 请求状态, ok或者error
	Ch      string      `json:"ch"`     // 数据所属的Channel, 格式: market.$symbol.trade.detail
	Ts      int64       `json:"ts"`     // 发送时间
	Data    []TradeTick `json:"data"`   // 成交记录
	ErrCode string      `json:"err-code"`
	ErrMsg  string      `json:"err-msg"`
}

type TradeDetailData struct {
	ID        int64   `json:"id"`        // 成交ID
	Price     float64 `json:"price"`     // 成交价
	Amount    float64 `json:"amount"`    // 成交量
	Direction string  `json:"direction"` // 主动成交方向
	Ts        int64   `json:"ts"`        // 成交时间
}

type TradeDetail struct {
	ID   int64             `json:"id"`   // 消息ID
	Ts   int64             `json:"ts"`   // 最新成交时间
	Data []TradeDetailData `json:"data"` // 交易细节数据
}

type TradeDetailReturn struct {
	Status  string      `json:"status"`   // 请求处理结果, "ok"、"error"
	Ts      int64       `json:"ts"`       // 响应生成时间点, 单位毫秒
	Tick    TradeDetail `json:"tick"`     // TradeDetail数据
	Ch      string      `json:"ch"`       // 数据所属的Channel, 格式: market.$symbol.trade.detail
	ErrCode string      `json:"err-code"` // 错误代码
	ErrMsg  string      `json:"err-msg"`  // 错误提示
}
