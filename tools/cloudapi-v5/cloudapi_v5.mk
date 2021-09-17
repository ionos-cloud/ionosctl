## This Makefile contains operations
## for CloudApiV5 resources:
## Tests, Mocks, Documentation

.PHONY: cloudapiv5_test_unit
cloudapiv5_test_unit:
	@echo "--- Run unit tests for CloudApiV5 ---"
	@go test -cover ./commands/cloudapi-v5/... ./services/cloudapi-v5/...
	@echo "DONE"

.PHONY: cloudapiv5_test
cloudapiv5_test: cloudapiv5_test_unit

.PHONY: cloudapiv5_mocks_update
cloudapiv5_mocks_update:
	@echo "--- Update mocks for CloudApiV5 ---"
	@tools/cloudapi-v5/regenerate_mocks.sh
	@echo "DONE"
