#!/usr/bin/env bash

set -euo pipefail

VERSION="$(git tag -l | sort --version-sort | tail -n1 | cut -c 2-)"
GIT_COMMIT="$(git rev-parse --short HEAD)"
[[ -n $(git status -s) ]] && echo 'modified and/or untracked diff' && GIT_COMMIT="${GIT_COMMIT}.modified"

ldflags="-X github.com/ionos-cloud/ionosctl/v6/commands.Version=${VERSION} -X github.com/ionos-cloud/ionosctl/v6/commands.GitCommit=${GIT_COMMIT}"

(
    export GO111MODULE=on
    export CGO_ENABLED=0
    go install -ldflags "$ldflags"
)
