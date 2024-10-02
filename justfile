#!/usr/bin/env just --justfile

run:
  go run main.go

build:
	go build -o build/ucli main.go

test-server:
	cd test_app && go run server.go
