package main

import (
	"cdb/bb"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	writer := bb.Create_writer("/tmp/BB")
	bb.Connect("/tmp/PROCESSSA", writer)
}
