DOCS_OUT_CONTAINER_REGISTRY?=$(shell pwd)/docs/subcommands/container-registry/

.PHONY: contreg_docs_update
contreg_docs_update:
	@echo "--- Generate Markdown documentation for Container Registry in ${DOCS_OUT_CONTAINER_REGISTRY} ---"
	@mkdir -p ${DOCS_OUT_CONTAINER_REGISTRY}
	@DOCS_OUT_CONTAINER_REGISTRY=${DOCS_OUT_CONTAINER_REGISTRY} go run tools/container-registry/doc.go
	@echo "DONE"