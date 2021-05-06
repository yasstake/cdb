package bb

import (
	"cdb/trans"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type LiquidRec struct {
	Id     json.Number `json:"id"`
	Price  json.Number `json:"price"`
	Volume json.Number `json:"qty"`
	Time   json.Number `json:"time"`
	Side   string      `json:"side"`
}

func (c *LiquidRec) ToString() (r string) {
	t, _ := c.Time.Int64()
	time := trans.DateTime(t * int64(time.Millisecond))
	r += time.UTC().String() + "(" + strconv.Itoa(int(t)) + ")"
	r += c.Id.String() + " "
	r += c.Price.String() + " "
	r += c.Volume.String() + " "
	r += c.Side

	return r
}

func LiquidRequest(from_ms int64) (body string, time time.Time, err error) {
	url := "https://api.bybit.com/v2/public/liq-records?symbol=BTCUSD"
	if from_ms != 0 {
		url += "&start_time=" + strconv.Itoa(int(from_ms))
	}

	result, time, err := RestRequest(url)

	return result, time, err
}

func LiquidMessage(message string) (response RestResponse, liquid []LiquidRec, err error) {
	// var liquid []LiquidRec
	err = json.Unmarshal([]byte(message), &liquid)

	if err != nil {
		log.Println(err)
		return response, liquid, err
	}
	return response, liquid, nil
}

func LiquidMessageStr(message string) string {
	_, liquid, _ := LiquidMessage(message)

	var s string
	for i := range liquid {
		s += liquid[i].ToString() + "\n"
	}

	return s
}
