DOCS_OUT_DATAPLATFORM?=$(shell pwd)/docs/subcommands/dataplatform/

.PHONY: dataplatform_test_unit
dataplatform_test_unit:
	@echo "--- Run Unit tests for Dataplatform ---"
	@go test -tags=unit -cover ./commands/dataplatform/...
	@echo "DONE"

.PHONY: dataplatform_test_integration
dataplatform_test_integration:
	@echo "--- Run Integration tests for Dataplatform ---"
	@go test -tags=integration -cover ./commands/dataplatform/...
	@echo "DONE"

.PHONY: dataplatform_test
dataplatform_test: dataplatform_test_unit dataplatform_test_integration

.PHONY: dataplatform_mocks_update
dataplatform_mocks_update:
	@echo "--- Update mocks for Dataplatform ---"
	@tools/dataplatform/regenerate_mocks.sh
	@echo "DONE"
