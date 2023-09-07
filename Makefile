# PROJECT SETTINGS
_PROJECT_DIRECTORY = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
_GOARCH = "amd64"

ifeq ("$(shell uname -m)", "arm64")
	_GOARCH = "arm64"
endif

.PHONY: prepare-ci

prepare-ci:
	@go install gotest.tools/gotestsum@latest

.PHONY: generate

generate:
	@echo "Generating Go stuff..."
	@go generate ./...
	@echo "Done"

.PHONY: format

format: fmt fumpt imports gci

fmt:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofmt -w -s {}'

fumpt:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofumpt -w -extra {}'

imports:
	@goimports -v -w -e -local github.com/unconditionalday main.go
	@goimports -v -w -e -local github.com/unconditionalday cmd/
	@goimports -v -w -e -local github.com/unconditionalday internal/

gci:
	@find . -name "*.go" -type f -not -path '*/vendor/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting imports for {}.." && \
	gci write --skip-generated -s standard,default,"prefix(github.com/unconditionalday)" {}'

.PHONY: test-unit test-integration

test-unit:
	@gotestsum --no-color=false -- -tags=unit ./...

test-integration:
	@gotestsum --no-color=false -- -tags=integration ./...

.PHONY: build

build:
	@go build --tags=release -o ${_PROJECT_DIRECTORY}/bin/unconditional-server

# Helpers
check-variable-%: # detection of undefined variables.
	@[[ "${${*}}" ]] || (echo '*** Please define variable `${*}` ***' && exit 1)
