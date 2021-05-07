package bb

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type LsRatioRec struct {
	BuyRatio  float64 `json:"buy_ratio"`
	SellRatio float64 `json:"sell_ratio"`
	Time      int64   `json:"timestamp"`
}

// Request LS ratio rest API
func LsRatioRequest(from_ms int64) (body string, time time.Time, err error) {
	url := "https://api.bybit.com//v2/public/account-ratio?symbol=BTCUSD&period=5min"

	if from_ms != 0 {
		url += "&start_time=" + strconv.Itoa(int(from_ms))
		fmt.Println("[requiest]", url)
	}

	body, time, err = RestRequest(url)

	return body, time, err
}

// Parse LS ratio JSON message and returns Header and body
func LsRatioMessage(message string) (ls_ratio []LsRatioRec, err error) {
	// var liquid []LiquidRec
	err = json.Unmarshal([]byte(message), &ls_ratio)

	if err != nil {
		log.Println(err)
		return ls_ratio, err
	}
	return ls_ratio, nil
}
