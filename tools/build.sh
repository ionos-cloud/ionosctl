#!/usr/bin/env bash

set -euo pipefail

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../"
OUT_D=${OUT_D:-${DIR}/builds}
mkdir -p "${OUT_D}"

VERSION="$(git tag -l | sort --version-sort | tail -n1 | cut -c 2-)"
GIT_COMMIT="$(git rev-parse --short HEAD)"
[[ -n $(git status -s) ]] && echo 'modified and/or untracked diff' && GIT_COMMIT="${GIT_COMMIT}.modified"

ldflags="-X github.com/ionos-cloud/ionosctl/commands.Version=${VERSION} -X github.com/ionos-cloud/ionosctl/commands.GitCommit=${GIT_COMMIT}"

(
    export GO111MODULE=on
    export CGO_ENABLED=0
    go build -ldflags "$ldflags" -o "${OUT_D}/ionosctl_${GOOS}_${GOARCH}"
)
