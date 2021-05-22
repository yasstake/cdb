package trans

import (
	"fmt"
	"testing"
)

func TestCalcExecPrice(t *testing.T) {
	tr := LogLoad("../DATA/2021-05-05T23-07-09.log.gz")

	tr.TimeSort()

	tr = AppendExecPrice(tr)
	CsvWriteToFile(tr, "/tmp/tran.csv")
	fmt.Println(len(tr))
	//fmt.Println(buy_price)
	//fmt.Println(sell_price)
}
