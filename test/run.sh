#!/bin/bash

# absolute directory where this script is located
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)
BASE_BRANCH=${BASE_BRANCH:-master}
export LIBS_PATH="${SCRIPT_DIR}/bats/libs"

# all .bats files excluding the 'libs' directory within 'bats'
BATS_FILES=$(find "${SCRIPT_DIR}/bats" -path "${SCRIPT_DIR}/bats/libs" -prune -o -name '*.bats' -print)

# modified files compared to the base
MODIFIED_FILES=$(git diff --name-only $BASE_BRANCH)

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
    # relative path to absolute path
    file_absolute_path=$(realpath "$file")

    # Get tags for the current file using absolute path
    read -ra tag_array <<< $("${SCRIPT_DIR}/bats/parse_tags.sh" "$file_absolute_path")

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
        echo "Running $file_absolute_path due to modified file: $matched_file (matched tag: $matched_tag)"
        bats "$file_absolute_path"
    else
        echo "Skipping $file_absolute_path because none of its tags match any modified files."
    fi
done
