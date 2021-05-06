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

	fmt.Println(db.chunk.GetOhlcv(db.current_start, db.current_end))
	fmt.Println(db.chunk.GetOhlcvSec())
}

func TestDbOpenNext(t *testing.T) {
	var db Db

	db.Open("/tmp/")
	fmt.Println(db)

	db.LoadTime(db.time_chunks[0].start)
	fmt.Println(db.chunk.GetOhlcv(db.current_start, db.current_end))

	for {
		_, err := db.LoadNext()
		if err != nil {
			break
		}
		fmt.Println(db.chunk.GetOhlcv(db.current_start, db.current_end))
	}
}

func TestCheckBound(t *testing.T) {
	t1 := DateTime(time.Hour.Nanoseconds())
	t2 := DateTime(time.Hour.Nanoseconds() * 2)

	t3 := DateTime(time.Hour.Nanoseconds() * 4)
	// t4 := date_time(time.Hour.Nanoseconds()*4 + 1)
	t5 := DateTime(time.Hour.Nanoseconds() * 5)

	frame1 := TimeFrame{t1, t2}
	frame2 := TimeFrame{t3, t5}

	frame := TimeFrames{frame1, frame2}

	bs, be, err := check_bounds(frame, t1, t2)

	fmt.Println(bs, be, err)
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
