package trans

import "time"

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
func (c *Db) LoadTime(t time.Time) Chunk {
	if c.current_start.Before(t) && t.Before(c.current_end) {
		return c.chunk
	} else {
		c.chunk.load_time(t)
		c.current_start = c.chunk.start_time()
		c.current_end = c.chunk.end_time()
	}
	return c.chunk
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

func (c *Db) SelectTrans(s time.Time, e time.Time) {

}

func (c *Db) SelectOhlcv(s time.Time, e time.Time) {

}
