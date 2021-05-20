package trans

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
