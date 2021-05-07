package bb

import (
	"fmt"
	"testing"
)

const LS_RATIO_MESSAGE = `[{"symbol":"BTCUSD","buy_ratio":0.5744,"sell_ratio":0.4256,"timestamp":1620390000},{"symbol":"BTCUSD","buy_ratio":0.5745,"sell_ratio":0.4255,"timestamp":1620389700}]`

func TestLsRateRequiest(t *testing.T) {
	r, _, _ := LsRatioRequest(0)

	fmt.Println(r)
}

func TestLsRatioMessage(t *testing.T) {
	ls, err := LsRatioMessage(LS_RATIO_MESSAGE)

	fmt.Println(ls)
	fmt.Println(err)
}
