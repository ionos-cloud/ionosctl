#!/usr/bin/env bash

## Temporarily, the existing documentation for the CLI is updated using the current script.

set -euo pipefail
DOCS_OUT=${DOCS_OUT:-$(shell pwd)/docs/subcommands}

## Generate all documentation in one place
DOCS_OUT_TMP=${DOCS_OUT_TMP:-${DOCS_OUT}/tmp}
echo "Generate Markdown documentation in ${DOCS_OUT_TMP}"
mkdir -p ${DOCS_OUT_TMP}
DOCS_OUT=${DOCS_OUT_TMP} go run tools/doc.go

## Separate documentation based on service

## cli-setup
DOCS_OUT_SETUP=${DOCS_OUT_SETUP:-${DOCS_OUT}/cli-setup/}
echo "Move CLI Setup documentation in ${DOCS_OUT_SETUP}"
mkdir -p ${DOCS_OUT_SETUP}
# mv -f docs/subcommands/completion-* ${DOCS_OUT_SETUP}
mv -f docs/subcommands/tmp/version.md ${DOCS_OUT_SETUP}
cp docs/subcommands/tmp/login.md ${DOCS_OUT_SETUP}

## authentication
DOCS_OUT_AUTH=${DOCS_OUT_AUTH:-${DOCS_OUT}/authentication/}
echo "Move Authentication documentation in ${DOCS_OUT_AUTH}"
mkdir -p ${DOCS_OUT_AUTH}
mv -f docs/subcommands/tmp/token-* ${DOCS_OUT_AUTH}
cp docs/subcommands/tmp/contract-* ${DOCS_OUT_AUTH}
mv -f docs/subcommands/tmp/login.md ${DOCS_OUT_AUTH}

## managed-kubernetes
DOCS_OUT_K8S=${DOCS_OUT_K8S:-${DOCS_OUT}/managed-kubernetes/}
echo "Move Kubernetes documentation in ${DOCS_OUT_K8S}"
mkdir -p ${DOCS_OUT_K8S}
mv -f docs/subcommands/tmp/k8s-* ${DOCS_OUT_K8S}

## networkloadbalancer
DOCS_OUT_NLB=${DOCS_OUT_NLB:-${DOCS_OUT}/networkloadbalancer/}
echo "Move Network Load Balancer documentation in ${DOCS_OUT_NLB}"
mkdir -p ${DOCS_OUT_NLB}
mv -f docs/subcommands/tmp/networkloadbalancer-* ${DOCS_OUT_NLB}



## natgateway
DOCS_OUT_NAT_GATEWAY=${DOCS_OUT_NAT_GATEWAY:-${DOCS_OUT}/natgateway/}
echo "Move NAT Gateway documentation in ${DOCS_OUT_NAT_GATEWAY}"
mkdir -p ${DOCS_OUT_NAT_GATEWAY}
mv -f docs/subcommands/tmp/natgateway-* ${DOCS_OUT_NAT_GATEWAY}

## managed-backup
DOCS_OUT_BACKUP=${DOCS_OUT_BACKUP:-${DOCS_OUT}/managed-backup/}
echo "Move Backup documentation in ${DOCS_OUT_BACKUP}"
mkdir -p ${DOCS_OUT_BACKUP}
mv -f docs/subcommands/tmp/backupunit-* ${DOCS_OUT_BACKUP}

## user-management
DOCS_OUT_USER=${DOCS_OUT_USER:-${DOCS_OUT}/user-management/}
echo "Move User documentation in ${DOCS_OUT_USER}"
mkdir -p ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/user-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/group-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/resource-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/share-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/contract-* ${DOCS_OUT_USER}

## application-load-balancer
DOCS_OUT_ALB=${DOCS_OUT_ALB:-${DOCS_OUT}/application-load-balancer/}
echo "Move Application Load Balancer documentation in ${DOCS_OUT_ALB}"
mkdir -p ${DOCS_OUT_ALB}
mv -f docs/subcommands/tmp/applicationloadbalancer-* ${DOCS_OUT_ALB}
mv -f docs/subcommands/tmp/targetgroup-* ${DOCS_OUT_ALB}

## compute-engine
DOCS_OUT_COMPUTE_ENGINE=${DOCS_OUT_COMPUTE_ENGINE:-${DOCS_OUT}/compute-engine/}
echo "Move Compute Engine documentation in ${DOCS_OUT_COMPUTE_ENGINE}"
mkdir -p ${DOCS_OUT_COMPUTE_ENGINE}
mv -f docs/subcommands/tmp/* ${DOCS_OUT_COMPUTE_ENGINE}

# Remove temporary folder
rm -rf ${DOCS_OUT_TMP}
