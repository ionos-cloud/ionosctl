.DEFAULT_GOAL := build

export CGO_ENABLED = 0
export GO111MODULE := on

GOFILES_NOVENDOR=$(shell find . -type f -name '*.go' | grep -v vendor)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

OUT_D?=$(shell pwd)/builds
DOCS_OUT?=$(shell pwd)/docs/subcommands

# Want to test verbosely? (i.e. see what test is failing?) Run like:
# make test TEST_FLAGS="-v [optionally other flags]"

TEST_DIRS := $(shell go list ./... | \
				grep -v /commands/container-registry | \
				grep -v /commands/dbaas/mongo | \
				grep -v /commands/logging-service) # All commands except ...
TEST_FLAGS := "-cover"
.PHONY: utest
utest:
	@echo "--- Run unit tests ---"
	@go test $(TEST_FLAGS) $(TEST_DIRS)

# Note about test file tagging:
# `//go:build integration` was introduced in Go 1.17
# `// +build integration` is still maintained for compatibility reasons
# `go fmt` still maintains these lines, if one is removed. If it stops this behaviour, then we can remove them

# run the go-based e2e tests
.PHONY: itest
itest:
	@echo "--- Run unit tests and go-based integration tests ---"
	@go test $(TEST_FLAGS) -tags=integration $(TEST_DIRS)

.PHONY: test
test:
	@echo "--- Run bats-core tests ---"
	@test/run.sh # bats-core tests and other
	@$(MAKE) itest # go-based tests (unit and integration)

.PHONY: mocks
mocks:
	@echo "--- Update mocks ---"
	@echo "--- Ensure gomock is up to date. Run: 'go install github.com/golang/mock/mockgen@v1.6.0' ---"
	@tools/mocks.sh && echo "DONE"

.PHONY: docs generate-docs
docs generate-docs:
	@echo "--- Purging docs ---"
	rm -rf docs/subcommands
	touch docs/summary.md
	rm docs/summary.md
	@echo "--- Regenerating docs ---"
	@go run tools/doc.go


.PHONY: gofmt_check
gofmt_check:
	@echo "--- Ensure code adheres to gofmt and list files whose formatting differs(vendor directory excluded) ---"
	@if [ "$(shell echo $$(gofmt -l ${GOFILES_NOVENDOR}))" != "" ]; then (echo "Format files: $(shell echo $$(gofmt -l ${GOFILES_NOVENDOR})) Hint: use \`make gofmt\`"; exit 1); fi
	@echo "DONE"

.PHONY: gofmt
gofmt:
	@echo "--- Ensure code adheres to gofmt and change files accordingly(vendor directory excluded) ---"
	@gofmt -w ${GOFILES_NOVENDOR} && echo "DONE"

.PHONY: goimports
goimports:
	@echo "--- Ensure code adheres to goimports and change files accordingly(vendor directory excluded) ---"
	@goimports -w ${GOFILES_NOVENDOR} && echo "DONE"

.PHONY: vendor_check
vendor_check:
	@govendor status

.PHONY: vendor
vendor:
	@echo "--- Update vendor dependencies ---"
	@go mod vendor
	@go mod tidy
	@echo "DONE"

BINARY_NAME ?= ionosctl # To install with a custom name, e.g. `io`, do `make install BINARY_NAME=io`

.PHONY: build
build: vendor
	@echo "--- Building ionosctl via go build ---"
	@OUT_D=${OUT_D} GOOS=$(GOOS) GOARCH=$(GOARCH) tools/build.sh build
	@echo "built ${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
	@echo "DONE"

.PHONY: install
install:
	@echo "--- Install ionosctl via go install ---"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) tools/build.sh install
	@echo "DONE"

.PHONY: clean
clean:
	@echo "--- Remove built / installed artifacts ---"
	@go clean -i
	@rm -rf builds
	@echo "DONE"

.PHONY: help
help:
	@echo "TARGETS: "
	@echo " - utest:\tRun unit tests"
	@echo " - itest:\tRun all go-based tests (unit, integration)"
	@echo " - test:\tRuns bats-core tests, then runs all go-based tests (CI Target)"
	@echo " - mocks:\tUpdate mocks. WARNING: Do not interrupt early!"
	@echo " - gofmt:\tFormat code to adhere to gofmt [gofmt_check for checking only]"
	@echo " - vendor:\tUpdate vendor dependencies. [vendor_check for checking only]"
	@echo " - goimports:\tFormat, sort imports"
	@echo " - docs:\tRegenerate docs"
	@echo " - build/install/clean"
