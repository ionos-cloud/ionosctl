#!/bin/bash

# Regenerate Mocks for all services

GO111MODULE=off go get -d github.com/golang/mock/mockgen

# Loop through all service directories
for dir in $(find services -type d); do
    echo "$dir:"
    pushd "$dir" >/dev/null || exit

    if [ -d "resources" ]; then
        pushd "resources" >/dev/null || exit

        for file in $(find . -name '*.go' ! -name '*_test.go' -not -path './mocks/*'); do
            if [ -f "$file" ]; then
                echo "  - $file"
                mkdir -p mocks

                # Capitalize first letter of filename
                filename=$(basename "$file" .go)
                filename=$(echo "${filename^}")

                mockgen -source "$file" >"mocks/${filename}Service.go"
            fi
        done

        popd >/dev/null || exit
    fi

    popd >/dev/null || exit
done

# Go back to the root directory
pushd "$(dirname "$0")"/.. >/dev/null || exit

popd >/dev/null || exit
