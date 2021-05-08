package trans

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func Open() Db {
	var database Db
	database.Open("/tmp/")

	return database
}

func TestDbOpen(t *testing.T) {
	var db Db

	db.Open("/tmp/")
	fmt.Println(db)

	db.CreateSession(db.time_chunks[0].start)
}

func TestOpenSession(t *testing.T) {
	database := Open()

	time := database.time_chunks[0].start
	fmt.Println(time.String())

	session, _ := database.CreateSession(time)

	fmt.Println(session)
}

func TestGetTransaction(t *testing.T) {
	database := Open()

	start_time := database.time_chunks[0].start
	fmt.Println(start_time.String())

	session, _ := database.CreateSession(start_time)

	reader, err := session.SelectTrans(start_time, start_time.Add(10*time.Minute))

	if err != nil {
		log.Println(err)
		t.Error()
	}

	for {
		r, err := reader.ReadTran()

		if err != nil {
			break
		}
		fmt.Println(r)
	}

}

/*
func TestLoadAndOhlcv(t *testing.T) {
	var c Chunk

	Open()

	s_time := database.time_chunks[0].start.Add(time.Second)
	e_time := s_time.Add(31 * time.Second)

	ohlcv, err := c.GetOhlcv(s_time, e_time)
	fmt.Println(ohlcv, err)

	c.LoadTime(s_time)
	ohlcv, num_rec := c.GetOhlcv(s_time, e_time)
	fmt.Println(ohlcv, num_rec)
}
*/

/*
func TestDbOpenNext(t *testing.T) {
	var db Db

	db.Open("/tmp/")
	fmt.Println(db)

	session, _ := db.CreateSession(db.time_chunks[0].start)

	fmt.Println(session.SelectOhlcv().GetOhlcv(db.current_start, db.current_end))

	for {
		_, err := db.LoadNext()
		if err != nil {
			break
		}
		fmt.Println(db.chunk.GetOhlcv(db.current_start, db.current_end))
	}
}
*/

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
	database := Open()
	session, _ := database.CreateSession(database.time_chunks[0].start)

	bid, ask, err := session.GetBoard(database.time_chunks[0].start)

	fmt.Println(err)
	fmt.Println(bid.Sort(false))
	fmt.Println(ask.Sort(true))

	bid, ask, err = session.GetBoard(database.time_chunks[0].start.Add(time.Second * 50))

	fmt.Println(err)
	fmt.Println(bid.Sort(false))
	fmt.Println(ask.Sort(true))
}

func BenchmarkGetBoard(b *testing.B) {
	database := Open()
	session, _ := database.CreateSession(database.time_chunks[0].start)

	for i := 0; i < b.N; i++ {
		_, _, err := session.GetBoard(database.time_chunks[0].start.Add(time.Second * 50))
		if err != nil {
			fmt.Println(err)
		}
	}
}
