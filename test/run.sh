#!/bin/bash

# absolute directory where this script is located
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
BASE_BRANCH=${BASE_BRANCH:-master}
export LIBS_PATH="${SCRIPT_DIR}/libs"

# all .bats files under suites/, excluding setup.bats
BATS_FILES=$(find "${SCRIPT_DIR}/suites" -name '*.bats' ! -name 'setup.bats')

# modified files compared to the base
MODIFIED_FILES=$(git diff --name-only $BASE_BRANCH | grep -vE 'vendor|docs')

# If shared code is modified, run all test suites
RUN_ALL=false
for file in $MODIFIED_FILES; do
    if [[ "$file" == internal/* ]] || [[ "$file" == pkg/* ]] || \
       [[ "$file" == services/* ]] || [[ "$file" == commands/compute/helpers/* ]] || \
       [[ "$file" == commands/compute/waiter/* ]]; then
        echo "Shared code modified ($file) — running all test suites"
        RUN_ALL=true
        break
    fi
done

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

    if $RUN_ALL; then
        echo "Running $file_absolute_path (shared code changed)"
        bats "$file_absolute_path"
        if [ $? -ne 0 ]; then
            exit_code=1
        fi
        continue
    fi

    # Get path globs for the current file
    if ! paths_output=$("${SCRIPT_DIR}/parse_paths.sh" "$file_absolute_path"); then
        echo "Skipping $file_absolute_path (no paths directive)"
        continue
    fi
    read -ra path_array <<< "$paths_output"

    matched_glob=""
    matched_file=""
    should_run=false

    for glob in "${path_array[@]}"; do
        read match match_file <<< "$(matches_path "$glob")"
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
