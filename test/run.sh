#!/bin/bash

# Get the absolute directory where this script is located
SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" &>/dev/null && pwd)

BATS_FILES=$(find "${SCRIPT_DIR}/bats" -path "${SCRIPT_DIR}/bats/libs" -prune -o -name '*.bats' -print)
echo "$BATS_FILES"

contains_tag() {
    for file in $MODIFIED_FILES; do
        if [[ $file == *"$1"* ]]; then
            return 0
        fi
    done
    return 1
}

# Get a list of modified files compared to the master branch
MODIFIED_FILES=$(git diff --name-only master)

for file in $BATS_FILES; do
    # Get tags for the current file
    read -ra tag_array <<< $("$SCRIPT_DIR/bats/parse_tags.sh" "$file")

    # Check if any tag matches or is close to any modified file
    for tag in "${tag_array[@]}"; do
        if contains_tag "$tag"; then
            # Run the .bats file
            bats "$file"
            break
        fi
    done
done
