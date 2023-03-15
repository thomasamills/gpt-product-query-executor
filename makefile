SHELL := /bin/bash

.PHONY: build
build:
	./config.sh
	go mod download
	go mod vendor
	go build -o ./gpt-product-gen
	./gpt-product-gen
