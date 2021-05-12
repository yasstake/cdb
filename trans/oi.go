package trans

import (
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

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
