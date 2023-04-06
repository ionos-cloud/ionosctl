DOCS_OUT_CONTAINER_REGISTRY?=$(shell pwd)/docs/subcommands/container-registry/

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
