#!/usr/bin/env bash

set -euo pipefail

version="$(git tag -l | sort --version-sort | tail -n1 | cut -c 2-)"
GIT_COMMIT="$(git rev-parse --short HEAD)"

ldflags="-X github.com/ionos-cloud/ionosctl/commands.Version=${version} -X github.com/ionos-cloud/ionosctl/commands.GitCommit=${GIT_COMMIT}"

(
    export GO111MODULE=on
    export CGO_ENABLED=0
    go install -ldflags "$ldflags"
)
