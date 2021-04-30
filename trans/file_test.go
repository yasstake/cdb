package trans

import (
	"fmt"
	"testing"
	"time"
)

func TestFileList(t *testing.T) {
	files := file_list(DB_ROOT)
	fmt.Println(files)
}

func TestInTimeFrame(t *testing.T) {
	t1 := date_time(time.Hour.Nanoseconds())
	t2 := date_time(time.Hour.Nanoseconds() * 2)

	t3 := date_time(time.Hour.Nanoseconds() * 3)
	t4 := date_time(time.Hour.Nanoseconds()*3 + 1)
	t5 := date_time(time.Hour.Nanoseconds() * 4)

	frame1 := TimeFrame{t1, t2}
	frame2 := TimeFrame{t3, t5}

	if frame1.In(t1) != true {
		t.Error()
	}

	if frame1.In(t2) != false {
		t.Error()
	}

	if frame1.In(t3) != false {
		t.Error(t1, t2, t3)
	}

	if frame2.In(t4) != true {
		t.Error(t3, t5, t4)
	}
}

func TestInTimeFrames(t *testing.T) {
	t1 := date_time(time.Hour.Nanoseconds())
	t2 := date_time(time.Hour.Nanoseconds() * 2)

	t3 := date_time(time.Hour.Nanoseconds() * 4)
	t4 := date_time(time.Hour.Nanoseconds()*4 + 1)
	t5 := date_time(time.Hour.Nanoseconds() * 5)

	frame1 := TimeFrame{t1, t2}
	frame2 := TimeFrame{t3, t5}

	frame := TimeFrames{frame1, frame2}

	if frame.In(t1) != true {
		t.Error()
	}
	if frame.In(t2) != false {
		t.Error()
	}
	if frame.In(t3) != true {
		t.Error()
	}
	if frame.In(t4) != true {
		t.Error()
	}
	if frame.In(t5) != false {
		t.Error()
	}
}

//   t <  Frame.s Frame.e

func TestBefore(t *testing.T) {
	t1 := date_time(time.Hour.Nanoseconds())
	t2 := date_time(time.Hour.Nanoseconds() * 2)

	t3 := date_time(time.Hour.Nanoseconds() * 4)
	// t4 := date_time(time.Hour.Nanoseconds()*4 + 1)
	t5 := date_time(time.Hour.Nanoseconds() * 5)

	frame1 := TimeFrame{t1, t2}
	frame2 := TimeFrame{t3, t5}

	frame := TimeFrames{frame1, frame2}

	if frame.Before(t2) != true {
		t.Error(t1, t2)
	}

}

func TestTimeChunks(t *testing.T) {
	chunks := Time_chunks(DB_ROOT)
	fmt.Println(chunks[0].start, chunks[0].end)
	fmt.Println(chunks.ToString())
}
