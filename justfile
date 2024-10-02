#!/usr/bin/env just --justfile

run:
  cd source && go run main.go

build:
	cd source && go build -o ../build/ucli main.go

test-server:
	cd test_app && go run server.go
