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

	FindCombination(&item, 0, 510)
	fmt.Println(item)

	r := make_oc_slice(item.mask)
	fmt.Println(r)
}
