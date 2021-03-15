export CGO_ENABLED = 0
export GO111MODULE := on

GOFILES_NOVENDOR=$(shell find . -type f -name '*.go' | grep -v vendor)
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

OUT_D=$${OUT_D:-$(shell pwd)/builds}
DOCS_OUT=$${DOCS_OUT:-$(shell pwd)/docs/commands/}

.PHONY: test_unit
test_unit:
	@echo "--- Run unit tests ---"
	@go test -cover ./commands/... ./pkg/utils/...
	@echo "DONE"

.PHONY: test
test: test_unit

.PHONY: docs_update
docs_update:
	@echo "--- Generate Markdown documentation in ${DOCS_OUT} ---"
	@DOCS_OUT=${DOCS_OUT} go run tools/doc.go
	@echo "DONE"

.PHONY: gofmt_check
gofmt_check:
	@echo "--- Ensure code adheres to gofmt and list files whose formatting differs from gofmt's(vendor directory excluded) ---"
	@if [ "$(shell echo $$(gofmt -l ${GOFILES_NOVENDOR}))" != "" ]; then (echo "Format files: $(shell echo $$(gofmt -l ${GOFILES_NOVENDOR})) Hint: use \`make gofmt_update\`"; exit 1); fi
	@echo "DONE"

.PHONY: gofmt_update
gofmt_update:
	@echo "--- Ensure code adheres to gofmt and change files accordingly(vendor directory excluded) ---"
	@gofmt -w ${GOFILES_NOVENDOR}
	@echo "DONE"

.PHONY: vendor_status
vendor_status:
	@govendor status

.PHONY: vendor_update
vendor_update:
	@echo "--- Update vendor dependencies ---"
	@go mod vendor
	@go mod tidy
	@echo "DONE"

.PHONY: mocks_update
mocks_update:
	@echo "--- Update mocks ---"
	@tools/regenerate_mocks.sh
	@echo "DONE"

.PHONY: build
build:
	@echo "--- Building ionosctl via go build ---"
	@OUT_D=${OUT_D} GOOS=$(GOOS) GOARCH=$(GOARCH) tools/build.sh
	@echo "built ${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
	@echo "DONE"

.PHONY: install
install:
	@echo "--- Install ionosctl via go install ---"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) tools/install.sh
	@echo "DONE"

.PHONY: clean
clean:
	@echo "--- Remove built / installed artifacts ---"
	@go clean -i
	@rm -rf builds
	@echo "DONE"
