package trans

import (
	"fmt"
	"log"
	"time"
)

type Db struct {
	base_path   string
	time_chunks TimeFrames

	current_start time.Time
	current_end   time.Time
	chunk         Chunk
}

// Open Data store directory
func (c *Db) Open(path string) {
	SetDbRoot(path)
	c.base_path = DB_ROOT
	c.time_chunks = Time_chunks(DB_ROOT)
}

// Return chunk which contains time.
func (c *Db) LoadTime(t time.Time) (Chunk, error) {
	if c.current_start.Before(t) && t.Before(c.current_end) {
		// retrun cached info
		return c.chunk, nil
	}
	err := c.chunk.load_time(t)

	if err != nil {
		// load failed unchanged.
		log.Println("DB chunk file is not found time=", t)
		return c.chunk, err
	}

	c.current_start = c.chunk.start_time()
	c.current_end = c.chunk.end_time()

	return c.chunk, nil
}

/*
// Search Board info and return board info
// TODO: implement time offset
func (c *Db) GetBoard(t time.Time) (bit Board, ask Board) {
	c.LoadTime(t)
	bit = c.chunk.bit_board
	ask = c.chunk.ask_board

	return bit, ask
}
*/

// Retrive order book board information from logdb
func (c *Db) GetBoard(t time.Time) (bid Board, ask Board, err error) {
	if !c.time_chunks.In(t) {
		return nil, nil, fmt.Errorf("out of range time=%s in[%s %s]", t, c.chunk.start_time(), c.chunk.end_time())
	}
	chunk, err := c.LoadTime(t)
	if err != nil {
		return nil, nil, err
	}

	bid, ask, err = chunk.GetOrderBook(t)

	return bid, ask, err
}

func (c *Db) check_boud(s time.Time, e time.Time) {
	if !c.time_chunks.In(s) {
		log.Println("[WARN] out bound in select trans before")
	}
	if !c.time_chunks.In(e) {
		log.Println("[WARN] out bound in select trans end")
	}
	// TODO: not implemented
}

func (c *Db) SelectTrans(s time.Time, e time.Time) error {
	c.check_boud(s, e)
	// TODO: not implemented

	return nil
}

func (c *Db) SelectOhlcv(s time.Time, e time.Time) error {
	//c.check_boud(s, e)

	load_time := s

	c.LoadTime(load_time)
	// TODO: not imlemented
	// ohlcv, _ := c.chunk.ohlcv(s, e)

	return nil
}
