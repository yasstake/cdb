package main

import (
	"cdb/bb"
	"flag"
)

func main() {
	var log_dir = flag.String("log_dir", "/tmp/BB", "log store directory")
	var flag_file = flag.String("flag_file", "", "flag file name, if not specified no flag file used.")
	var exit_wait = flag.Int("exit_wait", 0, "Exit wait minute, when terminated by peer process")

	flag.Parse()

	writer := bb.Create_writer(*log_dir)
	bb.Connect(*flag_file, writer, *exit_wait)
}
