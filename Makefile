GOFILES = $(shell find . -name '*.go')
GOPATH = $(go env GOPATH)
GOBIN = $(GOPATH)/bin

default: benchmark test

.PHONY: test benchmark coverage

test: ;$(info ▷ testing salesforceid)
	@go test -v -coverprofile coverage.out

benchmark: ;$(info ▷ running benchmarks)
	@go test -v -bench=. -benchmem

coverage: test ;$(info ▷ generating coverage.html)
	@go tool cover -html=coverage.out -o coverage.html
