export GO111MODULE = on

.PHONY: test test-cover

# for test
test:
	go test -race -cover ./...

test-all:
	go test -race -cover -tags brotli ./...

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench:
	go test -tags brotli -bench=. ./...

release:
	go mod tidy
