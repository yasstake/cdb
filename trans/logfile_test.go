package trans

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"
	"time"
)

var t1, t2, t3, t4, t5, t6 Transaction
var tr Transactions
var bd Board

func init2() {
	t1 = Transaction{1, 2, 3, 4, 7}
	t2 = Transaction{2, 2, 3, 4, 9}
	t3 = Transaction{3, 2, 3, 4, 0}
	t4 = Transaction{4, 2, 3, 4, 0}
	t5 = Transaction{5, 2, 3, 4, 0}
	t6 = Transaction{6, 2, 3, 4, 0}
	tr = Transactions{t1, t2, t3, t4, t5, t6}
	bd.init()
	bd.set(1, 1)
	bd.set(2, 1)
	bd.set(3, 1)
}

func init() {
	init2()
}

func create_file(path string) (stream io.WriteCloser) {
	stream, _ = os.Create(path)
	return stream
}

func open_file(path string) (stream io.ReadCloser) {
	stream, _ = os.Open(path)
	return stream
}

func TestDateTime(t *testing.T) {
	time := DateTime(1613864752487134000)
	fmt.Println(time)
	time = DateTime(1613864752587206)
	fmt.Println(time)

}

func TestMakePath(t *testing.T) {
	time := DateTime(1613864752487134000)
	dir, path := make_path(time)
	fmt.Println(dir, path)
}

func TestMakeDBFile(t *testing.T) {
	time := DateTime(1613864752487134000)
	fw := create_db_file(time)
	fmt.Println(fw)
}

func TestParseFilePath(t *testing.T) {
	path := "sfdasdfasdfas2021-02-20/23-45.log.gz"
	time := file_to_time(path)
	fmt.Println(path, time)
}

func TestTransactionInfoString(t *testing.T) {
	fmt.Println(t1.info_string())
}

func TestTransactionInit(t *testing.T) {
	fmt.Println(tr)
	tr.init()
	fmt.Println(tr)
	init2()
	fmt.Println(tr)
}

func TestTransactionSaveLoad(t *testing.T) {
	f := create_file("/tmp/savedata.bin")
	t1.save(f)

	var r1 Transaction
	f1 := open_file("/tmp/savedata.bin")
	r1.load(f1)

	if t1 != r1 {
		t.Error("Dosenot match", t1, r1)
	}
	fmt.Println(t1, r1)
}

func TestTransactionsSaveLoad(t *testing.T) {
	f := create_file("/tmp/savedata2.bin")
	tr.save(f)

	var r1 Transactions
	f1 := open_file("/tmp/savedata2.bin")
	r1.load(f1)

	if tr[1] != r1[1] {
		t.Error("Dosenot match", tr[1], r1[1])
	}
	fmt.Println(tr, r1)
}

func TestInitBoard(t *testing.T) {
	fmt.Println(bd)
	fmt.Println(bd.depth())
	if bd.depth() != 3 {
		t.Error("fail depth count", bd)
	}
	bd.init()
	fmt.Println(bd)
	fmt.Println(bd.depth())
	if bd.depth() != 0 {
		t.Error("fail to init", bd)
	}
	init2()
}

func TestSaveAndLoadBoard(t *testing.T) {
	f := create_file("/tmp/savedata3.bin")
	bd.save(f)

	var r1 Board
	f1 := open_file("/tmp/savedata3.bin")
	r1.load(f1)

	fmt.Println(bd, r1)

	if !reflect.DeepEqual(bd, r1) {
		t.Error("does not match", bd, r1)
	}
}

func TestLoadTime(t *testing.T) {
	var c Chunk

	Open()

	time := database.time_chunks[0].start
	fmt.Println(time.String())
	c.LoadTime(time)
	fmt.Println(c.info_string())
}

func TestGetTransaction(t *testing.T) {
	var c Chunk

	Open()

	time := database.time_chunks[0].start
	fmt.Println(time.String())
	c.LoadTime(time)
	tran := c.GetTran()

	fmt.Println(tran)
}

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

func TestLoadAndOhlcvSec(t *testing.T) {
	var c Chunk
	s_time := DateTime(1613864762187260 * 1000)
	c.LoadTime(s_time)
	ohlcvs := c.GetOhlcvSec()
	fmt.Println(ohlcvs, len(ohlcvs))
}

func TestOhlcv(t *testing.T) {
	o1 := Ohlcv{100, 1, 2, 3, 4, 5, 6, 0}
	o2 := Ohlcv{200, 7, 8, 9, 10, 11, 12, 0}
	target := Ohlcv{200, 1, 8, 3, 10, 16, 18, 34}

	o3 := o1.add(o2)

	if o3 != target {
		t.Error("Does not match", o3, target)
	}
	fmt.Println(o3)
	var o4 Ohlcv
	fmt.Println(o4)
}

func TestOhlcv2(t *testing.T) {
	var o1 Ohlcv
	o1.init()

	o1.sell(1000, 100, 1)
	fmt.Println(o1)

	o1.buy(1000, 101, 2)
	fmt.Println(o1)

	o1.buy(1000, 101, 2)
	fmt.Println(o1)

	o1.sell(1000, 100, 1)
	fmt.Println(o1)
}

func TestBoardLoad(t *testing.T) {
	time := 1613864762187260 * 1000
	s_time := DateTime(int64(time))

	var c Chunk
	c.LoadTime(s_time)

	// if out of chunk err(before)
	book_time := DateTime(int64(time - 5*SEC_IN_NS))
	bit, ask, err := c.GetOrderBook(book_time)
	if err != nil {
		t.Error("must be null", bit, ask, err)
	}

	book_time = DateTime(int64(time + 61*SEC_IN_NS))
	bit, ask, err = c.GetOrderBook(book_time)
	if err != nil {
		t.Error("must be null", bit, ask, err)
	}

	book_time = DateTime(int64(time + 5*SEC_IN_NS))
	bit, ask, err = c.GetOrderBook(book_time)
	fmt.Println(bit, ask, err)

	book_time = DateTime(int64(time + 10*SEC_IN_NS))
	bit, ask, err = c.GetOrderBook(book_time)
	fmt.Println(bit, ask, err)
}

func TestOpenInterest(t *testing.T) {
	time := 1613864762187260 * 1000
	s_time := DateTime(int64(time + 5*SEC_IN_NS))

	var c Chunk
	c.LoadTime(s_time)

	oi, err := c.open_interest(s_time)
	fmt.Println(oi, err)
}

func TestLoadLog(t *testing.T) {
	tr := Load_log("../../../../DATA/bb2.log")

	fmt.Println(tr)
	fmt.Println(tr.info_string())
}

func TestLoadLogBig(t *testing.T) {
	Load_log("../../../../DATA/BB2-2021-02-20T23-45-52.008914Z.log.gz")
}

func BenchmarkLoadLogBig(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Load_log("../../../../DATA/BB2-2021-02-20T23-45-52.008914Z.log.gz")
	}
}

func TestInitTransactions(t *testing.T) {
	trans := Transactions{t1, t2, t3}
	fmt.Println(trans)

	trans.init()
	fmt.Println(trans)
}

func TestOrderBookSort(t *testing.T) {
	fmt.Println(bd)

	order := bd.Sort(true)
	fmt.Println(order)

	order = bd.Sort(false)
	fmt.Println(order)
}
