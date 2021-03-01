#!/usr/bin/env bash

set -euo pipefail

version="$(git tag -l | sort --version-sort | tail -n1 | cut -c 2-)"

major="$(echo "$version" | cut -d . -f1)"
minor="$(echo "$version" | cut -d . -f2)"
patch="$(echo "$version" | cut -d . -f3)"

ldflags="-X github.com/ionos-cloud/ionosctl/commands.Major=${major}"
ldflags="${ldflags} -X github.com/ionos-cloud/ionosctl/commands.Minor=${minor}"
ldflags="${ldflags} -X github.com/ionos-cloud/ionosctl/commands.Patch=${patch}"

(
    export GO111MODULE=on
    export CGO_ENABLED=0
    go install -ldflags "$ldflags"
)
