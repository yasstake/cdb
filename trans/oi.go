package trans

import (
	"fmt"
	"log"
)

const OI_OPEN_OPEN = 3
const OI_OPEN_CLOSE = 2
const OI_CLOSE_CLOSE = 1

type OiItem interface {
	Get(int) int
	Hit(int)
	Len() int
}

type IntOiItem struct {
	item []int
	mask []bool
}

func (c *IntOiItem) Get(index int) int {
	pos := int(index / 2)

	value := c.item[pos]

	return value
}

func (c *IntOiItem) Hit(index int) {
	fmt.Println("HIT", c.Get(index))
	c.mask[index] = true
}

func (c *IntOiItem) Len() int {
	return len(c.item) * 2
}

func make_oc_slice(mask []bool) (result []int) {
	l := len(mask)
	if l%2 != 0 {
		log.Println("mask length must be 2, 4, 6, ,,,")
	}

	result = make([]int, l/2)

	for i := 0; i < l; i += 2 {
		result_index := int(i / 2)
		if mask[i] && mask[i+1] {
			result[result_index] = OI_OPEN_OPEN
		} else if !mask[i] && !mask[i+1] {
			result[result_index] = OI_CLOSE_CLOSE
		} else {
			result[result_index] = OI_OPEN_CLOSE
		}
	}
	return result
}

func FindCombination(item OiItem, offset int, target int) (remain int) {
	// fmt.Println("CALL ", offset, target)

	l := item.Len()

	if offset == l {
		if target == 0 {
			return 0
		} else {
			return target
		}
	} else {
		r1 := FindCombination(item, offset+1, target)
		if r1 == 0 {
			return 0
		}

		diff := target - item.Get(offset)
		if 0 <= diff {
			r2 := FindCombination(item, offset+1, diff)
			if r2 == 0 {
				item.Hit(offset)
				return 0
			} else {
				return diff
			}
		}
	}
	return target
}

/*
// FInd best match mask for transaction.
// TODO: Even if there are multiple combintion to fit, this function return only first match
func FindBestMatchMask(tr Transactions, target int) (result []int) {
	combi := MakeCombination(len(tr))
	l := len(combi)

	type MaskValue struct {
		mask  []int
		value int
	}

	var values []MaskValue
	for i := 0; i < l; i += 1 {
		v := CalcHamadar(tr, combi[i])
		if v == target {
			return combi[i]
		} else {
			values = append(values, MaskValue{combi[i], v})
		}
	}

	sort.Slice(values, func(i, j int) bool { return values[i].value < values[j].value })

	for i := 0; i < l; i += 1 {
		if values[i].value == target {
			result = values[i].mask
			break
		} else if target < values[i].value {
			log.Printf("approx value")
			result = values[i].mask
			break
		}
	}
	return result
}

func CalcHamadar(tr Transactions, mask []int) (result int) {
	l := len(tr)
	if l != len(mask) {
		log.Println("error mismatch len")
	}

	result = 0
	for i := 0; i < l; i += 1 {
		result += int(tr[i].Volume) * mask[i]
	}

	return result
}

// Make All combination of [combination] length in 3digits
func MakeCombination(combination int) (result [][]int) {
	s := makeCombinationString(combination)

	l := len(s)
	result = make([][]int, l)

	for i := 0; i < l; i++ {
		result[i] = string_to_array(s[i])
	}

	return result
}

// Make All combination of 3("0", "1", "2") in string
func makeCombinationString(combination int) (result []string) {
	repeat := math.Pow(3, float64(combination))
	for i := 0; i < int(repeat); i += 1 {
		num := strconv.FormatInt(int64(i), 3)
		diff := combination - len(num)

		if 0 < diff {
			filler := strings.Repeat("0", diff)
			num = filler + num
		}
		result = append(result, num)
	}

	return result
}

func string_to_array(var_string string) (result []int) {

	l := len(var_string)
	result = make([]int, l)

	for i := 0; i < l; i += 1 {
		ch := string(var_string[i])
		if ch == "0" {
			result[i] = -2
		} else if ch == "1" {
			result[i] = 0
		} else if ch == "2" {
			result[i] = 2
		} else {
			log.Println("error format char", ch)
		}
	}

	return result
}
*/
