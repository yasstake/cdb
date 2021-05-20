package trans

import (
	"fmt"
	"log"
	"strconv"
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

var tri_cache map[string]int

// TODO: not implemented
//
func FindTriMatch(item OiItem, target int) (result int) {
	tri_cache = make(map[string]int)
	l := item.Len()

	var total int
	for i := 0; i < l; i++ {
		total += item.Get(i)
	}
	fmt.Println("TriTarget=", total+target)

	diff := FindTriMatchRaw(item, 0, total+target)

	if diff == 0 {
		return diff
	}

	fmt.Println("RETRY DIFF", diff)
	return FindTriMatchRaw(item, 0, total+target-diff)
}

func FindTriMatchRaw(item OiItem, offset int, target int) int {

	v, ok := tri_cache[strconv.Itoa(offset)+"-"+strconv.Itoa(target)]
	if ok {
		// fmt.Print("cache")
		return v
	}

	result := FindTriMatchRawNoCache(item, offset, target)
	tri_cache[strconv.Itoa(offset)+"-"+strconv.Itoa(target)] = result

	return result
}

func FindTriMatchRawNoCache(item OiItem, offset int, target int) int {
	l := item.Len()

	if offset == l {
		if target == 0 {
			return 0
		} else {
			return target
		}
	} else {
		r1 := FindTriMatchRaw(item, offset+1, target)
		if r1 == 0 {
			return 0
		}

		diff := target - item.Get(offset)*2
		if 0 <= diff {
			r2 := FindTriMatchRaw(item, offset+1, diff)
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

func FindBestMatch(item OiItem, target int) int {
	l := item.Len()

	var total int
	for i := 0; i < l; i++ {
		total += item.Get(i)
	}

	fmt.Println("Total", total, target, total+target)

	return FindBestMatchRaw(item, 0, total+target)
}

func FindBestMatchRaw(item OiItem, offset int, target int) int {

	l := item.Len()

	if offset == l {
		if target == 0 {
			return 0
		} else {
			return target
		}
	}

	r1 := FindBestMatchRaw(item, offset+1, target)
	var r2 int
	var item_value = item.Get(offset) * 2
	diff := target - item_value
	if 0 <= diff {
		r2 = FindBestMatchRaw(item, offset+1, diff)
	}

	if r1 < r2 {
		fmt.Print("[", item.Get(offset), ",")
		return r2
	} else {
		return r1
	}
}
