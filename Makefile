.PHONY: test test-cover bench lint

# for test
test:
	go test -race -cover ./...

test-cover:
	go test -race -coverprofile=test.out ./... && go tool cover --html=test.out

bench:
	go test -bench=. -run='^$$' ./...

lint:
	golangci-lint run
