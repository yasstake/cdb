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

func (c *Db) GetTimeChunks() TimeFrames {
	return c.time_chunks
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
	err := c.chunk.LoadTime(t)

	if err != nil {
		// load failed unchanged.
		log.Println("DB chunk file is not found time=", t)
		return c.chunk, err
	}

	c.current_start = c.chunk.start_time()
	c.current_end = c.chunk.end_time()

	return c.chunk, nil
}

// Open and Load next chunk(1 min after)
func (c *Db) LoadNext() (Chunk, error) {
	next_time := c.current_start.Add(time.Minute + time.Second)
	return c.LoadTime(next_time)
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

// CASE1:   s < e < chunk(start) < chunk(end)
// CASE2:   s < chunk(start) < e < chunk(end)
// CASE3:   chunk(start) < s < e < chunk(end)
// CASE4:   s < chunk(start) < chunk(end) < e
// CASE5:   chunk(start) < s < chunk(end) < e
// CASE6:   chunk(start) < chunk(end) < s < e

const BOUND_BEFORE = -1
const BOUND_IN = 0
const BOUND_AFTER = 1

type BoundStatus int

// Check bounds return
//   CASE1 BOUND_BEFORE 	t   < time_frame
//   CASE1 BOUND_BEFORE 	t   < time_frame
//   CASE1 BOUND_BEFORE 	t   < time_frame
func check_bounds(frames TimeFrames, s time.Time, e time.Time) (bound_s, bound_e BoundStatus, err error) {
	err = nil

	if frames.Before(s) {
		bound_s = BOUND_BEFORE
	} else if frames.In(s) {
		bound_s = BOUND_IN
	} else if frames.After(s) {
		bound_s = BOUND_AFTER
	}

	if frames.Before(e) {
		bound_e = BOUND_BEFORE
	} else if frames.In(e) {
		bound_e = BOUND_IN
	} else if frames.After(e) {
		bound_e = BOUND_AFTER
	}

	if e.Before(s) {
		err = fmt.Errorf("start time %s must be before %s", s, e)
	}

	return bound_s, bound_e, err
}

/*
// TODO: not implemented
func check_bouds(c *Db, s time.Time, e time.Time) (err error) {
	bs, be, err := check_bounds(c.time_chunks, s, e)
	if bs == BOUND_AFTER || be == BOUND_BEFORE || e.Before(s) {
		err = fmt.Errorf("select time is out of chunk %s, %s", s, e)
		return err
	}
	return nil
}
*/

// TODO: not implemented
func (c *Db) SelectTrans(s time.Time, e time.Time) (err error) {
	bs, be, err := check_bounds(c.time_chunks, s, e)
	if bs == BOUND_AFTER || be == BOUND_BEFORE || e.Before(s) {
		err = fmt.Errorf("select time is out of chunk %s, %s", s, e)
		return err
	}

	if err != nil {
		return err
	}

	// TODO: not implemneted
	return err
}

func (c *Db) SelectOhlcv(s time.Time, e time.Time) error {
	//c.check_boud(s, e)

	load_time := s

	c.LoadTime(load_time)
	// TODO: not imlemented
	// ohlcv, _ := c.chunk.ohlcv(s, e)

	return nil
}
