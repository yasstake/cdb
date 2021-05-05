package bb

import (
	"cdb/trans"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	r += time.String() + " "
	r += c.Id.String() + " "
	r += c.Price.String() + " "
	r += c.Volume.String() + " "
	r += c.Side

	return r
}

type RestResponse struct {
	Code    json.Number     `json:"ret_code"`
	Message string          `json:"ret_msg"`
	ExtCode string          `json:"ext_code"`
	Result  json.RawMessage `json:"result"`
	Time    json.Number     `json:"time_now"`
}

func LiquidRequest(from_ms int64) (body string, err error) {
	url := "https://api.bybit.com/v2/public/liq-records?symbol=BTCUSD"
	if from_ms != 0 {
		url += "&start_time=" + strconv.Itoa(int(from_ms))
	}

	req, _ := http.NewRequest("GET", url, nil)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	byteArray, err := ioutil.ReadAll(resp.Body)

	return string(byteArray), err
}

func LiquidMessage(message string) (response RestResponse, liquid []LiquidRec, err error) {
	err = json.Unmarshal([]byte(message), &response)

	if err != nil {
		log.Println(err)
		return response, liquid, err
	}

	// var liquid []LiquidRec
	err = json.Unmarshal([]byte(response.Result), &liquid)

	if err != nil {
		log.Println(err)
		return response, liquid, err
	}
	return response, liquid, nil
}

func LiquidMessageStr(message string) string {
	response, liquid, _ := LiquidMessage(message)

	t, _ := response.Time.Float64()
	fmt.Println(trans.DateTime(int64(t * float64(time.Second))).String())

	var s string
	for i := range liquid {
		s += liquid[i].ToString() + "\n"
	}

	return s
}
