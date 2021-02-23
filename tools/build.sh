#!/usr/bin/env bash

set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"
OUT_D=${OUT_D:-${DIR}/builds}
mkdir -p "${OUT_D}"

(
    export GO111MODULE=on
    export CGO_ENABLED=0
    go build -o "${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
)
