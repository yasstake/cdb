package trans

import (
	"fmt"
	"os"
	"testing"
)

func TestCsvWrite(t *testing.T) {
	t1 := Transaction{1, 1, 1, 1, 1}
	t2 := Transaction{2, 4, 1, 1, 1}
	t3 := Transaction{3, 3, 1, 1, 1}

	trs := Transactions{t1, t2, t3}

	CsvWrite(trs, os.Stdout)

	trs.time_sort()
	CsvWrite(trs, os.Stdout)
}

func TestLogLoad(t *testing.T) {
	r := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	r.time_sort()
}

func TestLogLoad2(t *testing.T) {
	tran := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	tran.time_sort()

	var max_rec int

	for i := range tran {
		action := tran[i].Action
		if action == PARTIAL || action == UPDATE_BUY || action == UPDATE_SELL {
			continue
		}
		t := tran[i]
		fmt.Printf("%2d %20d %10d %20d %16d\n", t.Action, t.Time_stamp, t.Price, t.Volume, t.NextTime)

		if 200 < max_rec {
			break
		}
		max_rec += 1
	}
}
