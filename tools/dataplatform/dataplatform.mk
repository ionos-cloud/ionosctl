## This Makefile contains operations
## for Data Platform resources:
## Tests, Mocks, Documentation
DOCS_OUT_DATAPLATFORM?=$(shell pwd)/docs/subcommands/dataplatform/

.PHONY: dataplatform_test_unit
dataplatform_test_unit:
	@echo "--- Run unit tests for Data Platform---"
	@go test -cover ./commands/dataplatform/... ./services/dataplatform/...
	@echo "DONE"

.PHONY: dataplatform_test
dataplatform_test: dataplatform_test_unit

.PHONY: dataplatform_mocks_update
dataplatform_mocks_update:
	@echo "--- Update mocks for Data Platform ---"
	@tools/dataplatform/regenerate_mocks.sh
	@echo "DONE"

.PHONY: dataplatform_docs_update
dataplatform_docs_update:
	@echo "--- Generate Markdown documentation for Data Platform in ${DOCS_OUT_DATAPLATFORM} ---"
	@mkdir -p ${DOCS_OUT_DATAPLATFORM}
	@DOCS_OUT_DATAPLATFORM=${DOCS_OUT_DATAPLATFORM} go run tools/dataplatform/doc.go
	@echo "DONE"
