all: build

.PHONY: build

ifeq ($(GOPATH),)
	PATH := $(HOME)/go/bin:$(PATH)
else
	PATH := $(GOPATH)/bin:$(PATH)
endif

export GO111MODULE=on

PATH := $(GOPATH)/bin:$(PATH)
VERSION = $(shell git describe --tags --always --dirty)
BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
REVISION = $(shell git rev-parse HEAD)
REVSHORT = $(shell git rev-parse --short HEAD)
USER = $(shell whoami)
GOVERSION = $(shell go version | awk '{print $$3}')
NOW	= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
SHELL = /bin/bash

ifneq ($(OS), Windows_NT)
	CURRENT_PLATFORM = linux
	ifeq ($(shell uname), Darwin)
		SHELL := /bin/bash
		CURRENT_PLATFORM = darwin
	endif
else
	CURRENT_PLATFORM = windows
endif

BUILD_VERSION = "\
	-X github.com/micromdm/go4/version.appName=${APP_NAME} \
	-X github.com/micromdm/go4/version.version=${VERSION} \
	-X github.com/micromdm/go4/version.branch=${BRANCH} \
	-X github.com/micromdm/go4/version.buildUser=${USER} \
	-X github.com/micromdm/go4/version.buildDate=${NOW} \
	-X github.com/micromdm/go4/version.revision=${REVISION} \
	-X github.com/micromdm/go4/version.goVersion=${GOVERSION}"

gomodcheck: 
	@go help mod > /dev/null || (@echo micromdm requires Go version 1.11 or higher && exit 1)

deps: gomodcheck
	@go mod download

test:
	go test -cover -race ./...

build: squirrel

clean:
	rm -rf build/
	rm -f *.zip

.pre-build:
	mkdir -p build/darwin
	mkdir -p build/linux

install-local: \
	install-squirrel 

.pre-squirrel:
	$(eval APP_NAME = squirrel)

squirrel: .pre-build .pre-squirrel
	go build -i -o build/$(CURRENT_PLATFORM)/squirrel -ldflags ${BUILD_VERSION} ./cmd/squirrel

install-squirrel: .pre-squirrel
	go install -ldflags ${BUILD_VERSION} ./cmd/squirrel

xp-squirrel: .pre-build .pre-squirrel
	GOOS=darwin go build -i -o build/darwin/squirrel -ldflags ${BUILD_VERSION} ./cmd/squirrel
	GOOS=linux CGO_ENABLED=0 go build -i -o build/linux/squirrel  -ldflags ${BUILD_VERSION} ./cmd/squirrel

release-zip: xp-squirrel
	zip -r squirrel_${VERSION}.zip build/

