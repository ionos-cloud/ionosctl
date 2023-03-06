DOCS_OUT_CERT_MANAGER?=$(shell pwd)/docs/subcommands/certificate-manager/

.PHONY: certmanager_mocks_update
certmanager_mocks_update:
	@echo "--- Update mocks for CloudApi CertManager ---"
	@tools/certmanager/regenerate_mocks.sh
	@echo "DONE"

.PHONY: certmanager_test_integration
certmanager_test_integration:
	@echo "--- Run tests for Certificate Manger ---"
	@go test -cover ./commands/certmanager/...
	@echo "DONE"

.PHONY: certmanager_test
certmanager_test: certmanager_test_integration

.PHONY: certmanager_docs_update
certmanager_docs_update:
	@echo "--- Generate Markdown documentation for CertManager in ${DOCS_OUT_CERT_MANAGER} ---"
	@mkdir -p ${DOCS_OUT_CERT_MANAGER}
	@DOCS_OUT_CERT_MANAGER=${DOCS_OUT_CERT_MANAGER} go run tools/certmanager/doc.go
	@echo "DONE"