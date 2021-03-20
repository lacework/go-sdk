default: ci

ci: lint test fmt-check imports-check integration

GOLANGCILINTVERSION?=1.23.8
CIARTIFACTS?=circleci-artifacts
COVERAGEOUT?=coverage.out
COVERAGEHTML?=coverage.html
PACKAGENAME?=lacework-cli
CLINAME?=lacework
# Honeycomb variables
HONEYDATASET?=lacework-cli-dev
# => HONEYAPIKEY should be exported on every developers workstation or else events
#                won't be recorded in Honeycomb. Inside our CI/CD pipeline this
#                secret is set as well as a different dataset for production
GO_LDFLAGS="-X github.com/lacework/go-sdk/cli/cmd.Version=$(shell cat VERSION) \
            -X github.com/lacework/go-sdk/cli/cmd.GitSHA=$(shell git rev-parse HEAD) \
            -X github.com/lacework/go-sdk/cli/cmd.HoneyApiKey=$(HONEYAPIKEY) \
            -X github.com/lacework/go-sdk/cli/cmd.HoneyDataset=$(HONEYDATASET) \
            -X github.com/lacework/go-sdk/cli/cmd.BuildTime=$(shell date +%Y%m%d%H%M%S)"
GOFLAGS=-mod=vendor
CGO_ENABLED?=0
export GOFLAGS GO_LDFLAGS CGO_ENABLED

prepare: install-tools go-vendor

test:
	go test -v -cover -coverprofile=$(COVERAGEOUT) $(shell go list ./... | grep -v integration)

integration: build-cli-cross-platform integration-only

integration-only:
	PATH=$(PWD)/bin:${PATH} go test -v github.com/lacework/go-sdk/integration -timeout 30m

coverage: test
	go tool cover -func=$(COVERAGEOUT)

coverage-html: test
	go tool cover -html=$(COVERAGEOUT)

coverage-ci: test
	mkdir -p $(CIARTIFACTS)
	go tool cover -html=$(COVERAGEOUT) -o "$(CIARTIFACTS)/$(COVERAGEHTML)"

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
            -os="linux windows" \
            -arch="amd64 386" \
            -osarch="darwin/amd64 linux/arm linux/arm64" \
            -ldflags=$(GO_LDFLAGS) \
            github.com/lacework/go-sdk/cli

generate-databox:
	go generate internal/databox/box.go

generate-docs:
	go generate cli/cmd/docs.go

test-resources:
	scripts/prepare_test_resources.sh all

install-cli: build-cli-cross-platform
ifeq (x86_64, $(shell uname -m))
	mv bin/$(PACKAGENAME)-$(shell uname -s | tr '[:upper:]' '[:lower:]')-amd64 /usr/local/bin/$(CLINAME)
else
	mv bin/$(PACKAGENAME)-$(shell uname -s | tr '[:upper:]' '[:lower:]')-386 /usr/local/bin/$(CLINAME)
endif
	@echo "\nThe lacework cli has been installed at /usr/local/bin"

release: lint test fmt-check imports-check build-cli-cross-platform
	scripts/release.sh prepare

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
