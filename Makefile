dir = generated/$(shell date +%s)
host-target := $(shell grep -s 'arch:' ./cmd/modo/version.yaml | cut -d' ' -f2| cut -d- -f1,2)
version := $(shell grep -s 'version:' ./cmd/modo/version.yaml | cut -d' ' -f2 )


init:
	./script/gen-version-yaml.sh

build:
	go build -ldflags="$(shell go run ./script/gen_ldflag.go)" -o modo ./cmd/modo


clean:
	rm -rf generated
	rm -rf modo
	rm -rf dist
	rm -rf cmd/modo/version.yaml

build-debug:
	mkdir -p $(dir)
	go build -a -p 1 -x -work -o ./$(dir)/modo ./cmd/modo

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


dist:
	mkdir -p dist/modo/bin
	cp modo dist/modo/bin/
	cp cmd/modo/version.yaml dist/modo/
	tar -czf modo"$(version)"."$(host-target)".tar.gz -C dist modo
