
export CGO_ENABLED = 0
export GO111MODULE := on

GOFILES= $(shell find . -type f -name '*.go')

.PHONY: test_unit
test_unit:
	@echo "==> run unit tests"
	go test -cover ./commands/... ./pkg/utils/...

# Generate Markdown documentation for IonosCTL commands
.PHONY: docs
docs:
	@echo "==> Generate Markdown documentation in ${DOCS_OUT}"
	@DOCS_OUT=${DOCS_OUT} go run tools/doc.go
	@echo "DONE"

.PHONY: gofmt_check
gofmt_check:
	@echo "==> ensure code adheres to gofmt"
	@gofmt -w ${GOFILES}
	@echo "DONE"
