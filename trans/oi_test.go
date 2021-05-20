package trans

import (
	"fmt"
	"testing"
)

func TestFindCombination(t *testing.T) {
	item := IntOiItem{[]int{100, 200, 10}, []bool{false, false, false, false, false, false}}

	FindCombination(&item, 0, 510)
	fmt.Println(item)

	//FindCombination(item, 0, 300)
}

func TestFindCombination2(t *testing.T) {
	item := IntOiItem{[]int{100, 200, 10}, []bool{false, false, false, false, false, false}}

	r := FindCombination(&item, 0, 500)
	fmt.Println(item, r)

	r2 := make_oc_slice(item.mask)
	fmt.Println(r2)
}

func TestFindCombination3(t *testing.T) {
	item := IntOiItem{[]int{100000, 29040}, []bool{false, false, false, false}}

	r := FindCombination(&item, 0, 129040)
	fmt.Println(item, r)

	r2 := make_oc_slice(item.mask)
	fmt.Println(r2)
}

func TestTrimatch(t *testing.T) {
	t1 := Transaction{1, 1, 1, 1, 0}
	t2 := Transaction{2, 4, 1, 1, 0}
	t3 := Transaction{3, 3, 1, 1, 0}

	trs := TransactionSlice{t1, t2, t3}

	r := FindTriMatch(&trs, -2)
	fmt.Println(trs, r)
}

func TestBestmatch(t *testing.T) {
	t1 := Transaction{1, 1, 1, 1, 0}
	t2 := Transaction{2, 4, 1, 1, 0}
	t3 := Transaction{3, 3, 2, 1, 0}

	trs := TransactionSlice{t1, t2, t3}

	var target = 10
	r := FindBestMatch(&trs, target)
	fmt.Println(trs, target, r)
}
