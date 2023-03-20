DOCS_OUT_CONTAINER_REGISTRY?=$(shell pwd)/docs/subcommands/container-registry/

.PHONY: contreg_docs_update
contreg_docs_update:
	@echo "--- Generate Markdown documentation for Container Registry in ${DOCS_OUT_CONTAINER_REGISTRY} ---"
	@mkdir -p ${DOCS_OUT_CONTAINER_REGISTRY}
	@DOCS_OUT_CONTAINER_REGISTRY=${DOCS_OUT_CONTAINER_REGISTRY} go run tools/container-registry/doc.go
	@echo "DONE"

.PHONY: contreg_test_unit
contreg_test_unit:
	@echo "--- Run Unit tests for Container Registry ---"
	@go test -tags=unit ./commands/container-registry/...
	@echo "DONE"

.PHONY: contreg_test_integration
contreg_test_integration:
	@echo "--- Run Integration tests for Container Registry ---"
	@go test -tags=integration ./commands/container-registry/...
	@echo "DONE"

.PHONY: contreg_test
contreg_test: contreg_test_unit contreg_test_integration