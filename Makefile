


build:
	-mkdir bin
	go build -o ./bin/load ./loader/main.go

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
	go run ./logger/main.go
