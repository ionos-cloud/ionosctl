.DEFAULT_GOAL := build

export CGO_ENABLED = 0
export GO111MODULE := on
export GOFLAGS := -cover

GOFILES_NOVENDOR=$(shell find . -type f -name '*.go' | grep -v vendor)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

OUT_D?=$(shell pwd)/builds
DOCS_OUT?=$(shell pwd)/docs/subcommands
TEST_DIRS := ./commands/... ./pkg/... ./services/... ./internal/...

.PHONY: utest test_unit
utest test_unit:
	@echo "--- Run unit tests ---"
	@go test -v $(TEST_DIRS) && echo "DONE"

.PHONY: test itest test_integration
itest test test_integration:
	@echo "--- Run integration and unit tests ---"
	@go test -v -tags=integration $(TEST_DIRS) && echo "DONE"

# Note about test file tagging:
# `//go:build integration` was introduced in Go 1.17
# `// +build integration` is still maintained for compatibility reasons
# `go fmt` still maintains these lines, if one is removed. If it stops this behaviour, then we can remove them

.PHONY: mocks
mocks:
	@echo "--- Update mocks ---"
	@tools/mocks.sh && echo "DONE"

.PHONY: docs generate-docs
docs generate-docs:
	@echo "--- Purging docs ---"
	rm -rf docs/subcommands
	rm docs/summary.md
	@echo "--- Regenerating docs ---"
	@go run tools/doc.go


.PHONY: gofmt_check
gofmt_check:
	@echo "--- Ensure code adheres to gofmt and list files whose formatting differs(vendor directory excluded) ---"
	@if [ "$(shell echo $$(gofmt -l ${GOFILES_NOVENDOR}))" != "" ]; then (echo "Format files: $(shell echo $$(gofmt -l ${GOFILES_NOVENDOR})) Hint: use \`make gofmt_update\`"; exit 1); fi
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

.PHONY: build
build: vendor
	@echo "--- Building ionosctl via go build ---"
	@OUT_D=${OUT_D} GOOS=$(GOOS) GOARCH=$(GOARCH) tools/build.sh build
	@echo "built ${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
	@echo "DONE"

.PHONY: install
install: vendor
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
	@echo " - test:\tRun integration and unit tests (CI Target)"
	@echo " - mocks:\tUpdate mocks. WARNING: Do not interrupt early!"
	@echo " - gofmt:\tFormat code to adhere to gofmt [gofmt_check for checking only]"
	@echo " - vendor:\tUpdate vendor dependencies. [vendor_check for checking only]"
	@echo " - goimports:\tFormat, sort imports"
	@echo " - docs:\tRegenerate docs"
	@echo " - build/install/clean"
