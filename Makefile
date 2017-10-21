all: build

.PHONY: build

ifndef ($(GOPATH))
	GOPATH = $(HOME)/go
endif

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
	-X github.com/micromdm/squirrel/vendor/github.com/micromdm/go4/version.appName=${APP_NAME} \
	-X github.com/micromdm/squirrel/vendor/github.com/micromdm/go4/version.version=${VERSION} \
	-X github.com/micromdm/squirrel/vendor/github.com/micromdm/go4/version.branch=${BRANCH} \
	-X github.com/micromdm/squirrel/vendor/github.com/micromdm/go4/version.buildUser=${USER} \
	-X github.com/micromdm/squirrel/vendor/github.com/micromdm/go4/version.buildDate=${NOW} \
	-X github.com/micromdm/squirrel/vendor/github.com/micromdm/go4/version.revision=${REVISION} \
	-X github.com/micromdm/squirrel/vendor/github.com/micromdm/go4/version.goVersion=${GOVERSION}"

WORKSPACE = ${GOPATH}/src/github.com/micromdm/squirrel
check-deps:
ifneq ($(shell test -e ${WORKSPACE}/Gopkg.lock && echo -n yes), yes)
	@echo "folder is clonded in the wrong place, copying to a Go Workspace"
	@echo "See: https://golang.org/doc/code.html#Workspaces"
	@git clone git@github.com:micromdm/squirrel ${WORKSPACE}
	@echo "cd to ${WORKSPACE} and run make deps again."
	@exit 1
endif
ifneq ($(shell pwd), $(WORKSPACE))
	@echo "cd to ${WORKSPACE} and run make deps again."
	@exit 1
endif

deps: check-deps
	go get -u github.com/golang/dep/...
	dep ensure -vendor-only

test:
	go test -cover -race -v $(shell go list ./... | grep -v /vendor/)

build: squirrel

clean:
	rm -rf build/
	rm -f *.zip

.pre-build:
	mkdir -p build/darwin
	mkdir -p build/linux

INSTALL_STEPS := \
	install-squirrel 

install-local: $(INSTALL_STEPS)

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

