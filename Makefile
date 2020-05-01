default: ci

ci: lint test fmt-check imports-check build-cli-cross-platform

GOLANGCILINTVERSION?=1.23.8
COVERAGEOUT?=coverage.out
PACKAGENAME?=lacework-cli
CLINAME?=lacework
GO_LDFLAGS="-X github.com/lacework/go-sdk/cli/cmd.Version=$(shell cat VERSION) \
            -X github.com/lacework/go-sdk/cli/cmd.GitSHA=$(shell git rev-parse HEAD) \
            -X github.com/lacework/go-sdk/cli/cmd.BuildTime=$(shell date +%Y%m%d%H%M%S)"
GOFLAGS=-mod=vendor
export GOFLAGS GO_LDFLAGS

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
	@test -z $(shell gofmt -l $(shell go list -f {{.Dir}} ./...))

imports-check:
	@test -z $(shell goimports -l $(shell go list -f {{.Dir}} ./...))

build-cli-cross-platform:
	gox -output="bin/$(PACKAGENAME)-{{.OS}}-{{.Arch}}" \
            -os="darwin linux windows" \
            -arch="amd64 386" \
            -ldflags=$(GO_LDFLAGS) \
            github.com/lacework/go-sdk/cli

install-cli: build-cli-cross-platform
ifeq (x86_64, $(shell uname -m))
	mv bin/$(PACKAGENAME)-$(shell uname -s | tr '[:upper:]' '[:lower:]')-amd64 /usr/local/bin/$(CLINAME)
else
	mv bin/$(PACKAGENAME)-$(shell uname -s | tr '[:upper:]' '[:lower:]')-386 /usr/local/bin/$(CLINAME)
endif
	@echo "\nThe lacework cli has been installed at /usr/local/bin"

release-prepare: lint fmt-check imports-check test
	scripts/release.sh prepare

release-do: lint fmt-check imports-check test
	scripts/release.sh release

release-containers:
	scripts/release_containers.sh

install-tools:
ifeq (, $(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCILINTVERSION)
endif
ifeq (, $(shell which goimports))
	go get golang.org/x/tools/cmd/goimports
endif
ifeq (, $(shell which gox))
	go get github.com/mitchellh/gox
endif
