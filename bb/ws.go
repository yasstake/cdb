package bb

import (
	"cdb/trans"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const CHANNEL_ORDER_BOOK_200 = "orderBook_200.100ms.BTCUSD"
const CHANNEL_TRADE = "trade.BTCUSD"
const CHANNEL_INFO = "instrument_info.100ms.BTCUSD"

var (
	last_time  int64
	last_price int
)

func make_rec(action int, time int64, price float64, volume float64) (result string) {
	price10 := int(price * 10)

	result = fmt.Sprintf("%d,%d,%d,%d\n", action, time-last_time, price10-last_price, int(volume*10))
	last_time = time
	last_price = price10

	return result
}

func make_rec_op(action int, time int64, price float64, volume float64, option string) (result string) {
	price10 := int(price * 10)
	result = fmt.Sprintf("%d,%d,%d,%d,%s\n", action, time-last_time, price10-last_price, int(volume*10), option)
	last_time = time
	last_price = price10

	return result
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

	return make_rec_op(action, c.Time, price, volume, strconv.Itoa(int(c.Id)))
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

	price, err := c.Price.Float64()
	if err != nil {
		log.Println(err)
	}
	volume, err := c.Size.Float64()
	if err != nil {
		log.Println(err)
	}

	return make_rec(action, c.Time, price, volume)
}

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
	NextFundingTime string `json:"next_funding_time"` //:"2021-04-26T16:00:00Z",
	//`json:"countdown_hour"`//:1},
	Time int64
}

func (c *Instrument) update(d Instrument) {
	if d.Id != 0 {
		c.Id = d.Id
	}
	if d.Symbol != "" {
		c.Symbol = d.Symbol
	}
	if d.LastPrice != 0 {
		c.LastPrice = d.LastPrice
	}
	if d.BitPrice != 0 {
		c.BitPrice = d.BitPrice
	}
	if d.AskPrice != 0 {
		c.AskPrice = d.AskPrice
	}
	if d.LastTickDirection != "" {
		c.LastTickDirection = d.LastTickDirection
	}
	if d.PrevPrice != 0 {
		c.PrevPrice = d.PrevPrice
	}
	if d.HighPrice != 0 {
		c.HighPrice = d.HighPrice
	}
	if d.LowPrice != 0 {
		c.LowPrice = d.LowPrice
	}
	if d.MarkPrice != 0 {
		c.MarkPrice = d.MarkPrice
	}
	if d.IndexPrice != 0 {
		c.IndexPrice = d.IndexPrice
	}
	if d.OpenInterest != 0 {
		c.OpenInterest = d.OpenInterest
	}
	if d.OpenValue != 0 {
		c.OpenValue = d.OpenValue
	}
	if d.TotalTurnOver != 0 {
		c.TotalTurnOver = d.TotalTurnOver
	}
	if d.TurnOver24h != 0 {
		c.TurnOver24h = d.TurnOver24h
	}
	if d.TotalVolume != 0 {
		c.TotalVolume = d.TotalVolume
	}
	if d.Volume24h != 0 {
		c.Volume24h = d.Volume24h
	}
	if d.FundingRate != 0 {
		c.FundingRate = d.FundingRate
	}
	if d.PredictedFundingRate != 0 {
		c.PredictedFundingRate = d.PredictedFundingRate
	}
	if d.NextFundingTime != "" {
		c.NextFundingTime = d.NextFundingTime
	}
	if d.Time != 0 {
		c.Time = d.Time
	}

}

func (c *Instrument) ToLog() (result string) {
	result = ""

	t := c.Time / 1000

	// Open Interest
	if c.OpenInterest != 0 {
		result += make_rec_op(trans.OPEN_INTEREST, t, 0, float64(c.OpenInterest), "")
	}

	// Open Value
	if c.OpenValue != 0 {
		result += make_rec_op(trans.OPEN_VALUE, t, 0, float64(c.OpenValue), "")
	}

	// TurnOver
	if c.TotalTurnOver != 0 {
		result += make_rec_op(trans.TURN_OVER, t, 0, float64(c.TotalTurnOver), "")
	}

	if c.FundingRate != 0 {
		if c.NextFundingTime != "" {
			result += make_rec_op(trans.FUNDING_RATE, t, 0, float64(c.FundingRate), time_to_ms_str(c.NextFundingTime))
		} else {
			result += make_rec_op(trans.FUNDING_RATE, t, 0, float64(c.FundingRate), time_to_ms_str(instrument_data.NextFundingTime))
		}
	}

	if c.PredictedFundingRate != 0 {
		if c.NextFundingTime != "" {
			result += make_rec_op(trans.PREDICTED_FUNDING_RATE, t, 0, float64(c.PredictedFundingRate), time_to_ms_str(c.NextFundingTime))
		} else {
			result += make_rec_op(trans.FUNDING_RATE, t, 0, float64(c.FundingRate), time_to_ms_str(instrument_data.NextFundingTime))
		}
	}

	return result
}

func parse_iso_time(t string) time.Time {
	const layout = "2006-01-02T15:04:05Z"
	result, err := time.Parse(layout, t)

	if err != nil {
		log.Println("Dateformat error in log ", err, t)
	}

	return result
}

func time_to_ms(ts string) int64 {
	t := parse_iso_time(ts)
	ns := t.UnixNano()

	return int64(ns / 1_000_000)
}

func time_to_ms_str(ts string) string {
	ms := time_to_ms(ts)

	return strconv.Itoa(int(ms))
}

func order_book(m string) (result string) {
	var message Message

	err := json.Unmarshal([]byte(m), &message)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	switch message.Type {
	case "snapshot":
		return order_book_snap(message.Data, int64(message.Time/1000))
	case "delta":
		return order_book_delta(message.Data, int64(message.Time/1000))
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
	result = ""

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

	result = ""

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

var (
	instrument_data Instrument
)

func instrument_snapshot(message json.RawMessage, time int64) (result string) {
	result = ""

	err := json.Unmarshal(message, &instrument_data)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	instrument_data.Time = time
	result += instrument_data.ToLog()

	return result
}

func instrument_delta(message json.RawMessage, time int64) (result string) {
	var data InstrumentDelta

	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	result = ""

	for i := range data.Update {
		instrument_data.update(data.Update[i])

		data.Update[i].Time = time
		result += data.Update[i].ToLog()
	}

	// Assume Delete and Insert message is not implemented
	for i := range data.Delete {
		data.Delete[i].Time = time
		log.Println("INFO delete ", data.Delete[i])

		result += data.Delete[i].ToLog()
	}

	for i := range data.Insert {
		log.Println("INFO Insert", data.Insert[i])
		data.Insert[i].Time = time
		result += data.Insert[i].ToLog()
	}

	return result
}

func Connect(flag_file_name string, w io.WriteCloser, close_wait_min int) {

	var flag_file FlagFile
	flag_file.Init(flag_file_name)
	flag_file.Create()
	peer_reset := make(chan struct{})

	// wait 300 sec to terminate
	go flag_file.Check_other_process_loop(300, peer_reset)

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

	var mutex sync.Mutex

	write := func(s string) {
		mutex.Lock()
		defer mutex.Unlock()

		w.Write([]byte(s))
	}

	go func() {
		var last_liquid_time int64
		var sleep_time int

		for {
			liqs, _, err := LiquidRequest(&last_liquid_time)
			if err != nil {
				log.Println(err)
			}

			if len(liqs) != 0 {
				write(liqs.ToLog())
				log.Println("liquid ", len(liqs), " records")
				sleep_time = 1
			} else {
				log.Println("liquid sleep", sleep_time)
				sleep_time = sleep_time + 5
				if 30 <= sleep_time {
					sleep_time = 30
				}
			}

			time.Sleep(time.Duration(sleep_time * int(time.Second)))
		}
	}()

	go func() {
		var message_count int
		var board_update_count int
		var trade_count int
		var info_count int

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("[ERROR] ReadMessaeg:", err)
				close(interrupt)
				return
			}
			var decoded Message
			err = json.Unmarshal([]byte(message), &decoded)
			if err != nil {
				// skip message
				log.Println("Parse error", err)
				continue
			}

			if message_count%1000 == 0 {
				log.Printf("%d total / %d board update/ %d trade/ %d info",
					message_count, board_update_count, trade_count, info_count)
			}
			message_count += 1

			switch decoded.Topic {
			case CHANNEL_ORDER_BOOK_200:
				board_update_count += 1
				s := order_book(string(message))
				write(s)
			case CHANNEL_TRADE:
				trade_count += 1
				s := trade(decoded.Data)
				write(s)
			case CHANNEL_INFO:
				info_count += 1
				s := ""

				if decoded.Type == "snapshot" {
					s = instrument_snapshot(decoded.Data, decoded.Time)
				} else if decoded.Type == "delta" {
					s = instrument_delta(decoded.Data, decoded.Time)
				} else {
					log.Println("unknown instrument info type", string(message))
				}
				write(s)

			default:
				log.Println("[OTHER CHANNEL]", string(message))
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
		case <-peer_reset:
			log.Println("Peer reset")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			w.Close()
			goto close_wait

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
			}
			w.Close()
			goto exit
		}
	}

close_wait:
	{
		log.Println("Peer reset close")

		s := 0
		for s < close_wait_min {
			s += 1
			time.Sleep(time.Minute) // sleep min
			log.Printf("[wait min] %4d/%d", s, close_wait_min)
		}
	}
exit:
	log.Println("Logger End")
}
