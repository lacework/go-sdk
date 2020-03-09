default: ci

ci: test

export GOFLAGS=-mod=vendor

test:
	go test -v -cover ./...
