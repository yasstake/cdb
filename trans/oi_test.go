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
