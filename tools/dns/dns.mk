DOCS_OUT_CONTAINER_REGISTRY?=$(shell pwd)/docs/subcommands/dns/

.PHONY: dns_test_unit
dns_test_unit:
	@echo "--- Run Unit tests for DNS ---"
	@go test -tags=unit ./commands/dns/...
	@echo "DONE"

.PHONY: dns_test_integration
dns_test_integration:
	@echo "--- Run Integration tests for DNS ---"
	@go test -tags=integration ./commands/dns/...
	@echo "DONE"

.PHONY: dns_test
dns_test: dns_test_unit dns_test_integration
