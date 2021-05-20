package trans

import (
	"fmt"
	"log"
	"sort"
	"testing"
)

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

func TestCalcExecPrice(t *testing.T) {
	r := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	r.TimeSort()

	execute := func(tr Transaction) bool {
		action := tr.Action

		if action == TRADE_SELL || action == TRADE_BUY {
			return true
		}
		return false
	}

	exec_tran := r.Where(execute)

	const DURATION = 10 * 1_000_000_000 // sec

	sell_price := map[int]int{}
	buy_price := map[int]int{}

	l := len(exec_tran)
	var i, j int
	var count int
	var lowest_buy, highest_sell int
	var last_lowest_buy, last_highest_sell int

	for i = 0; i < l; i++ {
		start_time := exec_tran[i].Time_stamp
		action := exec_tran[i].Action

		// price[int(exec_tran[i].Price)] += int(exec_tran[i].Volume)

		// add

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

		if last_lowest_buy != lowest_buy {
			fmt.Println("BUY CHANGE", lowest_buy)
			last_lowest_buy = lowest_buy
		}

		if last_highest_sell != highest_sell {
			fmt.Println("SELL CAHGEN", highest_sell)
			last_highest_sell = highest_sell
		}

		if action == TRADE_BUY {
			buy_price[int(exec_tran[i].Price)] -= int(exec_tran[i].Volume)
		} else if action == TRADE_SELL {
			sell_price[int(exec_tran[i].Price)] -= int(exec_tran[i].Volume)
		}
	}

	//fmt.Println(buy_price)
	//fmt.Println(sell_price)
	fmt.Println(count)
}
