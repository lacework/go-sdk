default: ci

ci: lint test fmt-check

export GOFLAGS=-mod=vendor

GOLANGCILINTVERSION?=1.23.8
COVERAGEOUT?=coverage.out

prepare:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCILINTVERSION)
	go get golang.org/x/tools/cmd/goimports

test:
	go test -v -cover -coverprofile=$(COVERAGEOUT) ./...

coverage: test
	go tool cover -func=$(COVERAGEOUT)

coverage-html: test
	go tool cover -html=$(COVERAGEOUT)

go-vendor:
	go mod tidy
	go mod vendor
	go mod verify

lint:
	golangci-lint run

fmt:
	@gofmt -w -l ./
	@goimports -w -l ./

fmt-check:
	@test -z $(shell gofmt -l ./)
	@test -z $(shell goimports -l ./)

cli:
	go run cli/main.go
