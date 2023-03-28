.DEFAULT_GOAL := build

## Include Services Makefile Targets
include ./tools/cloudapi-v6/cloudapi_v6.mk
include ./tools/dbaas-postgres/dbaas_postgres.mk
include ./tools/dbaas-mongo/dbaas_mongo.mk
include ./tools/auth-v1/auth_v1.mk
include ./tools/certmanager/certmanager.mk
include ./tools/dataplatform/dataplatform.mk

include ./tools/container-registry/contregistry.mk

export CGO_ENABLED = 0
export GO111MODULE := on

GOFILES_NOVENDOR=$(shell find . -type f -name '*.go' | grep -v vendor)
GOOS?=$(shell go env GOOS)
GOARCH?=$(shell go env GOARCH)

OUT_D?=$(shell pwd)/builds
DOCS_OUT?=$(shell pwd)/docs/subcommands

.PHONY: test_unit
test_unit:
	@echo "--- Run unit tests ---"
	@go test -cover ./commands/ ./pkg/...
	@echo "DONE"

# run unit tests for all services
.PHONY: utest
utest: test_unit cloudapiv6_test auth_v1_test dbaas_postgres_test dbaas_mongo_test_unit certmanager_test_unit dataplatform_test_unit contreg_test_unit

# run integration tests for all services
.PHONY: itest
itest: dbaas_mongo_test_integration certmanager_test_integration dataplatform_test # contreg_test_integration # Temp Skip because 409 Conflict

# run all tests
.PHONY: test
test: utest itest

.PHONY: mocks_update
mocks_update: cloudapiv6_mocks_update auth_v1_mocks_update dbaas_postgres_mocks_update certmanager_mocks_update dbaas_mongo_mocks_update
	@echo "--- Update mocks ---"
	@tools/regenerate_mocks.sh
	@echo "DONE"

.PHONY: docs_update
docs_update: dbaas_postgres_docs_update certmanager_docs_update dbaas_mongo_docs_update dataplatform_docs_update
	@echo "--- Update documentation in ${DOCS_OUT} ---"
	@mkdir -p ${DOCS_OUT}
	@DOCS_OUT=${DOCS_OUT} tools/regenerate_doc.sh
	@echo "DONE"

.PHONY: gofmt_check
gofmt_check:
	@echo "--- Ensure code adheres to gofmt and list files whose formatting differs(vendor directory excluded) ---"
	@if [ "$(shell echo $$(gofmt -l ${GOFILES_NOVENDOR}))" != "" ]; then (echo "Format files: $(shell echo $$(gofmt -l ${GOFILES_NOVENDOR})) Hint: use \`make gofmt_update\`"; exit 1); fi
	@echo "DONE"

.PHONY: gofmt_update
gofmt_update:
	@echo "--- Ensure code adheres to gofmt and change files accordingly(vendor directory excluded) ---"
	@gofmt -w ${GOFILES_NOVENDOR}
	@echo "DONE"

.PHONY: goimports_update
goimports_update:
	@echo "--- Ensure code adheres to goimports and change files accordingly(vendor directory excluded) ---"
	@goimports -w ${GOFILES_NOVENDOR}
	@echo "DONE"

.PHONY: vendor_status
vendor_status:
	@govendor status

.PHONY: vendor_update
vendor_update:
	@echo "--- Update vendor dependencies ---"
	@go mod vendor
	@go mod tidy
	@echo "DONE"

.PHONY: build
build:
	@echo "--- Building ionosctl via go build ---"
	@OUT_D=${OUT_D} GOOS=$(GOOS) GOARCH=$(GOARCH) tools/build.sh
	@echo "built ${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
	@echo "DONE"

.PHONY: install
install:
	@echo "--- Install ionosctl via go install ---"
	@GOOS=$(GOOS) GOARCH=$(GOARCH) tools/install.sh
	@echo "DONE"

.PHONY: clean
clean:
	@echo "--- Remove built / installed artifacts ---"
	@go clean -i
	@rm -rf builds
	@echo "DONE"
