## This Makefile contains operations
## for CloudApiDBaaSPgsql resources:
## Tests, Mocks

.PHONY: dbaas_pg_test_unit
dbaas_pg_test_unit:
	@echo "--- Run unit tests for CloudApi DBaaS Postgres ---"
	@go test -cover ./commands/pg/... ./services/dbaas-pg/...
	@echo "DONE"

.PHONY: dbaas_pg_test
dbaas_pg_test: dbaas_pg_test_unit

.PHONY: dbaas_pg_mocks_update
dbaas_pg_mocks_update:
	@echo "--- Update mocks for CloudApi DBaaS Postgres ---"
	@tools/dbaas-pg/regenerate_mocks.sh
	@echo "DONE"
