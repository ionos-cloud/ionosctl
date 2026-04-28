#!/bin/bash

# absolute directory where this script is located
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
BASE_BRANCH=${BASE_BRANCH:-master}
export LIBS_PATH="${SCRIPT_DIR}/libs"

# all .bats files under suites/, excluding setup.bats
BATS_FILES=$(find "${SCRIPT_DIR}/suites" -name '*.bats' ! -name 'setup.bats')

# modified files compared to the base
MODIFIED_FILES=$(git diff --name-only $BASE_BRANCH | grep -vE 'vendor|docs')

# Check if any modified file matches a glob pattern.
# Globs use * to match any path segment (e.g. commands/cdn/* matches commands/cdn/foo/bar.go).
# Returns the matched glob and file on stdout, exit code 0 on match.
matches_path() {
    local glob=$1
    for file in $MODIFIED_FILES; do
        # Use bash pattern matching: convert glob/* to match any depth
        # commands/cdn/* should match commands/cdn/distribution/routingrules/list.go
        local pattern="${glob%/\*}/"
        if [[ "$file" == ${pattern}* ]]; then
            echo "$glob $file"
            return 0
        fi
    done
    return 1
}

exit_code=0
for file in $BATS_FILES; do
    file_absolute_path=$(realpath "$file")

    # Get path globs for the current file
    read -ra path_array <<< $("${SCRIPT_DIR}/parse_paths.sh" "$file_absolute_path")

    matched_glob=""
    matched_file=""
    should_run=false

    for glob in "${path_array[@]}"; do
        read match match_file <<< $(matches_path "$glob")
        if [[ -n $match ]]; then
            should_run=true
            matched_glob=$match
            matched_file=$match_file
            break
        fi
    done

    if $should_run; then
        echo "Running $file_absolute_path due to modified file: $matched_file (matched path: $matched_glob)"
        bats "$file_absolute_path"
        if [ $? -ne 0 ]; then
            exit_code=1
        fi
    else
        echo "Skipping $file_absolute_path (no path matches)"
    fi
done

exit $exit_code
