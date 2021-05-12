package trans

import (
	"fmt"
	"testing"
)

func TestFindBestMatchMask(t *testing.T) {
	t1 := Transaction{1, 1, 1, 3, 1}
	t2 := Transaction{2, 4, 1, 2, 1}
	t3 := Transaction{3, 3, 1, 1, 1}

	trs := Transactions{t1, t2, t3}

	result := FindBestMatchMask(trs, 4)
	fmt.Println(result)
}

func TestCalcHamadar(t *testing.T) {
	t1 := Transaction{1, 1, 1, 3, 1}
	t2 := Transaction{2, 4, 1, 2, 1}
	t3 := Transaction{3, 3, 1, 1, 1}

	trs := Transactions{t1, t2, t3}
	mask := []int{-2, 0, 2}

	result := CalcHamadar(trs, mask)
	fmt.Println(result)
}

func TestMakeCombinationString(t *testing.T) {
	fmt.Println(makeCombinationString(2))
	fmt.Println(makeCombinationString(3))
	fmt.Println(makeCombinationString(4))
}

func TestMakeCombination(t *testing.T) {
	fmt.Println(MakeCombination(1))
	fmt.Println(MakeCombination(2))
	fmt.Println(MakeCombination(3))
	fmt.Println(MakeCombination(4))
	fmt.Println(MakeCombination(5))
}

func TestStringToArray(t *testing.T) {
	result := string_to_array("012012")
	fmt.Println(result)
}

// Calc combination to fit target
func FindCombination(value []int, target int) (result []int) {
	l := len(value)
	combination := make([]int, l)

	loop := int(3 ^ l)
	for i := 0; i < loop; i++ {
		// remain = int(i)
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
