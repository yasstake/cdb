package trans

import (
	"fmt"
	"testing"
	"time"
)

var database Db

func Open() {

	database.Open("/tmp/")

	database.LoadTime(database.time_chunks[0].start)
}

func TestDbOpen(t *testing.T) {
	var db Db

	db.Open("/tmp/")
	fmt.Println(db)

	db.LoadTime(db.time_chunks[0].start)

	fmt.Println(db.chunk.ohlcv(db.current_start, db.current_end))
	fmt.Println(db.chunk.ohlcvSec())
}

func TestGetBoard(t *testing.T) {
	Open()

	bid, ask, err := database.chunk.GetOrderBook(database.chunk.start_time())

	fmt.Println(err)
	fmt.Println(bid.Sort(false))
	fmt.Println(ask.Sort(true))

	bid, ask, err = database.chunk.GetOrderBook(database.chunk.start_time().Add(time.Second * 10))

	fmt.Println(err)
	fmt.Println(bid.Sort(false))
	fmt.Println(ask.Sort(true))
}

func BenchmarkGetBoard(b *testing.B) {
	Open()
	for i := 0; i < b.N; i++ {
		_, _, err := database.chunk.GetOrderBook(database.chunk.start_time().Add(time.Second * 50))
		if err != nil {
			fmt.Println(err)
		}
	}
}
