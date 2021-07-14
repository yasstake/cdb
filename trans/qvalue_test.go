package trans

import (
	"fmt"
	"testing"
)

func TestCalcQValue(t *testing.T) {
	t1 := Transaction{1, 1, 1, 1, 1}
	t2 := Transaction{1, 1, 1, 1, 1}
	t3 := Transaction{1, 1, 1, 1, 1}
	t4 := Transaction{1, 1, 1, 1, 1}
	t5 := Transaction{1, 1, 1, 1, 1}
	t6 := Transaction{1, 1, 1, 1, 1}
	t7 := Transaction{1, 1, 1, 1, 1}
	t8 := Transaction{1, 1, 1, 1, 1}

	tr := TransactionSlice{t1, t2, t3, t4, t5, t6, t7, t8}
	fmt.Println(tr)

	CalcQ(tr)
}
