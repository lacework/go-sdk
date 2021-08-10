default: ci

ci: lint test fmt-check imports-check integration

# Tooling versions
GOLANGCILINTVERSION?=1.23.8
GOIMPORTSVERSION?=v0.1.2
GOXVERSION?=v1.0.1
GOTESTSUMVERSION?=v1.6.4

CIARTIFACTS?=ci-artifacts
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

# CI variables
CI_V2_ACCOUNT?=customerdemo
export CI_V2_ACCOUNT

prepare: install-tools go-vendor

test: prepare
	gotestsum -f testname -- -v -cover -coverprofile=$(COVERAGEOUT) $(shell go list ./... | grep -v integration)

integration: build-cli-cross-platform integration-only

integration-only:
	PATH=$(PWD)/bin:${PATH} go test -v github.com/lacework/go-sdk/integration -timeout 30m -tags="\
		account \
		agent_token \
		compliance \
		configure \
		event \
		help \
		integration \
		migration \
		policy-disabled \
		query-disabled \
		version \
		vulnerability"

integration-lql: build-cli-cross-platform integration-lql-only

integration-lql-only:
	PATH=$(PWD)/bin:${PATH} go test -v github.com/lacework/go-sdk/integration -timeout 30m -tags="query"

integration-policy: build-cli-cross-platform integration-policy-only

integration-policy-only:
	PATH=$(PWD)/bin:${PATH} go test -v github.com/lacework/go-sdk/integration -timeout 30m -tags="policy"

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
            -osarch="darwin/amd64 darwin/arm64 linux/arm linux/arm64" \
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
	go get golang.org/x/tools/cmd/goimports@$(GOIMPORTSVERSION)
endif
ifeq (, $(shell which gox))
	go get github.com/mitchellh/gox@$(GOXVERSION)
endif
ifeq (, $(shell which gotestsum))
	go get gotest.tools/gotestsum@$(GOTESTSUMVERSION)
endif

git-env:
	scripts/git_env.sh

vagrant-macos-up: build-cli-cross-platform
	$(call run_vagrant,macos-sierra,up)
vagrant-macos-login: build-cli-cross-platform
	$(call run_vagrant,macos-sierra,ssh)
vagrant-macos-destroy:
	$(call run_vagrant,macos-sierra,destroy -f)

vagrant-linux-up: build-cli-cross-platform
	$(call run_vagrant,ubuntu-1804,up)
vagrant-linux-login: build-cli-cross-platform
	$(call run_vagrant,ubuntu-1804,ssh)
vagrant-linux-destroy:
	$(call run_vagrant,ubuntu-1804,destroy -f)

vagrant-windows-up: build-cli-cross-platform
	$(call run_vagrant,windows-10,up)
vagrant-windows-login: build-cli-cross-platform
	$(call run_vagrant,windows-10,powershell)
vagrant-windows-destroy:
	$(call run_vagrant,windows-10,destroy -f)

define run_vagrant
	cd cli/vagrant/${1}; vagrant ${2}
endef
