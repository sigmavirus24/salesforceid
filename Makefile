GOFILES = $(shell find . -name '*.go')
GOPATH = $(go env GOPATH)
GOBIN = $(GOPATH)/bin

default: benchmark test

.PHONY: test benchmark

test: ;$(info ▷ testing salesforceid)
	@go test -v

benchmark: ;$(info ▷ running benchmarks)
	@go test -v -bench=. -benchmem
