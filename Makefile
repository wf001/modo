dir = generated/$(shell date +%s)

build:
	mkdir -p $(dir)
	go build ./cmd/modo

build-debug:
	mkdir -p $(dir)
	go build -a -p 1 -x -work -o ./$(dir)/modo ./cmd/modo

clean:
	rm -rf generated
	rm -rf modo


test-all:
	make test-compiler
	make test-go

test-compiler:
	./script/test-lite.sh
	./script/test-full.sh

test-full-compiler:
	./script/test-full.sh

test-lite-compiler:
	./script/test-lite.sh

test-go:
	go test -v ./... | tc
