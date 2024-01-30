#!/bin/bash

# Get the absolute directory where this script is located
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

BATS_FILES=$(find "${SCRIPT_DIR}/bats" -path "${SCRIPT_DIR}/bats/libs" -prune -o -name '*.bats' -print)

# Get a list of modified files compared to the master branch
MODIFIED_FILES=$(git diff --name-only master)

contains_tag() {
    tag=$1
    for file in $MODIFIED_FILES; do
        file_base_name=$(basename "$file")
        file_dir_name=$(dirname "$file")

        if [[ $file_base_name == *"$tag"* ]] || [[ $file_dir_name == *"$tag"* ]]; then
            echo "$tag $file"
            return 0
        fi
    done
    return 1
}

for file in $BATS_FILES; do
    # Get tags for the current file
    read -ra tag_array <<< $("$SCRIPT_DIR/bats/parse_tags.sh" "$file")
    matched_tag=""
    matched_file=""
    should_run=false

    # Check if any tag matches or is close to any modified file
    for tag in "${tag_array[@]}"; do
        read match match_file <<< $(contains_tag "$tag")
        if [[ -n $match ]]; then
            should_run=true
            matched_tag=$match
            matched_file=$match_file
            break
        fi
    done

    if $should_run; then
        echo "Running $file due to modified file: $matched_file (matched tag: $matched_tag)"
        bats "$file"
    else
        echo "Skipping $file because none of its tags match any modified files."
    fi
done
