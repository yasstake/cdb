package bb

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

/*
func parse_message(message string) (result Message) {
	err := json.Unmarshal([]byte(message), &result)

	if err != nil {
		log.Println("Parse Error")
	}

	return result
}
*/

func get_data(m string) json.RawMessage {
	var message Message

	err := json.Unmarshal([]byte(m), &message)
	if err != nil {
		log.Fatalln("Fail to pase message", err, message)
	}

	return message.Data
}

func TestParseIosTime(t *testing.T) {
	time := parse_iso_time("2021-04-26T16:00:00Z")

	fmt.Println(time)
}

func TestParseTimeToMs(t *testing.T) {
	ms := time_to_ms("2021-04-26T16:00:00Z")

	fmt.Println(ms)
}

func TestParseTimeTomSS(t *testing.T) {
	ts := time_to_ms_str("2021-04-26T16:00:00Z")

	fmt.Println(ts)
}

func TestOrderBook(t *testing.T) {
	result := order_book(ORDER_BOOK_SNAP_RECORD)
	fmt.Println(result)

	result = order_book(ORDER_BOOK_DELTA_RECORD)
	fmt.Println(result)
}

func TestTradeRecord(t *testing.T) {
	result := trade(get_data(TRADE_RECORD))

	fmt.Println(result)
}

func TestInstrument(t *testing.T) {
	result := instrument_snapshot(get_data(INSTRUMENT_SNAPSHOT), 1000)

	fmt.Println(result)
}
