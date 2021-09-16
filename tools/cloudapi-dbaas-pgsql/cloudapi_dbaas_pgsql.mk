## This Makefile contains operations
## for CloudApiDbaasPgsql resources:
## Tests, Mocks, Documentation

.PHONY: cloudapi_dbaas_pgsql_test_unit
cloudapi_dbaas_pgsql_test_unit:
	@echo "--- Run unit tests for CloudApiDbaasPgsql ---"
	@go test -cover ./commands/cloudapi-dbaas-pgsql/... ./services/cloudapi-dbaas-pgsql/...
	@echo "DONE"

.PHONY: cloudapi_dbaas_pgsql_test
cloudapi_dbaas_pgsql_test: cloudapi_dbaas_pgsql_test_unit

.PHONY: cloudapi_dbaas_pgsql_mocks_update
cloudapi_dbaas_pgsql_mocks_update:
	@echo "--- Update mocks for CloudApiDbaasPgsql ---"
	@tools/cloudapi-dbaas-pgsql/regenerate_mocks.sh
	@echo "DONE"
