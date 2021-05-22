package trans

import (
	"log"
	"sort"
)

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

// Find edge price
func findPrice(price map[int]int, size int, buy_side bool) int {
	if len(price) == 0 {
		return 0
	}

	orders := make(Order, len(price))

	i := 0
	for price, volume := range price {
		orders[i].Price = price
		orders[i].Volume = volume
		i++
	}

	if buy_side {
		// FInd lowest buy price
		sort.Slice(orders, func(i, j int) bool { return orders[i].Price < orders[j].Price })
	} else {
		// Find highest sell
		sort.Slice(orders, func(i, j int) bool { return orders[i].Price > orders[j].Price })
	}

	target_vol := orders[0].Price * size

	for i := range orders {
		target_vol -= orders[i].Volume

		if target_vol < 0 {
			return orders[i].Price
		}
	}

	return 0
}

func AppendExecPrice(tr TransactionSlice) (result TransactionSlice) {

	execute := func(tr Transaction) bool {
		action := tr.Action

		if action == TRADE_SELL || action == TRADE_BUY {
			return true
		}
		return false
	}

	exec_tran := tr.Where(execute)

	const DURATION = 10 * 1_000_000_000 // sec
	const EXEC_WAIT = 1 * 1_000_000_000 // sec

	sell_price := map[int]int{}
	buy_price := map[int]int{}

	l := len(exec_tran)
	var i, j int
	var lowest_buy, highest_sell int
	var last_lowest_buy, last_highest_sell int
	var last_lowest_buy_time, last_highest_sell_time int64
	var price_change_tran TransactionSlice

	for i = 0; i < l; i++ {
		start_time := exec_tran[i].Time_stamp
		action := exec_tran[i].Action

		for ; j < l; j++ {
			tr := exec_tran[j]
			if start_time+DURATION < tr.Time_stamp {
				break
			}

			if tr.Action == TRADE_BUY {
				buy_price[int(tr.Price)] += int(tr.Volume)
			} else if tr.Action == TRADE_SELL {
				sell_price[int(tr.Price)] += int(tr.Volume)
			} else {
				log.Fatal("Unkonw log", tr)
			}
		}

		lowest_buy = findPrice(buy_price, 0, true)
		highest_sell = findPrice(sell_price, 0, false)
		t := exec_tran[i].Time_stamp

		if last_lowest_buy != lowest_buy && last_lowest_buy_time != t {
			tr := Transaction{TRADE_BUY_PRICE, exec_tran[i].Time_stamp - EXEC_WAIT, int32(last_lowest_buy), 0, 0}
			price_change_tran = append(price_change_tran, tr)
			last_lowest_buy = lowest_buy
			last_lowest_buy_time = t
		}

		if last_highest_sell != highest_sell && last_highest_sell_time != t {
			// fmt.Println("SELL CAHGEN", highest_sell)
			tr := Transaction{TRADE_SELL_PRICE, exec_tran[i].Time_stamp - EXEC_WAIT, int32(last_highest_sell), 0, 0}
			price_change_tran = append(price_change_tran, tr)
			last_highest_sell = highest_sell
			last_highest_sell_time = t
		}

		if action == TRADE_BUY {
			buy_price[int(exec_tran[i].Price)] -= int(exec_tran[i].Volume)
		} else if action == TRADE_SELL {
			sell_price[int(exec_tran[i].Price)] -= int(exec_tran[i].Volume)
		}
	}

	tr = append(tr, price_change_tran...)
	tr.TimeSort()

	return tr
}
