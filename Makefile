.DEFAULT_GOAL := build

.PHONY: lint
lint:
	golangci-lint run

.PHONY: test
test:
	go test ./...

.PHONY: build
build:
	go build github.com/shibataka000/dailyreport/cmd/dailyreport

.PHONY: install
install:
	go install github.com/shibataka000/dailyreport/cmd/dailyreport

.PHONY: clean
clean:
	go clean -testcache
