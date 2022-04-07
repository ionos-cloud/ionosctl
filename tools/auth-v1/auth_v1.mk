## This Makefile contains operations
## for Auth resources:
## Tests, Mocks, Documentation

.PHONY: auth_v1_test_unit
auth_v1_test_unit:
	@echo "--- Run unit tests for Auth V1 ---"
	@go test -cover ./commands/auth-v1/... ./services/auth-v1/...
	@echo "DONE"

.PHONY: auth_v1_test
auth_v1_test: auth_v1_test_unit

.PHONY: auth_v1_mocks_update
auth_v1_mocks_update:
	@echo "--- Update mocks for Auth V1 ---"
	@tools/auth-v1/regenerate_mocks.sh
	@echo "DONE"
