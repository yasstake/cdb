package main

import (
	"cdb/bb"
	"flag"
	"log"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	bb.Connect()
}
