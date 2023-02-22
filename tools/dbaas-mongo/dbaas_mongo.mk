## This Makefile contains operations
## for CloudApiDBaaSPgsql resources:
## Tests, Mocks, Documentation
DOCS_OUT_DBAAS_MONGO?=$(shell pwd)/docs/subcommands/database-as-a-service/mongo/

.PHONY: dbaas_mongo_test_unit
dbaas_mongo_test_unit:
	@echo "--- Run unit tests for CloudApi DBaaS mongo ---"
	@go test -cover ./commands/dbaas/mongo/... ./services/dbaas-mongo/...
	@echo "DONE"

.PHONY: dbaas_mongo_test
dbaas_mongo_test: dbaas_mongo_test_unit

.PHONY: dbaas_mongo_mocks_update
dbaas_mongo_mocks_update:
	@echo "--- Update mocks for CloudApi DBaaS mongo ---"
	@tools/dbaas-mongo/regenerate_mocks.sh
	@echo "DONE"

.PHONY: dbaas_mongo_docs_update
dbaas_mongo_docs_update:
	@echo "--- Generate Markdown documentation for DBaaS Mongo in ${DOCS_OUT_DBAAS_MONGO} ---"
	@mkdir -p ${DOCS_OUT_DBAAS_MONGO}
	@DOCS_OUT_DBAAS_MONGO=${DOCS_OUT_DBAAS_MONGO} go run tools/dbaas-mongo/doc.go
	@echo "DONE"
