VERSION  ?= unknown
LDFLAGS  := -w -s
NAME     := go-files-watcher 
GIT_SHA  ?= $(shell git rev-parse --short HEAD)
GOLANGCI_VERSION = v1.36.0

GOLANGCI := $(shell which golangci-lint 2>/dev/null)
ifeq ($(GOLANGCI),)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin ${GOLANGCI_VERSION}
endif

build: LDFLAGS += -X 'main.Timestamp=$(shell date +%s)'
build: LDFLAGS += -X 'main.Version=${VERSION}'
build: LDFLAGS += -X 'main.GitSHA=${GIT_SHA}'
build: LDFLAGS += -X 'main.ServiceName=${NAME}'
build:
	$(info building binary to cmd/bin/$(NAME) with flags $(LDFLAGS))
	@go build -race -o cmd/bin/$(NAME) -ldflags "$(LDFLAGS)" ./cmd/go-files-watcher/main.go

deps:
	@go mod download
	@go mod tidy

lint:
	${GOLANGCI} -v run --out-format=line-number

test:
	go test --race -covermode=atomic -coverprofile=coverage.out ./...

cover:
	@LOG_LEVEL=debug TMP_COV=$(shell mktemp); \
	go test -failfast -coverpkg=./... -coverprofile=$$TMP_COV ./... && \
	go tool cover -html=$$TMP_COV && rm $$TMP_COV

run:
	go run ./cmd/go-files-watcher/main.go

all: deps lint test build

.PHONY: deps lint test cover build run
