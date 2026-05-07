#!/usr/bin/env bats

# paths: commands/object-storage/*

load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-obj-$(randStr 8 | tr '[:upper:]' '[:lower:]')"
    export TEST_KEY="test-object.txt"
    export TEST_FILE="/tmp/ionosctl-test-upload-$$.txt"
    export TEST_DOWNLOAD="/tmp/ionosctl-test-download-$$.txt"

    echo "hello ionosctl object test" > "$TEST_FILE"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION"
    assert_success

    echo "created test bucket for object tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    rm -f "$TEST_FILE" "$TEST_DOWNLOAD"
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- put: validation ---

@test "object-storage object put: missing --name flag returns error" {
    run ionosctl object-storage object put --key foo --source "$TEST_FILE"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object put: missing --key flag returns error" {
    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --source "$TEST_FILE"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object put: missing --source flag returns error" {
    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key foo
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object put: nonexistent source file returns error" {
    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key foo --source /nonexistent
    assert_failure
    assert_stderr -p "opening source file"
}

@test "object-storage object put: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key foo --source "$TEST_FILE" 
    assert_failure
    assert_stderr -p "object storage credentials not found"
}

# --- put: success ---

@test "object-storage object put: upload file succeeds" {
    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --source "$TEST_FILE"
    assert_success
    assert_output -p "uploaded to bucket"
}

# --- head ---

@test "object-storage object head: missing --name flag returns error" {
    run ionosctl object-storage object head --key foo
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object head: missing --key flag returns error" {
    run ionosctl object-storage object head --name "$TEST_BUCKET_NAME"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object head: returns metadata for uploaded object" {
    run ionosctl object-storage object head --name "$TEST_BUCKET_NAME" --key "$TEST_KEY"
    assert_success
    assert_output -p "$TEST_KEY"
}

@test "object-storage object head: json output contains expected fields" {
    run ionosctl object-storage object head --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" -o json
    assert_success
    echo "$output" | jq -e '.Key' >/dev/null
    echo "$output" | jq -e '.ContentType' >/dev/null
    echo "$output" | jq -e '.ETag' >/dev/null
}

@test "object-storage object head: nonexistent key returns error" {
    run ionosctl object-storage object head --name "$TEST_BUCKET_NAME" --key "nonexistent-key-$(randStr 10)"
    assert_failure
}

# --- get ---

@test "object-storage object get: missing --name flag returns error" {
    run ionosctl object-storage object get --key foo
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object get: missing --key flag returns error" {
    run ionosctl object-storage object get --name "$TEST_BUCKET_NAME"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object get: download file succeeds" {
    run ionosctl object-storage object get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --destination "$TEST_DOWNLOAD"
    assert_success
    assert_output -p "downloaded to"
    [ -f "$TEST_DOWNLOAD" ]
    diff "$TEST_FILE" "$TEST_DOWNLOAD"
}

@test "object-storage object get: nonexistent key returns error" {
    run ionosctl object-storage object get --name "$TEST_BUCKET_NAME" --key "nonexistent-key-$(randStr 10)" --destination /tmp/nope
    assert_failure
}

# --- copy ---

@test "object-storage object copy: missing --name flag returns error" {
    run ionosctl object-storage object copy --key foo --copy-source "/src/key"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object copy: missing --key flag returns error" {
    run ionosctl object-storage object copy --name "$TEST_BUCKET_NAME" --copy-source "/src/key"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object copy: missing --copy-source flag returns error" {
    run ionosctl object-storage object copy --name "$TEST_BUCKET_NAME" --key "dest"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object copy: copy object within same bucket succeeds" {
    run ionosctl object-storage object copy \
        --name "$TEST_BUCKET_NAME" \
        --key "copy-of-${TEST_KEY}" \
        --copy-source "/${TEST_BUCKET_NAME}/${TEST_KEY}"
    assert_success
    assert_output -p "ETag"
}

# --- delete ---

@test "object-storage object delete: missing --name flag returns error" {
    run ionosctl object-storage object delete --key foo -f
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object delete: missing --key and --all flags returns error" {
    run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" -f
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object delete: --all without --name returns error" {
    run ionosctl object-storage object delete --all -f
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object delete: delete object succeeds" {
    run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" -f
    assert_success
    assert_output -p "deleted from bucket"
}

@test "object-storage object delete: delete copied object succeeds" {
    run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --key "copy-of-${TEST_KEY}" -f
    assert_success
    assert_output -p "deleted from bucket"
}

@test "object-storage object delete: --all deletes all objects in bucket" {
    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key "all-test.txt" --source "$TEST_FILE"
    assert_success

    run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
    assert_success
    assert_output -p "All objects deleted"

    # Verify bucket is empty
    run ionosctl object-storage object list --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "No objects found"
}
