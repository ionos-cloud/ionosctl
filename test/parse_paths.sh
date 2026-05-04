#!/bin/bash

# This script takes a .bats file as input and outputs its trigger path globs
# Path globs for .bats files should follow this format: '# paths: glob1, glob2, glob3'
# Path globs should ideally reside at the start of the file
# Output: glob1 glob2 glob3 (comma-separated globs split into space-separated)
# Exit code 0: Paths outputted successfully
# Exit code 1: Usage error - file not found
# Exit code 2: No paths found in file

BATS_FILE=$1

if [ $# -eq 0 ]; then
    echo "No .bats file provided"
    echo "Usage: $0 path/to/file.bats"
    exit 1
fi

if [ ! -f "$BATS_FILE" ]; then
    echo "File not found: $BATS_FILE"
    exit 1
fi

# Extract the first line that contains paths
PATH_LINE=$(grep "^# paths:" "$BATS_FILE" | head -n 1)

if [ -z "$PATH_LINE" ]; then
    echo "No paths found in $BATS_FILE"
    exit 2
fi

# Remove the '# paths: ' prefix and replace ', ' with spaces
PATHS=${PATH_LINE/#\# paths: /}
PATHS=${PATHS//, / }

echo "$PATHS"
