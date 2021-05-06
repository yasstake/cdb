package trans

import (
	"compress/gzip"
	"encoding/binary"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Action
const PARTIAL = 1
const UPDATE_SELL = 2
const UPDATE_BUY = 3

// trade
const TRADE_BUY = 4
const TRADE_BUY_LIQUID = 5

const TRADE_SELL = 6
const TRADE_SELL_LIQUID = 7

// Open Interest
// action, time, 0,, volume,
const OPEN_INTEREST = 10

// Open Value
// action, time, 0, volume
const OPEN_VALUE = 11

// Turn Over
// action, time, 0, volume
const TURN_OVER = 12

// Funding Rate
// action, time, 0, volume, next time
const FUNDING_RATE = 20

// Next Funding Rate
// action, time, 0, volume, next time
const PREDICTED_FUNDING_RATE = 21

var DB_ROOT = "/tmp/BITLOG"

const VOLUME_MAG = 1000

func SetDbRoot(path string) {
	DB_ROOT = filepath.Join(path, "BITLOG")
}

const SEC_IN_NS = 1_000_000_000 // ns = sec

// Convert unix time(in ns) to Time object
func DateTime(nsec int64) time.Time {
	t := time.Unix(0, nsec).UTC()
	return t
}

// make log file path from Time
func make_path(time time.Time) (dir, path string) {
	yy := time.Year()
	mm := time.Month()
	dd := time.Day()
	h := time.Hour()
	m := time.Minute()

	dir_name := fmt.Sprintf("%04d-%02d-%02d", yy, mm, dd) + string(os.PathSeparator)
	file_name := fmt.Sprintf("%02d-%02d.log.gz", h, m)

	return dir_name, file_name
}

func make_full_path(time time.Time) (path string) {
	dir_name, file_name := make_path(time)
	return filepath.Join(DB_ROOT, dir_name, file_name)
}

func create_db_file(time time.Time) (db_file io.WriteCloser) {
	dir_name, file_name := make_path(time)

	dir_path := filepath.Join(DB_ROOT, dir_name)
	os.MkdirAll(dir_path, 0777)

	db_path := filepath.Join(dir_path, file_name)
	fw, _ := os.Create(db_path)
	gw := gzip.NewWriter(fw)
	return gw
}

// parse file name and return timestamp in TimeMs
// reverse function of make_path
func file_to_time(file_path string) time.Time {
	// path_exp := `(\d{4})-(\d{2})-(\d{2})` + os.PathSeparator + `(\d{2})-(\d{2}).log.gz$`
	//re := regexp.MustCompile(path_exp)
	re := regexp.MustCompile(`(\d{4})-(\d{2})-(\d{2})/(\d{2})-(\d{2}).log.gz$`)
	res := re.FindStringSubmatch(file_path)

	yy, err := strconv.Atoi(res[1])
	if err != nil {
		log.Println("[ERROR] Wrong time format", file_path)
	}
	mm, err := strconv.Atoi(res[2])
	if err != nil {
		log.Println("[ERROR] Wrong time format", file_path)
	}
	dd, err := strconv.Atoi(res[3])
	if err != nil {
		log.Println("[ERROR] Wrong time format", file_path)
	}
	h, err := strconv.Atoi(res[4])
	if err != nil {
		log.Println("[ERROR] Wrong time format", file_path)
	}
	m, err := strconv.Atoi(res[5])
	if err != nil {
		log.Println("[ERROR] Wrong time format", file_path)
	}

	return time.Date(yy, time.Month(mm), dd, h, m, 0, 0, time.UTC)
}

// Store each transaction data from Bybit Exchange
type Transaction struct {
	Action     int8
	Time_stamp int64
	Price      int32
	Volume     int64
	NextTime   int64
}

func (c *Transaction) info_string() (result string) {
	result += DateTime(c.Time_stamp).String()
	result += "{Action:" + strconv.Itoa(int(c.Action)) + "}"
	result += "{Price:" + strconv.Itoa(int(c.Price)) + "}"
	result += "{vol:" + strconv.Itoa(int(c.Volume)) + "}"
	result += "{next_time:" + strconv.Itoa(int(c.NextTime)) + "}"

	return result
}

func (t *Transaction) save(stream io.WriteCloser) {
	binary.Write(stream, binary.LittleEndian, t)
}

func (t *Transaction) load(stream io.ReadCloser) Transaction {
	binary.Read(stream, binary.LittleEndian, t)
	return *t
}

// Array of transaction
type Transactions []Transaction

func (t *Transactions) init() {
	*t = make(Transactions, 0, 1000)
}

func (t Transactions) save(stream io.WriteCloser) {
	length := int32(len(t))
	binary.Write(stream, binary.LittleEndian, &length)
	for i := 0; i < int(length); i++ {
		t[i].save(stream)
	}
}

func (t *Transactions) load(stream io.ReadCloser) Transactions {
	var length int32
	binary.Read(stream, binary.LittleEndian, &length)

	re := make(Transactions, length)

	for i := 0; i < int(length); i++ {
		re[i] = re[i].load(stream)
	}
	*t = re

	return *t
}

// Store order book
type Board map[int]int

func (c *Board) init() {
	*c = make(Board)
}

func (board *Board) set(price int, volume int) {
	if volume == 0 {
		delete(*board, price)
	} else {
		(*board)[price] = volume
	}
}

func (board *Board) copy() Board {
	copy_board := make(Board)

	for key, value := range *board {
		copy_board[key] = value
	}

	return copy_board
}

type boardBuf struct {
	Price uint32
	Vol   uint32
}

func (board *Board) save(stream io.WriteCloser) {
	length := uint16(len(*board))
	binary.Write(stream, binary.LittleEndian, &length)

	var buf boardBuf

	for price := range *board {
		buf.Price = uint32(price)
		buf.Vol = uint32((*board)[price])
		binary.Write(stream, binary.LittleEndian, &buf)
	}
}

func (board *Board) load(stream io.ReadCloser) Board {
	var len16 uint16
	binary.Read(stream, binary.LittleEndian, &len16)
	len := int(len16)

	var buf boardBuf

	board.init()

	for i := 0; i < len; i++ {
		binary.Read(stream, binary.LittleEndian, &buf)
		(*board)[int(buf.Price)] = int(buf.Vol)
	}

	return *board
}

func (board *Board) depth() int {
	return len(*board)
}

type Order []struct {
	Price  int
	Volume int
}

func (c Order) Len() int {
	return len(c)
}

func (c Order) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func (c Order) Less(i, j int) bool {
	return c[i].Price < c[j].Price
}

func (board *Board) Sort(asc bool) (orders Order) {
	orders = make(Order, len(*board))

	i := 0
	for price, volume := range *board {
		orders[i].Price = price
		orders[i].Volume = volume
		i++
	}

	if asc {
		sort.Slice(orders, func(i, j int) bool { return orders[i].Price < orders[j].Price })
	} else {
		sort.Slice(orders, func(i, j int) bool { return orders[i].Price > orders[j].Price })
	}

	return orders
}

type Chunk struct {
	bid_board Board
	ask_board Board
	trans     Transactions

	/* Ensure one funding record in one chunk we must buffer latest info
	funding_rate         Transaction
	funding_rate_predict Transaction
	*/
}

func (c Chunk) info_string() string {
	bid_len := len(c.bid_board)
	ask_len := len(c.ask_board)
	trans_len := len(c.trans)

	start := c.trans[0]
	end := c.trans[trans_len-1]

	result :=
		"Start->" + start.info_string() +
			" End->" + end.info_string() +
			" BID->" + strconv.Itoa(bid_len) +
			" ASK->" + strconv.Itoa(ask_len) +
			" Trans->" + strconv.Itoa(trans_len)

	return result
}

func (c *Chunk) start_time() time.Time {
	return DateTime(c.trans[0].Time_stamp)
}

func (c *Chunk) end_time() time.Time {
	i := len(c.trans) - 1

	return DateTime(c.trans[i].Time_stamp)
}

func (c *Chunk) init() {
	c.bid_board.init()
	c.ask_board.init()
	c.trans.init()
}

func (c *Chunk) Append(r Transaction) {
	c.trans = append(c.trans, r)
}

func (c *Chunk) dump() {
	time := DateTime(c.trans[0].Time_stamp)
	stream := create_db_file(time)
	defer stream.Close()

	c.bid_board.save(stream)
	c.ask_board.save(stream)

	c.trans.save(stream)
}

func (c *Chunk) load_file(path string) error {
	stream, err := os.Open(path)
	if err != nil {
		log.Printf("cannot open file%s %s", path, err)
		return err
	}

	gzip_reader, _ := gzip.NewReader(stream)
	defer gzip_reader.Close()

	c.bid_board.load(gzip_reader)
	c.ask_board.load(gzip_reader)
	c.trans.load(gzip_reader)

	return nil
}

func (c *Chunk) LoadTime(time time.Time) error {
	path := make_full_path(time)

	return c.load_file(path)
}

func (c *Chunk) GetTran() (result Transactions) {
	for i := range c.trans {
		action := c.trans[i].Action
		if action == UPDATE_BUY || action == UPDATE_SELL || action == PARTIAL {
			continue
		}

		result = append(result, c.trans[i])
	}

	return result
}

type Ohlcv struct {
	time     int64
	open     int
	high     int
	low      int
	close    int
	buy_vol  int
	sell_vol int
	vol      int
}

func (c *Ohlcv) init() {
	c.time = 0
	c.open = 0
	c.high = 0
	c.low = 0
	c.close = 0
	c.buy_vol = 0
	c.sell_vol = 0
	c.vol = 0
}

// Merge two ohlcv value and retun merged one
// when add value is {0, 0, 0, 0, 0} then return original
func (c Ohlcv) add(ohlcv Ohlcv) (result Ohlcv) {
	if ohlcv.open == 0 && ohlcv.close == 0 {
		return c
	}

	result.time = ohlcv.time
	result.open = c.open

	if c.high < ohlcv.high {
		result.high = ohlcv.high
	} else {
		result.high = c.high
	}

	if c.low < ohlcv.low {
		result.low = c.low
	} else {
		result.low = ohlcv.low
	}

	result.close = ohlcv.close

	result.buy_vol += c.buy_vol + ohlcv.buy_vol
	result.sell_vol += c.sell_vol + ohlcv.sell_vol

	result.vol = result.buy_vol + result.sell_vol

	return result
}

func (c *Ohlcv) buy(time int64, price int, volume int) {
	c.sell_buy(time, price, volume, true)
}

func (c *Ohlcv) sell(time int64, price int, volume int) {
	c.sell_buy(time, price, volume, false)
}

func (c *Ohlcv) sell_buy(time int64, price int, volume int, buy bool) {
	c.time = time

	if c.open == 0 {
		c.open = price
	}

	if c.high < price || c.high == 0 {
		c.high = price
	}

	if price < c.low || c.low == 0 {
		c.low = price
	}

	c.close = price

	if buy {
		c.buy_vol += volume
		c.vol += volume
	} else {
		c.sell_vol += volume
		c.vol += volume
	}
}

// Calculate OHLCV in chunk within time.
//  returns OHLCV and raw record within transaction
func (c *Chunk) GetOhlcv(from time.Time, end time.Time) (result Ohlcv, num_record int) {
	result.init()

	for i := range c.trans {

		time_stamp := DateTime(c.trans[i].Time_stamp)

		if time_stamp.Before(from) {
			continue
		}

		if end.Before(time_stamp) {
			break
		}

		action := c.trans[i].Action

		if action == TRADE_BUY || action == TRADE_BUY_LIQUID {
			num_record += 1
			result.buy(c.trans[i].Time_stamp, int(c.trans[i].Price), int(c.trans[i].Volume))
		} else if action == TRADE_SELL || action == TRADE_SELL_LIQUID {
			num_record += 1
			result.sell(c.trans[i].Time_stamp, int(c.trans[i].Price), int(c.trans[i].Volume))
		}
	}

	if num_record == 0 {
		log.Printf("no data in time frame　%s %s", from, end)
	}

	return result, num_record
}

func (c *Chunk) GetOhlcvSec() (result []Ohlcv) {
	start_time := int64((c.trans[0].Time_stamp+SEC_IN_NS/10)/SEC_IN_NS) * SEC_IN_NS
	current_end := start_time + SEC_IN_NS

	var ohlcv Ohlcv
	ohlcv.init()
	trans_len := len(c.trans)

	for i := 0; i < trans_len; i++ {
		time_stamp := c.trans[i].Time_stamp

		if time_stamp <= current_end || trans_len-i < 100 {
			action := c.trans[i].Action

			if action == TRADE_BUY || action == TRADE_BUY_LIQUID {
				ohlcv.buy(c.trans[i].Time_stamp, int(c.trans[i].Price), int(c.trans[i].Volume))
			} else if action == TRADE_SELL || action == TRADE_SELL_LIQUID {
				ohlcv.sell(c.trans[i].Time_stamp, int(c.trans[i].Price), int(c.trans[i].Volume))
			}
		} else {
			current_end += SEC_IN_NS
			result = append(result, ohlcv)
			ohlcv.init()
		}
	}
	if ohlcv.open != 0 {
		result = append(result, ohlcv)
	}

	return result
}

// Apply update transaction data until the time and return bords.
// TODO: cache result for performance improbement
func (c *Chunk) GetOrderBook(time time.Time) (bid, ask Board, err error) {
	bid = nil
	ask = nil

	trans_len := len(c.trans)
	t := time.UnixNano()

	// if chunk does not have data
	if trans_len == 0 {
		return nil, nil, fmt.Errorf("chunk have not data %s", time)
	}

	// check out of chunk time frame(with 100ms allowance)
	if t < c.trans[0].Time_stamp-SEC_IN_NS/10 ||
		c.trans[trans_len-1].Time_stamp+SEC_IN_NS/10 < t {
		return nil, nil, fmt.Errorf("select outside of chunk timeframe, %s", time)
	}

	// first setup initial
	bid = c.bid_board.copy()
	ask = c.ask_board.copy()

	// apply transactions until time
	for i := 0; i < trans_len; i++ {
		time_stamp := c.trans[i].Time_stamp
		// TODO: Consider action (execute and board update)
		// action := c.trans[i].Action

		// exceed time
		if t < time_stamp {
			break
		}

		action := int(c.trans[i].Action)

		if action == UPDATE_BUY {
			price := int(c.trans[i].Price)
			volume := int(c.trans[i].Volume)
			bid.set(price, volume)
		} else if action == UPDATE_SELL {
			price := int(c.trans[i].Price)
			volume := int(c.trans[i].Volume)
			ask.set(price, volume)
		}
	}

	return bid, ask, nil
}

func (c *Chunk) open_interest(time time.Time) (oi int, err bool) {
	trans_len := len(c.trans)
	t := time.UnixNano()

	// if chunk does not have data
	if trans_len == 0 {
		return 0, true
	}

	// check out of chunk time frame(with 100ms allowance)
	if t < c.trans[0].Time_stamp-SEC_IN_NS/10 ||
		c.trans[trans_len-1].Time_stamp+SEC_IN_NS/10 < t {
		return 0, true
	}

	for i := 0; i < trans_len; i++ {
		time_stamp := c.trans[i].Time_stamp
		// TODO: Consider action (execute and board update)
		// action := c.trans[i].Action

		// exceed time
		if t < time_stamp {
			break
		}

		action := int(c.trans[i].Action)

		if action == OPEN_INTEREST {
			fmt.Println(c.trans[i].info_string())
			return int(c.trans[i].Volume), false
		}
	}

	return 0, true
}

func Load_log(file string) (chunk Chunk) {
	f, err := os.Open(file)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	compress := strings.HasSuffix(file, ".gz")
	var r *csv.Reader
	if compress {
		gzipfile, _ := gzip.NewReader(f)
		r = csv.NewReader(gzipfile)
	} else {
		r = csv.NewReader(f)
	}
	r.FieldsPerRecord = -1 // ignore feild number varies

	var record Transaction

	last_min := int(-1)

	var bid_board Board
	bid_board.init()
	var ask_board Board
	ask_board.init()

	chunk.init()

	var current_time int64
	var current_price int32

	for {
		row, err := r.Read()
		if err == io.EOF {
			fmt.Println("[PROCESS DONE]")
			break
		}
		if err != nil {
			fmt.Println("[FILE READ ERROR]", err)
			break
		}
		for i, v := range row {
			switch i {
			case 0: // Action
				r, _ := strconv.Atoi(v)
				record.Action = int8(r)
			case 1: // Time(us)
				t, _ := strconv.ParseInt(v, 10, 64)
				record.Time_stamp = (t + current_time) * 1_000_000 // convert to ns
				current_time = t + current_time

			case 2: // Price
				r, _ := strconv.Atoi(v)
				price := int32(r)
				record.Price = price + current_price
				current_price = price + current_price
			case 3: // volume
				// TODO: FIX omit under floating point
				r, _ := strconv.ParseInt(v, 10, 64)
				record.Volume = r
			case 4: // Time Info
				t, _ := strconv.ParseInt(v, 10, 64)
				record.NextTime = t
			}
		}

		if record.Action == PARTIAL {
			bid_board.init()
			ask_board.init()
		} else if record.Action == UPDATE_BUY || record.Action == UPDATE_SELL {
			time := DateTime(record.Time_stamp)
			min := time.Minute()
			sec := time.Second()

			if min != last_min {
				if sec <= 1 {
					last_min = min
					tr_len := len(chunk.trans)
					if 100 < tr_len {
						duration := chunk.trans[tr_len-1].Time_stamp - chunk.trans[0].Time_stamp

						if 30*1000000 <= duration {
							chunk.dump()
							fmt.Println("DUMP", chunk.info_string())
						}
					}
				}

				chunk.bid_board = bid_board.copy() // CopyBuffer
				chunk.ask_board = ask_board.copy()
				chunk.trans.init()
			}

			if record.Action == UPDATE_BUY {
				bid_board.set(int(record.Price), int(record.Volume))
			} else if record.Action == UPDATE_SELL {
				ask_board.set(int(record.Price), int(record.Volume))
			} else {
				log.Fatal("Unknown action")
			}
		}

		chunk.Append(record)
	}

	return chunk
}
