GOLANGCILINTVERSION?=1.50.0
GOIMPORTSVERSION?=v0.1.12
GOXVERSION?=v1.0.1
GOTESTSUMVERSION?=v1.8.2
COVERAGEOUT?=coverage.out
COVERAGEHTML?=coverage.html
GOJUNITOUT?=go-junit.xml
PACKAGENAME?=[[.Component]]
CLINAME?=lacework
GO_LDFLAGS="-X [[.Component]]/internal/version.Version=$(shell cat VERSION) \
            -X [[.Component]]/internal/version.GitSHA=$(shell git rev-parse HEAD) \
            -X [[.Component]]/internal/version..BuildTime=$(shell date +%Y%m%d%H%M%S)"
GOFLAGS=-mod=vendor
CGO_ENABLED?=0
export GOFLAGS GO_LDFLAGS CGO_ENABLED


.PHONY: help
help:
	@echo "-------------------------------------------------------------------"
	@echo " Makefile helper:"
	@echo ""
	@grep -Fh "##" $(MAKEFILE_LIST) | grep -v grep | sed -e 's/\\$$//' | sed -E 's/^([^:]*):.*##(.*)/ \1 -\2/'
	@echo "-------------------------------------------------------------------"

.PHONY: prepare
prepare: install-tools go-vendor ## Initialize the go environment

.PHONY: go-vendor
go-vendor: ## Runs go mod tidy, vendor and verify to cleanup, copy and verify dependencies
	go mod tidy
	go mod vendor
	go mod verify

.PHONY: test
test: prepare ## Run unit tests
	gotestsum -f testname -- -v -cover -run=$(regex) -coverprofile=$(COVERAGEOUT) $(shell go list ./...)

.PHONY: imports-check
imports-check: ## Lists imports issues
	@test -z $(shell goimports -l $(shell go list -f {{.Dir}} ./...))

.PHONY: lint
lint: ## Runs go linter
	golangci-lint run

.PHONY: fmt
fmt: ## Runs and applies go formatting changes
	@gofmt -w -l $(shell go list -f {{.Dir}} ./...)
	@goimports -w -l $(shell go list -f {{.Dir}} ./...)

.PHONY: fmt-check
fmt-check: ## Lists formatting issues
	@test -z $(shell gofmt -l $(shell go list -f {{.Dir}} ./...))

.PHONY: build
build: ## Compiles binary for the running workstation (CDK support)
	go build .

.PHONY: build-cross-platform
build-cross-platform: ## Compiles binaries for all supported platforms
	gox -output="bin/$(PACKAGENAME)-{{.OS}}-{{.Arch}}" \
            -os="linux windows" \
            -arch="amd64 386" \
            -osarch="darwin/amd64 darwin/arm64 linux/arm linux/arm64" \
            -ldflags=$(GO_LDFLAGS)

.PHONY: install-tools
install-tools: ## Install go indirect dependencies
ifeq (, $(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v$(GOLANGCILINTVERSION)
endif
ifeq (, $(shell which goimports))
	GOFLAGS=-mod=readonly go install golang.org/x/tools/cmd/goimports@$(GOIMPORTSVERSION)
endif
ifeq (, $(shell which gox))
	GOFLAGS=-mod=readonly go install github.com/mitchellh/gox@$(GOXVERSION)
endif
ifeq (, $(shell which gotestsum))
	GOFLAGS=-mod=readonly go install gotest.tools/gotestsum@$(GOTESTSUMVERSION)
endif

.PHONY: uninstall-tools
uninstall-tools: ## Uninstall go indirect dependencies
ifneq (, $(shell which golangci-lint))
	rm $(shell go env GOPATH)/bin/golangci-lint
endif
ifneq (, $(shell which goimports))
	rm $(shell go env GOPATH)/bin/goimports
endif
ifneq (, $(shell which gox))
	rm $(shell go env GOPATH)/bin/gox
endif
ifneq (, $(shell which gotestsum))
	rm $(shell go env GOPATH)/bin/gotestsum
endif
