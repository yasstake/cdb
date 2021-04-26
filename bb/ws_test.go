package bb

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func parse_message(message string) (result Message) {
	err := json.Unmarshal([]byte(message), &result)

	if err != nil {
		log.Println("Parse Error")
	}

	return result
}

func TestOrderBook(t *testing.T) {
	result := order_book(ORDER_BOOK_SNAP_RECORD)
	fmt.Println(result)

	result = order_book(ORDER_BOOK_DELTA_RECORD)
	fmt.Println(result)
}

func TestTradeRecord(t *testing.T) {
	result := trade(TRADE_RECORD)

	fmt.Println(result)
}

func TestTrade(t *testing.T) {

}
