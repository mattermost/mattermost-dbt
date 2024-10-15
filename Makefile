# Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
# See LICENSE.txt for license information.

################################################################################
##                             VERSION PARAMS                                 ##
################################################################################

## Tool Versions
GOLANG_VERSION := $(shell cat go.mod | grep "^go " | cut -d " " -f 2)
ARCH ?= amd64

################################################################################

## Variables
GO ?= $(shell command -v go 2> /dev/null)
PACKAGES=$(shell go list ./...)
APPNAME := mmdbt

BUILD_HASH = $(shell git rev-parse HEAD)
BUILD_TIME := $(shell date -u +%Y%m%d.%H%M%S)

LDFLAGS += -X "github.com/mattermost/mattermost-dbt/model.BuildHash=$(BUILD_HASH)"

################################################################################

# Binaries.
TOOLS_BIN_DIR := $(abspath bin)

GOLINT_VER := master
GOLINT_BIN := golint
GOLINT_GEN := $(TOOLS_BIN_DIR)/$(GOLINT_BIN)

GOIMPORTS_VER := master
GOIMPORTS_BIN := goimports
GOIMPORTS := $(TOOLS_BIN_DIR)/$(GOIMPORTS_BIN)

################################################################################

export GO111MODULE=on

all: check-style unittest

.PHONY: check-style
check-style: govet goformat
	@echo Checking for style guide compliance

.PHONY: vet
govet:
	@echo Running govet
	$(GO) vet ./...
	@echo Govet success

## Checks if files are formatted with go fmt.
.PHONY: goformat
goformat:
	@echo Checking if code is formatted
	@for package in $(PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			gofmt_output=$$(gofmt -d -s $$files 2>&1); \
			if [ "$$gofmt_output" ]; then \
				echo "$$gofmt_output"; \
				echo "gofmt failed"; \
				echo "To fix it, run:"; \
				echo "go fmt [FAILED_PACKAGE]"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "gofmt success"; \


## Checks if imports are formatted correctly.
.PHONY: goimports
goimports: $(GOIMPORTS)
	@echo Checking if imports are sorted
	@for package in $(PACKAGES); do \
		echo "Checking "$$package; \
		files=$$(go list -f '{{range .GoFiles}}{{$$.Dir}}/{{.}} {{end}}' $$package); \
		if [ "$$files" ]; then \
			goimports_output=$$($(GOIMPORTS) -d $$files 2>&1); \
			if [ "$$goimports_output" ]; then \
				echo "$$goimports_output"; \
				echo "goimports failed"; \
				echo "To fix it, run:"; \
				echo "goimports -w [FAILED_PACKAGE]"; \
				exit 1; \
			fi; \
		fi; \
	done
	@echo "goimports success"; \


.PHONY: unittest
unittest:
	$(GO) test ./... -v -covermode=count -coverprofile=coverage.out

# Build for distribution
.PHONY: build
build:
	@echo Building mmdbt for ARCH=$(ARCH)
	@if [ "$(ARCH)" = "amd64" ]; then \
		export GOARCH="amd64"; \
	elif [ "$(ARCH)" = "arm64" ]; then \
		export GOARCH="arm64"; \
	elif [ "$(ARCH)" = "arm" ]; then \
		export GOARCH="arm"; \
	else \
		echo "Unknown architecture $(ARCH)"; \
		exit 1; \
	fi; \
	env GOOS=linux $(GO) build -ldflags '$(LDFLAGS)' -o $(APPNAME) ./cmd/$(APPNAME)

# Build and install locally
.PHONY: install
install:
	@echo Installing mmdbt...
	$(GO) install -ldflags '$(LDFLAGS)' ./cmd/$(APPNAME)

# Run go mod verify.
verify-gomod:
	$(GO) mod download
	$(GO) mod verify
