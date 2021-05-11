package trans

import (
	"fmt"
	"testing"
)

type Value interface {
	Value() int
}

// Calc combination to fit target
func FindCombination(value []int, target int) (result []int) {
	l := len(value)
	combination := make([]int, l)

	loop := int(3 ^ l)
	for i := 0; i < loop; i++ {
		remain = int(i)
		for j := l; 0 < j; j-- {

		}
	}

	fmt.Println(l)
	fmt.Println(combination)

	return combination
}

func TestCOmbination(t *testing.T) {
	FindCombination([]int{10, 20, 30, 40}, 5)
}
