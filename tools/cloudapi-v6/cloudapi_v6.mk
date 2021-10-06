## This Makefile contains operations
## for CloudApiV6 resources:
## Tests, Mocks, Documentation

DOCS_OUT_V6?=$(shell pwd)/docs/cloudapiv6/

.PHONY: cloudapiv6_test_unit
cloudapiv6_test_unit:
	@echo "--- Run unit tests for CloudApiV6 ---"
	@go test -cover ./commands/cloudapi-v6/... ./services/cloudapi-v6/...
	@echo "DONE"

.PHONY: cloudapiv6_test
cloudapiv6_test: cloudapiv6_test_unit

.PHONY: cloudapiv6_mocks_update
cloudapiv6_mocks_update:
	@echo "--- Update mocks for CloudApiV6 ---"
	@tools/cloudapi-v6/regenerate_mocks.sh
	@echo "DONE"

.PHONY: cloudapiv6_docs_update
cloudapiv6_docs_update:
	@echo "--- Generate Markdown documentation for CloudApiV6 in ${DOCS_OUT_V6} ---"
	@mkdir -p ${DOCS_OUT_V6}
	@DOCS_OUT_V6=${DOCS_OUT_V6} go run tools/cloudapi-v6/doc.go
	@echo "DONE"
