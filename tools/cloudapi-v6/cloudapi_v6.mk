## This Makefile contains operations
## for CloudApiV6 resources:
## Tests, Mocks, Documentation

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
