## This Makefile contains operations
## for CloudApiDBaaSPgsql resources:
## Tests, Mocks, Documentation
DOCS_OUT_DBAAS_MONGO?=$(shell pwd)/docs/subcommands/database-as-a-service/mongo/

.PHONY: dbaas_mongo_test_unit
dbaas_mongo_test_unit:
	@echo "--- Run Unit tests for CloudApi DBaaS mongo ---"
	@go test -cover ./commands/dbaas/mongo/... ./services/dbaas-mongo/...
	@echo "DONE"

.PHONY: dbaas_mongo_test_integration
dbaas_mongo_test_integration:
	@echo "--- Run Integration tests for CloudApi DBaaS mongo ---"
	@go test -cover -tags=integration ./commands/dbaas/mongo/... ./services/dbaas-mongo/...
	@echo "DONE"

.PHONY: dbaas_mongo_test
dbaas_mongo_test: dbaas_mongo_test_unit dbaas_mongo_test_integration

.PHONY: dbaas_mongo_mocks_update
dbaas_mongo_mocks_update:
	@echo "--- Update mocks for CloudApi DBaaS mongo ---"
	@tools/dbaas-mongo/regenerate_mocks.sh
	@echo "DONE"

