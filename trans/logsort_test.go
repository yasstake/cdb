package trans

import (
	"fmt"
	"log"
	"os"
	"testing"
)

const TEST_LOG_FILE1 = "../DATA/2021-05-21T00-06-32.log.gz"
const TEST_LOG_FILE2 = "../DATA/2021-05-21T04-11-43.log.gz"
const TEST_LOG_FILE3 = "../DATA/2021-05-21T08-16-54.log.gz"
const TEST_LOG_FILE4 = "../DATA/2021-05-21T12-21-55.log.gz"
const TEST_LOG_FILE5 = "../DATA/2021-05-21T16-27-06.log.gz"
const TEST_LOG_FILE6 = "../DATA/2021-05-21T20-32-07.log.gz"

func TestCsvWrite(t *testing.T) {
	t1 := Transaction{1, 1, 1, 1, 1}
	t2 := Transaction{2, 4, 1, 1, 1}
	t3 := Transaction{3, 3, 1, 1, 1}

	trs := TransactionSlice{t1, t2, t3}

	CsvWrite(trs, os.Stdout)

	trs.TimeSort()
	CsvWrite(trs, os.Stdout)
}

func TestLogLoad(t *testing.T) {
	r := LogLoad(TEST_LOG_FILE1)

	r.TimeSort()
}

func TestLogLoad2(t *testing.T) {
	tran := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	tran.TimeSort()

	var max_rec int

	for i := range tran {
		action := tran[i].Action
		if action == PARTIAL || action == UPDATE_BUY || action == UPDATE_SELL {
			continue
		}
		t := tran[i]
		fmt.Printf("%2d %20d %10d %20d %16d\n", t.Action, t.Time_stamp, t.Price, t.Volume, t.OtherInfo)

		if 200 < max_rec {
			break
		}
		max_rec += 1
	}
}

func TestLogLoad3(t *testing.T) {
	tran := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	tran.TimeSort()

	var max_rec int

	var trans TransactionSlice
	var last_oi int64
	var diff_oi int64

	for i := range tran {
		action := tran[i].Action
		if action == PARTIAL || action == UPDATE_BUY || action == UPDATE_SELL {
			continue
		}
		// t := tran[i]

		if action == TRADE_SELL || action == TRADE_BUY {
			trans = append(trans, tran[i])
		}

		if action == OPEN_INTEREST {
			diff_oi = tran[i].Volume - last_oi

			/*
				var total_volume int
				for i := range trans {
					total_volume += int(trans[i].Volume)
				}
			*/

			fmt.Println("-----")
			fmt.Println("DIFF oi=", diff_oi)

			l := len(trans)
			if l != 0 && l < 20 {
				r := FindTriMatch(&trans, int(diff_oi))
				fmt.Println(r, trans)
			}

			trans = make(TransactionSlice, 0)

			last_oi = tran[i].Volume
		}

		// fmt.Printf("%2d %20d %10d %20d %16d\n", t.Action, t.Time_stamp, t.Price, t.Volume, t.NextTime)

		/*
			if 200 < max_rec {
				break
			}
		*/
		max_rec += 1
	}
}

// {x1, x2, x3, x4} -> {x1, x1, x2, x2, x3, x3, x4, x4}
//
//
//

var result_item []int

func Dp(item []int, i int, max_vol int) (result int, weight int) {
	fmt.Println("index=", i, "maxvol=", max_vol)

	if i == -1 {
		fmt.Println("ERR")
		return -1, 0
	} else {
		w := item[i]

		if max_vol < w {
			fmt.Print("-", w, "-")
			return Dp(item, i-1, max_vol)
		} else if w < max_vol {
			a, wa := Dp(item, i-1, max_vol-w)
			b, wb := Dp(item, i-1, max_vol)
			if a < b {
				fmt.Print("-", w, "-")
				return b, wb
			} else {
				fmt.Print("[", w, "] ")
				return a, wa
			}
		} else if w == max_vol {
			fmt.Println("----HIT-----", w)
			return 0, 0
		}
	}

	log.Println("unknown status")
	return 0, 0 // err
}
