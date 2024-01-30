#!/bin/bash

# This script takes a .bats file as input and outputs its tags in space-separated format
# Tags for .bats files should follow this format: '# tags: tag1, tag2, tag3'
# Tags for .bats files should ideally reside at the start of the file
# Tags should ideally represent the resources which are tested in that file
# Output: tag1 tag2 tag3
# Exit code 0: Tags outputted successfully
# Exit code 1: Usage error - file not found
# Exit code 2: No tags found in file

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

# Extract the first line that contains tags
TAG_LINE=$(grep -i "^# tags:" "$BATS_FILE" | head -n 1)

if [ -z "$TAG_LINE" ]; then
    echo "No tags found in $BATS_FILE"
    exit 2
fi

# Remove the '# tags:' prefix and replace commas with spaces
TAGS=${TAG_LINE/#\# tags: /}
TAGS=${TAGS//, / }

echo "$TAGS"
