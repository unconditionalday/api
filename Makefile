# PROJECT SETTINGS
_PROJECT_DIRECTORY = $(dir $(realpath $(firstword $(MAKEFILE_LIST))))
_GOARCH = "amd64"

GOOSE_VER = "v3.17.0"

ifeq ("$(shell uname -m)", "arm64")
	_GOARCH = "arm64"
endif

.PHONY: install-tools

install-tools:
	@go install gotest.tools/gotestsum@latest
	@go install github.com/pressly/goose/v3/cmd/goose@${GOOSE_VER}

.PHONY: generate

generate:
	@echo "Generating Go stuff..."
	@go generate ./...
	@echo "Done"

.PHONY: format

format: fmt fumpt imports gci

fmt:
	@find . -name "*.go" -type f -not -path '*/api/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofmt -w -s {}'

fumpt:
	@find . -name "*.go" -type f -not -path '*/api/*' \
	| sed 's/^\.\///g' \
	| xargs -I {} sh -c 'echo "formatting {}.." && gofumpt -w -extra {}'

imports:
	@goimports -v -w -e -local github.com/unconditionalday main.go
	@goimports -v -w -e -local github.com/unconditionalday cmd/
	@goimports -v -w -e -local github.com/unconditionalday internal/

gci:
	@find . -name "*.go" -type f -not -path '*/api/*' \
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

.PHONY: prepare-db

prepare-db:
	@sh ./scripts/prepare_db.sh

.PHONY: migrate

migrate:
	@goose --dir db/migration postgres "postgresql://${UNCONDITIONAL_API_DATABASE_USER}:${UNCONDITIONAL_API_DATABASE_PASSWORD}@localhost:5432/${UNCONDITIONAL_API_DATABASE_NAME}?sslmode=disable" up

# Helpers
check-variable-%: # detection of undefined variables.
	@[[ "${${*}}" ]] || (echo '*** Please define variable `${*}` ***' && exit 1)
