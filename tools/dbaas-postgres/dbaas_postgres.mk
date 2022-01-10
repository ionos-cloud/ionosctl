## This Makefile contains operations
## for CloudApiDBaaSPgsql resources:
## Tests, Mocks

.PHONY: dbaas_postgres_test_unit
dbaas_postgres_test_unit:
	@echo "--- Run unit tests for CloudApi DBaaS Postgres ---"
	@go test -cover ./commands/dbaas/postgres/... ./services/dbaas-postgres/...
	@echo "DONE"

.PHONY: dbaas_postgres_test
dbaas_postgres_test: dbaas_postgres_test_unit

.PHONY: dbaas_postgres_mocks_update
dbaas_postgres_mocks_update:
	@echo "--- Update mocks for CloudApi DBaaS Postgres ---"
	@tools/dbaas-postgres/regenerate_mocks.sh
	@echo "DONE"
