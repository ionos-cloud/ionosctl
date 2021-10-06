## This Makefile contains operations
## for CloudApiDbaasPgsql resources:
## Tests, Mocks, Documentation

DOCS_OUT_DBAAS_PGSQL?=$(shell pwd)/docs/cloudapi-dbaas-pgsql/

.PHONY: cloudapi_dbaaspgsql_test_unit
cloudapi_dbaaspgsql_test_unit:
	@echo "--- Run unit tests for CloudApi DBaaS Postgres ---"
	@go test -cover ./commands/cloudapi-dbaas-pgsql/... ./services/cloudapi-dbaas-pgsql/...
	@echo "DONE"

.PHONY: cloudapi_dbaaspgsql_test
cloudapi_dbaaspgsql_test: cloudapi_dbaaspgsql_test_unit

.PHONY: cloudapi_dbaaspgsql_mocks_update
cloudapi_dbaaspgsql_mocks_update:
	@echo "--- Update mocks for CloudApi DBaaS Postgres ---"
	@tools/cloudapi-dbaas-pgsql/regenerate_mocks.sh
	@echo "DONE"

.PHONY: cloudapi_dbaaspgsql_docs_update
cloudapi_dbaaspgsql_docs_update:
	@echo "--- Generate Markdown documentation for CloudApi DBaaS Postgres in ${DOCS_OUT_DBAAS_PGSQL} ---"
	@mkdir -p ${DOCS_OUT_DBAAS_PGSQL}
	@DOCS_OUT_DBAAS_PGSQL=${DOCS_OUT_DBAAS_PGSQL} go run tools/cloudapi-dbaas-pgsql/doc.go
	@echo "DONE"
