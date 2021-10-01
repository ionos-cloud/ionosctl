#!/usr/bin/env bash

set -euo pipefail

version="$(git tag -l | sort --version-sort | tail -n1 | cut -c 2-)"

ldflags="-X github.com/ionos-cloud/ionosctl/commands.Version=${version}"

(
    export GO111MODULE=on
    export CGO_ENABLED=0
    go install -ldflags "$ldflags"
)
