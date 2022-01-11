#!/usr/bin/env bash

## Temporarily, the documentation for the CLI is updated using the current script.

set -euo pipefail
DOCS_OUT=${DOCS_OUT:-$(shell pwd)/docs/subcommands}

## Generate all documentation in one place
DOCS_OUT_TMP=${DOCS_OUT_TMP:-${DOCS_OUT}/tmp}
echo "Generate Markdown documentation in ${DOCS_OUT_TMP}"
mkdir -p ${DOCS_OUT_TMP}
DOCS_OUT=${DOCS_OUT_TMP} go run tools/doc.go

## Separate documentation based on service

## CLI Setup
DOCS_OUT_SETUP=${DOCS_OUT_SETUP:-${DOCS_OUT}/tool/}
echo "Move CLI Setup documentation in ${DOCS_OUT_SETUP}"
mkdir -p ${DOCS_OUT_SETUP}
# mv -f docs/subcommands/completion-* ${DOCS_OUT_SETUP}
mv -f docs/subcommands/tmp/version.md ${DOCS_OUT_SETUP}
cp docs/subcommands/tmp/login.md ${DOCS_OUT_SETUP}

## Authentication
DOCS_OUT_AUTH=${DOCS_OUT_AUTH:-${DOCS_OUT}/authentication/}
echo "Move Authentication documentation in ${DOCS_OUT_AUTH}"
mkdir -p ${DOCS_OUT_AUTH}
mv -f docs/subcommands/tmp/token-* ${DOCS_OUT_AUTH}
cp docs/subcommands/tmp/contract-* ${DOCS_OUT_AUTH}
mv -f docs/subcommands/tmp/login.md ${DOCS_OUT_AUTH}

## Kubernetes
DOCS_OUT_K8S=${DOCS_OUT_K8S:-${DOCS_OUT}/kubernetes/}
echo "Move Kubernetes documentation in ${DOCS_OUT_K8S}"
mkdir -p ${DOCS_OUT_K8S}
mv -f docs/subcommands/tmp/k8s-* ${DOCS_OUT_K8S}

## Network Load Balancer
DOCS_OUT_NLB=${DOCS_OUT_NLB:-${DOCS_OUT}/networkloadbalancer/}
echo "Move Network Load Balancer documentation in ${DOCS_OUT_NLB}"
mkdir -p ${DOCS_OUT_NLB}
mv -f docs/subcommands/tmp/networkloadbalancer-* ${DOCS_OUT_NLB}

## NAT Gateway
DOCS_OUT_NAT_GATEWAY=${DOCS_OUT_NAT_GATEWAY:-${DOCS_OUT}/natgateway/}
echo "Move NAT Gateway documentation in ${DOCS_OUT_NAT_GATEWAY}"
mkdir -p ${DOCS_OUT_NAT_GATEWAY}
mv -f docs/subcommands/tmp/natgateway-* ${DOCS_OUT_NAT_GATEWAY}

## Backup
DOCS_OUT_BACKUP=${DOCS_OUT_BACKUP:-${DOCS_OUT}/backup/}
echo "Move Backup documentation in ${DOCS_OUT_BACKUP}"
mkdir -p ${DOCS_OUT_BACKUP}
mv -f docs/subcommands/tmp/backupunit-* ${DOCS_OUT_BACKUP}

## User Management
DOCS_OUT_USER=${DOCS_OUT_USER:-${DOCS_OUT}/user/}
echo "Move User documentation in ${DOCS_OUT_USER}"
mkdir -p ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/user-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/group-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/resource-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/share-* ${DOCS_OUT_USER}
mv -f docs/subcommands/tmp/contract-* ${DOCS_OUT_USER}

## Compute Engine
DOCS_OUT_COMPUTE_ENGINE=${DOCS_OUT_COMPUTE_ENGINE:-${DOCS_OUT}/compute-engine/}
echo "Move Compute Engine documentation in ${DOCS_OUT_COMPUTE_ENGINE}"
mkdir -p ${DOCS_OUT_COMPUTE_ENGINE}
mv -f docs/subcommands/tmp/* ${DOCS_OUT_COMPUTE_ENGINE}

# Remove temporary folder
rm -rf ${DOCS_OUT_TMP}
