## This Makefile contains operations
## for CloudApiDBaaSPgsql resources:
## Tests, Mocks, Documentation

DOCS_OUT_DBAAS_PGSQL?=$(shell pwd)/docs/dbaas-pg/

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

.PHONY: dbaas_pg_docs_update
dbaas_pg_docs_update:
	@echo "--- Generate Markdown documentation for CloudApi DBaaS Postgres in ${DOCS_OUT_DBAAS_PGSQL} ---"
	@mkdir -p ${DOCS_OUT_DBAAS_PGSQL}
	@DOCS_OUT_DBAAS_PGSQL=${DOCS_OUT_DBAAS_PGSQL} go run tools/dbaas-pg/doc.go
	@echo "DONE"
