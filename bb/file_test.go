package bb

import (
	"fmt"
	"testing"
)

func TestCreatePath(t *testing.T) {
	dir, file := make_file_name("/base")
	fmt.Println(dir, file)
}

func TestCreateWriter(t *testing.T) {
	fw := Create_writer("/tmp/")
	defer fw.Close()

	fw.Write([]byte("TEST"))
}
