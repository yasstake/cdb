package trans

import (
	"encoding/json"
	"fmt"
	"testing"
)

type PlotlyOhlc struct {
	X     []string `json:"x"`
	Open  []int    `json:"open"`
	High  []int    `json:"high"`
	Low   []int    `json:"low"`
	Close []int    `json:"close"`
	Type  string   `json:"type"`
	Xaxis string   `json:"xaxis"`
	Yaxis string   `json:"yaxis"`
}

func (c *PlotlyOhlc) init() {
	c.Type = "ohlc"
	c.Xaxis = "x"
	c.Yaxis = "y"
}

func (c *PlotlyOhlc) Append(rec Ohlcv) {
	if rec.open == 0 {
		return
	}

	time := DateTime(rec.time)
	c.X = append(c.X, fmt.Sprintf("%04d-%02d-%02d %02d:%02d:%02d",
		time.Year(), int(time.Month()), time.Day(), time.Hour(), time.Minute(), time.Second()))
	c.Open = append(c.Open, rec.open)
	c.High = append(c.High, rec.high)
	c.Low = append(c.Low, rec.low)
	c.Close = append(c.Close, rec.close)
}

func (c *PlotlyOhlc) dump() string {
	s, _ := json.Marshal(c)

	return string(s)
}

func TestDumpPlotly(t *testing.T) {
	var plot PlotlyOhlc
	plot.init()

	o1 := Ohlcv{1000000000, 1, 2, 3, 4, 5, 6, 0}
	o2 := Ohlcv{2000000000, 7, 9, 10, 11, 12, 13, 14}

	plot.Append(o1)
	plot.Append(o2)
	s := plot.dump()
	fmt.Println(s)
}

func TestMakeData(t *testing.T) {
	var db Db
	db.Open("/tmp")

	s := db.time_chunks[0].start
	// session, _ := db.CreateSession(s)

	var c Chunk

	// s_time := DateTime(1613864762187260 * 1000)
	c.LoadTime(s)

	ohlcvs := c.GetOhlcvSec()

	var plot PlotlyOhlc
	plot.init()

	for i := range ohlcvs {
		if ohlcvs[i].open == 0 {
			continue
		}
		plot.Append(ohlcvs[i])
	}
	str, _ := json.Marshal(plot)
	fmt.Println(string(str))
}
