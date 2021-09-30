#!/usr/bin/env bash

set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"
OUT_D=${OUT_D:-${DIR}/builds}
mkdir -p "${OUT_D}"

version="$(git tag -l | sort --version-sort | tail -n1 | cut -c 2-)"

ldflags="-X github.com/ionos-cloud/ionosctl/commands.Version=${version}"

(
    export GO111MODULE=on
    export CGO_ENABLED=0
    go build -ldflags "$ldflags" -o "${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
)
