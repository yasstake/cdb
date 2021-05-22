package trans

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type Db struct {
	base_path   string
	time_chunks TimeFrames
	/*
		current_start time.Time
		current_end   time.Time
		chunk         Chunk
	*/
}

type DbSession struct {
	db            *Db
	current_start time.Time
	current_end   time.Time
	select_start  int64
	select_end    int64
	chunk_len     int
	current_index int
	chunk         Chunk
}

func (c *DbSession) ToString() (result string) {
	result = c.current_start.String() + " " + c.current_end.String() + " " + strconv.Itoa(c.chunk_len)
	return result
}

// Open Data store directory and setup time chunks(Data begin and ends)
func (c *Db) Open(path string) {
	SetDbRoot(path)
	c.base_path = DB_ROOT
	c.time_chunks = TimeChunks(DB_ROOT)
}

// Return chunk which contains time specified
func (c *Db) CreateSession(t time.Time) (session DbSession, err error) {
	session.db = c
	err = session.LoadTime(t)

	return session, err
}

type TranReader interface {
	ReadTran() (Transaction, error)
}

type OhlcvReader interface {
	ReadOhlcv() (Ohlcv, error)
}

// Load data at time to session
func (c *DbSession) LoadTime(t time.Time) (err error) {
	if c.current_start.Before(t) && t.Before(c.current_end) {
		// re-use cached data. but reset read index
		c.current_index = 0
		return nil
	}

	err = c.chunk.LoadTime(t)
	if err != nil {
		// load failed unchanged.
		log.Println("DB chunk file is not found time=", t)
		return err
	}

	c.current_start = c.chunk.start_time()
	c.current_end = c.chunk.end_time()
	c.chunk_len = len(c.chunk.trans)
	c.current_index = 0

	return nil
}

// Load Next data(1 min after)
func (c *DbSession) LoadNext() (err error) {
	next_time := c.current_end.Add(30 * time.Second) // Load 30 sec after the end
	return c.LoadTime(next_time)
}

// Load Before data(1 min before)
func (c *DbSession) LoadBefore() (err error) {
	next_time := c.current_start.Add(-30 * time.Second) // load 30 sec before the begining
	return c.LoadTime(next_time)
}

func (c *Db) GetTimeChunks() TimeFrames {
	return c.time_chunks
}

// Retrive order book board information from logdb
func (c *DbSession) GetBoard(t time.Time) (bid Board, ask Board, err error) {
	if !c.db.time_chunks.In(t) {
		return nil, nil, fmt.Errorf("out of range time=%s in[%s %s]", t, c.chunk.start_time(), c.chunk.end_time())
	}

	err = c.LoadTime(t)
	if err != nil {
		return nil, nil, err
	}

	bid, ask, err = c.chunk.GetOrderBook(t)

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

func (c *DbSession) ReadTran() (tran Transaction, err error) {
	for {
		if c.chunk_len <= c.current_index {
			err = c.LoadNext()
			if err != nil {
				break
			}
		}
		tran = c.chunk.trans[c.current_index]
		c.current_index += 1

		if tran.Action == UPDATE_BUY || tran.Action == UPDATE_SELL || tran.Action == PARTIAL {
			continue
		}

		if c.select_end <= tran.Time_stamp {
			err = fmt.Errorf("end of select")
			return tran, err
		} else {
			return tran, nil
		}
	}

	return tran, err
}

// Select transaction from session (setup time duration and returns reader inteface)
func (c *DbSession) SelectTrans(s time.Time, e time.Time) (reader TranReader, err error) {
	bs, be, err := check_bounds(c.db.time_chunks, s, e)

	if err != nil {
		return nil, err
	}
	if bs == BOUND_BEFORE || be == BOUND_AFTER || e.Before(s) {
		err = fmt.Errorf("select time is out of chunk %s, %s", s, e)
		return nil, err
	}

	session := *c

	session.select_start = TimeToNsec(s)
	session.select_end = TimeToNsec(e)

	err = session.LoadTime(s)

	return &session, err
}

func (c *DbSession) SelectOhlcv(s time.Time, e time.Time) error {
	//c.check_boud(s, e)

	// TODO: not imlemented
	// ohlcv, _ := c.chunk.ohlcv(s, e)

	return nil
}

func (c *DbSession) ReadOhlcv() (ohlcv Ohlcv, err error) {
	return ohlcv, err
}
