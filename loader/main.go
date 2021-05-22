package main

import (
	"cdb/trans"
	"flag"
	"fmt"
)

func main() {
	var db_dir = flag.String("db_dir", "/tmp/", "db dir")
	var load = flag.Bool("load", false, "load log file to db")

	flag.Parse()

	fmt.Println("Loading logfile to db directory")
	fmt.Println("[DB dir]", *db_dir)
	fmt.Println("[load]", *load)

	trans.SetDbRoot(*db_dir)

	if *load {
		load_files(flag.Args())
	}

	chunks := trans.TimeChunks(trans.DB_ROOT)
	fmt.Println(chunks.ToString())
}

func load_files(files []string) {
	for _, file := range files {
		fmt.Println("[PROCESSING] ", file)
		trans.Load_log(file)
	}
}
