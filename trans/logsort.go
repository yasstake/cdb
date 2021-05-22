package trans

import (
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const TIME_GRAD = 1_000_000 // 1_000_000[ns]=ms

func CsvWriteToFile(data TransactionSlice, file_name string) {
	fw, err := os.Create(file_name)
	if err != nil {
		log.Println(err)
	}

	CsvWrite(data, fw)
}

func CsvWrite(data TransactionSlice, stream io.Writer) {
	var current_time int64
	var current_price int32

	for i := range data {
		t := int64(data[i].Time_stamp / TIME_GRAD)
		r := fmt.Sprintf("%d,%d,%d,%d,%d\n",
			data[i].Action, t-current_time, data[i].Price-current_price, data[i].Volume, data[i].OtherInfo)

		stream.Write([]byte(r))
		current_time = t
		current_price = data[i].Price
	}
}

func LogLoad(from_file string) (result TransactionSlice) {
	f, err := os.Open(from_file)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.Println("loading... file=", from_file)

	compress := strings.HasSuffix(from_file, ".gz")
	var r *csv.Reader
	if compress {
		gzipfile, _ := gzip.NewReader(f)
		r = csv.NewReader(gzipfile)
	} else {
		r = csv.NewReader(f)
	}
	r.FieldsPerRecord = -1 // ignore feild number varies

	var record Transaction

	var current_time int64
	var current_price int32
	partial := false

	for {
		row, err := r.Read()
		if err == io.EOF {
			start_time := DateTime(result[0].Time_stamp)
			end_time := DateTime(result[len(result)-1].Time_stamp)
			log.Println("load done from=", start_time, "  end=", end_time)
			break
		}
		if err != nil {
			fmt.Println("[FILE READ ERROR]", err)
			break
		}

		for i, v := range row {
			switch i {
			case 0: // Action
				r, err := strconv.Atoi(v)
				if err != nil {
					log.Println("[ACTION]", err, v)
				}
				record.Action = int8(r)
				if record.Action == PARTIAL {
					partial = true
				}
			case 1: // Time(us)
				t, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					log.Println("[TIMESTAMP]", err, v)
				}
				record.Time_stamp = (t + current_time) * TIME_GRAD // convert to ns
				current_time = t + current_time

			case 2: // Price
				r, err := strconv.Atoi(v)
				if err != nil {
					log.Println("[PRICE]", err, v)
				}
				price := int32(r)
				record.Price = price + current_price
				current_price = price + current_price
			case 3: // volume
				// TODO: FIX omit under floating point
				r, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					log.Println("[VOL]", err, v)
				}

				record.Volume = r
			case 4: // Time Info
				if v == "" {
					record.OtherInfo = 0
				} else {
					t, err := strconv.ParseInt(v, 10, 64)
					if err != nil {
						log.Println("[TIMEINFO]", err, v)
					}

					record.OtherInfo = t
				}
			}
		}

		// ignore messages before partial message comes
		if partial {
			result = append(result, record)
		}
	}

	return result
}
