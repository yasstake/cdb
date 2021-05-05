


build:
	-mkdir bin
	go build -o ./bin/load ./loader/main.go
	go build -o ./bin/logger ./logger/main.go	

run:
	go run ./loader/main.go  

bench:
	go test -bench ./ -benchmem

test:
	go test ./trans/*_test.go 

clean:
	rm -rf ./bin

import: build
	./bin/load -load ../LOGFILES/comp.log.gz


ws:
	go build -o ./bin/logger ./logger/main.go


run-ws: 
	go run ./logger/main.go -log_dir /tmp/BB -flag_file /tmp/PROCESSA -exit_wait 1


loc:
	 find . -name \*.go | xargs wc
