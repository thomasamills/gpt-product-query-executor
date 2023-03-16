SHELL := /bin/bash

.PHONY: build
build:
	go mod download
	go mod vendor
	go build -o ./gpt-product-gen