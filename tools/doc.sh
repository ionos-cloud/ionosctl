#!/usr/bin/env bash

## Until a specific split on services is made,
## the documentation for the CLI is generated using
## the current script.

set -euo pipefail

DOCS_OUT=${DOCS_OUT:-$(shell pwd)/docs/subcommands}

## Generate all documentation in one place
## Use temporary files
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
mv -f docs/subcommands/tmp/contract-* ${DOCS_OUT_AUTH}
mv -f docs/subcommands/tmp/login.md ${DOCS_OUT_AUTH}

## Kubernetes
DOCS_OUT_K8S=${DOCS_OUT_K8S:-${DOCS_OUT}/kubernetes/}
echo "Move Kubernetes documentation in ${DOCS_OUT_K8S}"
mkdir -p ${DOCS_OUT_K8S}
mv -f docs/subcommands/tmp/k8s-* ${DOCS_OUT_K8S}

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

## Compute Engine
DOCS_OUT_COMPUTE_ENGINE=${DOCS_OUT_COMPUTE_ENGINE:-${DOCS_OUT}/compute-engine/}
echo "Move Compute Engine documentation in ${DOCS_OUT_COMPUTE_ENGINE}"
mkdir -p ${DOCS_OUT_COMPUTE_ENGINE}
mv -f docs/subcommands/tmp/* ${DOCS_OUT_COMPUTE_ENGINE}

# Remove temporary folder
rm -rf ${DOCS_OUT_TMP}
