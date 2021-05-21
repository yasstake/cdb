package trans

// Action
const PARTIAL = 1
const UPDATE_SELL = 2
const UPDATE_BUY = 3

// trade
const TRADE_BUY = 4
const TRADE_BUY_LIQUID = 5

const TRADE_SELL = 6
const TRADE_SELL_LIQUID = 7

// buy edge price
const TRADE_BUY_PRICE = 8

// sell edge price
const TRADE_SELL_PRICE = 9

// Open Interest
// action, time, 0,, volume,
const OPEN_INTEREST = 10

// Open Value
// action, time, 0, volume
const OPEN_VALUE = 11

// Turn Over
// action, time, 0, volume
const TURN_OVER = 12

// Funding Rate
// action, time, 0, volume, next time
const FUNDING_RATE = 20

// Next Funding Rate
// action, time, 0, volume, next time
const PREDICTED_FUNDING_RATE = 21

// Use gen: slice typewriter
// https://clipperhouse.com/gen/slice/

// Store each transaction data from Bybit Exchange
// +gen slice:"Where, GroupBy[int32], Count, SortBy"
type Transaction struct {
	Action     int8
	Time_stamp int64
	Price      int32
	Volume     int64
	OtherInfo  int64
}

//type TransactionSlice []TransactionSlice

// SelectExecute

// SelectTimeSpan

// FindExecutePrice (1BTC execute price)

// SelectLiquid

// DollarBar

//
