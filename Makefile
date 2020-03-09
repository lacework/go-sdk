default: ci

ci: lint test fmt-check

export GOFLAGS=-mod=vendor

GOLANGCILINTVERSION?=1.23.8
COVERAGEOUT?=coverage.out

prepare: install-tools go-vendor

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

install-tools:
ifeq (, $(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCILINTVERSION)
endif
ifeq (, $(shell which goimports))
	go get golang.org/x/tools/cmd/goimports
endif
