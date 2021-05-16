package bb

import (
	"cdb/trans"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

type LiquidRec struct {
	Id     int64       `json:"id"`
	Price  json.Number `json:"price"`
	Volume json.Number `json:"qty"`
	Time   json.Number `json:"time"`
	Side   string      `json:"side"`
}

func (c *LiquidRec) ToString() (r string) {
	t, _ := c.Time.Int64()
	time := trans.DateTime(t * int64(time.Millisecond))
	r += time.UTC().String() + "(" + strconv.Itoa(int(t)) + ")"
	r += strconv.Itoa(int(c.Id)) + " "
	r += c.Price.String() + " "
	r += c.Volume.String() + " "
	r += c.Side

	return r
}

func (c *LiquidRec) ToLog() (r string) {
	var action int
	if c.Side == "Sell" {
		action = trans.TRADE_SELL_LIQUID
	} else if c.Side == "Buy" {
		action = trans.TRADE_BUY_LIQUID
	} else {
		log.Println("unknown action side", c.Side)
	}

	price, _ := c.Price.Float64()
	volume, _ := c.Volume.Float64()
	time, _ := c.Time.Float64()

	return make_rec_op(action, int64(time), price, volume, strconv.Itoa(int(c.Id)))
}

type LiquidRecs []LiquidRec

func (c *LiquidRecs) ToLog() (result string) {
	l := len(*c)

	for i := 0; i < l; i++ {
		r := (*c)[i].ToLog()
		result += r
	}
	return result
}

// Request Liquid rest API
func LiquidRequest(from_id *int64) (liq LiquidRecs, time time.Time, err error) {
	url := "https://api.bybit.com/v2/public/liq-records?symbol=BTCUSD&limit=1000"

	if *from_id != 0 {
		url += "&from=" + strconv.Itoa(int(*from_id))
	}

	body, time, err := RestRequest(url)
	if err != nil {
		log.Println(err)
	}

	liq, err = LiquidMessage(body)
	if err != nil {
		log.Println(err)
	}

	l := len(liq)
	if l != 0 {
		*from_id = liq[l-1].Id + 1
	}

	return liq, time, err
}

// Parse Liquid JSON message and returns Header and body
func LiquidMessage(message string) (liquid LiquidRecs, err error) {
	// var liquid []LiquidRec
	err = json.Unmarshal([]byte(message), &liquid)

	if err != nil {
		log.Println(err)
		return liquid, err
	}
	return liquid, nil
}

// Convert Liquid JSON message to string representation(for debug purpose)
func LiquidMessageStr(message string) string {
	liquid, _ := LiquidMessage(message)

	var s string
	for i := range liquid {
		s += liquid[i].ToString() + "\n"
	}

	return s
}
