package bb

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func make_file_name(base string) (dir string, file string) {
	time := time.Now()

	yyyy := fmt.Sprintf("%04d", time.Year())
	mm := fmt.Sprintf("%02d", time.Month())
	dd := fmt.Sprintf("%02d", time.Day())

	dir = filepath.Join(base, yyyy, mm, dd)
	file = time.Format("2006-01-02-15-04-05") + ".log.gz"

	return dir, file
}

func Create_writer(base string) io.WriteCloser {
	dir, file := make_file_name(base)

	os.MkdirAll(dir, 0777)

	full_path := filepath.Join(dir, file)
	fp, err := os.Create(full_path)

	if err != nil {
		log.Println("CANNOT OPEN FILE for Log Write", full_path)
	}

	wc := io.WriteCloser(fp)
	//	return wc

	gz := gzip.NewWriter(wc)
	return gz
}
