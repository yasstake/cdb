package bb

import (
	"cdb/trans"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

const CHANNEL_ORDER_BOOK_200 = "orderBook_200.100ms.BTCUSD"
const CHANNEL_TRADE = "trade.BTCUSD"
const CHANNEL_INFO = "instrument_info.100ms.BTCUSD"

func make_rec(action int, time int64, price float64, volume float64) string {
	return fmt.Sprintf("%d,%d,%d,%d\n", action, time, int(price*10), int(volume*10))
}

type Response struct {
	Success bool    `json:"success"`
	Message string  `json:"ret_msg"`
	Id      string  `json:"conn_id"`
	Request Request `json:"request"`
}

type Request struct {
	Operation string   `json:"op"`
	Args      []string `json:"args"`
}

type Order struct {
	Price  json.Number `json:"price"`
	Symbol string      `json:"symbol"`
	Id     int64       `json:"id"`
	Side   string      `json:"side"`
	Size   json.Number `json:"size"`
	Time   int64
}

func (c *Order) ToLog() string {
	var action int
	if c.Side == "Buy" {
		action = trans.UPDATE_BUY
	} else if c.Side == "Sell" {
		action = trans.UPDATE_SELL
	} else {
		log.Fatalln("Unknown side", c.Side)
	}

	price, _ := c.Price.Float64()
	volume, _ := c.Size.Float64()

	return make_rec(action, c.Time, price, volume)
}

type SnapShot []Order

type Delta struct {
	Delete []Order `json:"delete"`
	Update []Order `json:"update"`
	Insert []Order `json:"insert"`
}

type Message struct {
	Topic    string          `json:"topic"`
	Type     string          `json:"type"`
	Data     json.RawMessage `json:"data"`
	Sequence int64           `json:"cross_seq"`
	Time     int64           `json:"timestamp_e6"`
}

type Trade struct {
	Topic string          `json:"topic"`
	Type  string          `json:"type"`
	Data  json.RawMessage `json:"data"`
}

type TradeRec struct {
	Time      int64       `json:"trade_time_ms"`
	Timestamp string      `json:"timestamp"`
	Symbol    string      `json:"symbol"`
	Side      string      `json:"side"`
	Size      json.Number `json:"size"`
	Price     json.Number `json:"price"`
}

type TradeRecs []TradeRec

func (c *TradeRec) ToLog() (result string) {
	var action int
	if c.Side == "Buy" {
		action = trans.TRADE_BUY
	} else if c.Side == "Sell" {
		action = trans.TRADE_SELL
	} else {
		log.Fatalln("Unknown side", c.Side)
	}

	price, _ := c.Price.Float64()
	volume, _ := c.Size.Float64()

	return make_rec(action, c.Time, price, volume)
}

type InstrumentSnapshot Instrument

type InstrumentDelta struct {
	Delete []Instrument `json:"delete"`
	Update []Instrument `json:"update"`
	Insert []Instrument `json:"insert"`
}

type Instrument struct {
	Id                int    `json:"id"`
	Symbol            string `json:"symbol"`
	LastPrice         int    `json:"last_price_e4"`       //:536925000,
	BitPrice          int    `json:"bid1_price_e4"`       //:536925000,
	AskPrice          int    `json:"ask1_price_e4"`       //:536930000,
	LastTickDirection string `json:"last_tick_direction"` //:"ZeroMinusTick",
	PrevPrice         int    `json:"prev_price_24h_e4"`   //:503145000,
	// `json:"price_24h_pcnt_e6"`//:67137,
	HighPrice int `json:"high_price_24h_e4"` //:539840000,
	LowPrice  int `json:"low_price_24h_e4"`  //:470000000,
	// `json:"prev_price_1h_e4"`//:537670000,
	//`json:"price_1h_pcnt_e6"`//:-1385,
	MarkPrice            int `json:"mark_price_e4"`             //:536850500,
	IndexPrice           int `json:"index_price_e4"`            //:536796200,
	OpenInterest         int `json:"open_interest"`             //:1905193274,
	OpenValue            int `json:"open_value_e8"`             //:1461680310351,
	TotalTurnOver        int `json:"total_turnover_e8"`         //:7597234355112527,
	TurnOver24h          int `json:"turnover_24h_e8"`           //:17278982841461,
	TotalVolume          int `json:"total_volume"`              //:1465809840880,
	Volume24h            int `json:"volume_24h"`                //:8787821958,
	FundingRate          int `json:"funding_rate_e6"`           //:-51,
	PredictedFundingRate int `json:"predicted_funding_rate_e6"` //:9,
	//`json:"cross_seq"`//:6183139859,
	//`json:"created_at"`//:"2018-11-14T16:33:26Z",
	//`json:"updated_at"`//:"2021-04-26T15:05:12Z",
	//`json:"next_funding_time"`//:"2021-04-26T16:00:00Z",
	//`json:"countdown_hour"`//:1},
	//`json:"cross_seq"`//:6183139917,
	Time int `json:"timestamp_e6"` //:1619449512968696
}

func order_book(m string) (result string) {
	var message Message

	err := json.Unmarshal([]byte(m), &message)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	switch message.Type {
	case "snapshot":
		return order_book_snap(message.Data, message.Time)
	case "delta":
		return order_book_delta(message.Data, message.Time)
	}

	log.Fatalln("Unknown Message type", message.Type)

	return ""
}

func order_book_snap(message json.RawMessage, time int64) (result string) {
	var data SnapShot

	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	l := len(data)
	if l == 0 {
		return ""
	}

	result = make_rec(trans.PARTIAL, time, 0, 0)

	for i := 0; i < l; i++ {
		data[i].Time = time
		result += data[i].ToLog()
	}

	return result
}

func order_book_delta(message json.RawMessage, time int64) (result string) {
	var data Delta

	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	l := len(data.Insert)
	for i := 0; i < l; i++ {
		data.Insert[i].Time = time
		result += data.Insert[i].ToLog()
	}

	l = len(data.Update)
	for i := 0; i < l; i++ {
		data.Update[i].Time = time
		result += data.Update[i].ToLog()
	}

	l = len(data.Delete)
	for i := 0; i < l; i++ {
		data.Delete[i].Time = time
		result += data.Delete[i].ToLog()
	}

	return result
}

func trade(message json.RawMessage) (result string) {
	var data TradeRecs

	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	l := len(data)

	for i := 0; i < l; i++ {
		result += data[i].ToLog()
	}

	return result
}

func instrument(message Message) (result string) {
	return result
}

func Connect() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	// wss://stream.bybit.com/realtime
	// u := url.URL{Scheme: "ws", Host: "stream.bybit.com", Path: "/realtime"}
	u := url.URL{Scheme: "wss", Host: "stream.bybit.com", Path: "/realtime"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var decoded Message
			err = json.Unmarshal([]byte(message), &decoded)
			if err != nil {
				log.Println("Parse error", err)
				return
			}
			// fmt.Println(decoded.Time, decoded.Sequence, decoded.Topic, decoded.Type, decoded.Sequence)

			switch decoded.Topic {
			case CHANNEL_ORDER_BOOK_200:
				s := order_book(string(message))
				fmt.Println(s)
			case CHANNEL_TRADE:
				s := trade(decoded.Data)
				fmt.Println(s)
			case CHANNEL_INFO:
				fmt.Println("INFO", string(message))
			case "":
				var response Response
				json.Unmarshal([]byte(message), &response)
				fmt.Println(response.Success, response.Id, response.Message, response.Request.Args)
			}

			if decoded.Time == 0 {
				fmt.Println(string(message))
			}
		}
	}()

	subscribe := func(ch string) {
		param := make(map[string]interface{})
		param["op"] = "subscribe"
		args := []string{ch}
		param["args"] = args
		req, _ := json.Marshal(param)
		c.WriteMessage(websocket.TextMessage, []byte(req))
	}

	subscribe("orderBook_200.100ms.BTCUSD")
	subscribe("trade.BTCUSD")
	subscribe("instrument_info.100ms.BTCUSD")

	for {
		select {
		case <-done:
			return
		/*
			case t := <-ticker.C:
				err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
				if err != nil {
					log.Println("write:", err)
					return
				}
		*/
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
