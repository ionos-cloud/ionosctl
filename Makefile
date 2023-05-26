.DEFAULT_GOAL := build

export CGO_ENABLED = 0
export GO111MODULE := on
export GOFLAGS := -cover

GOFILES_NOVENDOR=$(shell find . -type f -name '*.go' | grep -v vendor)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

OUT_D?=$(shell pwd)/builds
DOCS_OUT?=$(shell pwd)/docs/subcommands
TEST_DIRS := ./commands/... ./pkg/...

.PHONY: utest test_unit
utest test_unit:
	@echo "--- Run unit tests ---"
	@go test $(TEST_DIRS) && echo "DONE"

.PHONY: test itest test_integration
itest test test_integration:
	@echo "--- Run integration and unit tests ---"
	@go test -tags=integration $(TEST_DIRS) && echo "DONE"

# Note about test file tagging:
# `//go:build integration` was introduced in Go 1.17
# `// +build integration` is still maintained for compatibility reasons
# `go fmt` still maintains these lines, if one is removed. If it stops this behaviour, then we can remove them

.PHONY: mocks_update
mocks_update: cloudapiv6_mocks_update auth_v1_mocks_update dbaas_postgres_mocks_update certmanager_mocks_update dbaas_mongo_mocks_update
	@echo "--- Update mocks ---"
	@tools/regenerate_mocks.sh && echo "DONE"

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

.PHONY: gofmt_update
gofmt_update:
	@echo "--- Ensure code adheres to gofmt and change files accordingly(vendor directory excluded) ---"
	@gofmt -w ${GOFILES_NOVENDOR} && echo "DONE"

.PHONY: goimports_update
goimports_update:
	@echo "--- Ensure code adheres to goimports and change files accordingly(vendor directory excluded) ---"
	@goimports -w ${GOFILES_NOVENDOR} && echo "DONE"

.PHONY: vendor_status
vendor_status:
	@govendor status

.PHONY: vendor_update
vendor_update:
	@echo "--- Update vendor dependencies ---"
	@go mod vendor
	@go mod tidy
	@echo "DONE"

.PHONY: build
build: vendor_update
	@echo "--- Building ionosctl via go build ---"
	@OUT_D=${OUT_D} GOOS=$(GOOS) GOARCH=$(GOARCH) tools/build.sh build
	@echo "built ${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
	@echo "DONE"

.PHONY: install
install: vendor_update
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
	@echo "The following are some of the valid targets for this Makefile:"
	@echo "... utest: Run unit tests"
	@echo "... test: Run integration and unit tests (CI Target)"
	@echo "... mocks_update: Update mocks (Used in some legacy tests)"
	@echo "... docs: Regenerate docs"
	@echo "... gofmt_check: Check code adheres to gofmt"
	@echo "... gofmt_update: Format code to adhere to gofmt"
	@echo "... goimports_update: Format code to adhere to goimports"
	@echo "... vendor_status: Check vendor status"
	@echo "... vendor_update: Update vendor dependencies"
	@echo "... build: Build ionosctl"
	@echo "... install: Install ionosctl"
	@echo "... clean: Remove built / installed artifacts"
