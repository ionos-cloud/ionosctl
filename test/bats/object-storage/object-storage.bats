#!/usr/bin/env bats

# tags: object-storage

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --region "$TEST_REGION"
    assert_success

    echo "created test bucket: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" --recursive -f
    fi
}

# --- list ---

@test "object-storage bucket list: returns test bucket" {
    run ionosctl object-storage bucket list 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
}

@test "object-storage bucket list: --region filters by region" {
    run ionosctl object-storage bucket list --region "$TEST_REGION" 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
}

@test "object-storage bucket list: --region with wrong region excludes test bucket" {
    run ionosctl object-storage bucket list --region "eu-south-2" 2>/dev/null
    assert_success
    refute_output -p "$TEST_BUCKET_NAME"
}

@test "object-storage bucket list: json output" {
    run ionosctl object-storage bucket list -o json 2>/dev/null
    assert_success
    echo "$output" | jq -e '.items[0].Name // .[0].Name // .Name' >/dev/null
}

@test "object-storage bucket list: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket list \
        --config /dev/null 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

# --- get ---

@test "object-storage bucket get: missing --name flag returns error" {
    run ionosctl object-storage bucket get 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage bucket create: missing --name flag returns error" {
    run ionosctl object-storage bucket create 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage bucket get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket get --name some-bucket \
        --config /dev/null 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

@test "object-storage bucket get: retrieve existing bucket" {
    run ionosctl object-storage bucket get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
}

@test "object-storage bucket get: bucket not found returns error" {
    run ionosctl object-storage bucket get --name "nonexistent-bucket-$(randStr 10 | tr '[:upper:]' '[:lower:]')" 2>&1
    assert_failure
    assert_output -p "not found"
}

@test "object-storage bucket get: json output contains expected fields" {
    run ionosctl object-storage bucket get --name "$TEST_BUCKET_NAME" -o json 2>/dev/null
    assert_success

    echo "$output" | jq -e '.Name' >/dev/null
    echo "$output" | jq -e '.CreationDate' >/dev/null
}

@test "object-storage bucket head: missing --name flag returns error" {
    run ionosctl object-storage bucket head 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage bucket head: existing bucket is accessible" {
    run ionosctl object-storage bucket head --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "exists and is accessible"
}

@test "object-storage bucket head: nonexistent bucket returns error" {
    run ionosctl object-storage bucket head --name "nonexistent-bucket-$(randStr 10 | tr '[:upper:]' '[:lower:]')" 2>&1
    assert_failure
}

@test "object-storage bucket head: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket head --name some-bucket \
        --config /dev/null 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

@test "object-storage bucket versioning get: missing --name flag returns error" {
    run ionosctl object-storage bucket versioning get 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage bucket versioning get: returns versioning state" {
    run ionosctl object-storage bucket versioning get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
}

@test "object-storage bucket versioning get: new bucket has versioning disabled" {
    run ionosctl object-storage bucket versioning get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "Disabled"
}

@test "object-storage bucket versioning get: json output" {
    run ionosctl object-storage bucket versioning get --name "$TEST_BUCKET_NAME" -o json 2>/dev/null
    assert_success
    echo "$output" | jq -e '.Versioning' >/dev/null
}

@test "object-storage bucket versioning set: missing flags returns error" {
    run ionosctl object-storage bucket versioning set 2>&1
    assert_failure
    assert_output -p "requires at least 2 option"
}

@test "object-storage bucket versioning set: missing --status flag returns error" {
    run ionosctl object-storage bucket versioning set --name "$TEST_BUCKET_NAME" 2>&1
    assert_failure
}

@test "object-storage bucket versioning set: enable versioning" {
    run ionosctl object-storage bucket versioning set --name "$TEST_BUCKET_NAME" --status Enabled 2>/dev/null
    assert_success
    assert_output -p "Enabled"

    # Verify with get
    run ionosctl object-storage bucket versioning get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "Enabled"
}

@test "object-storage bucket versioning set: suspend versioning" {
    run ionosctl object-storage bucket versioning set --name "$TEST_BUCKET_NAME" --status Suspended 2>/dev/null
    assert_success
    assert_output -p "Suspended"

    # Verify with get
    run ionosctl object-storage bucket versioning get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "Suspended"
}

@test "object-storage bucket head: json output" {
    run ionosctl object-storage bucket head --name "$TEST_BUCKET_NAME" -o json 2>/dev/null
    assert_success
    echo "$output" | jq -e '.Status' >/dev/null
}

@test "object-storage bucket list-objects: missing --name flag returns error" {
    run ionosctl object-storage bucket list-objects 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage bucket list-objects: empty bucket returns no objects" {
    run ionosctl object-storage bucket list-objects --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "No objects found"
}

@test "object-storage bucket list-objects: nonexistent prefix returns no objects" {
    run ionosctl object-storage bucket list-objects --name "$TEST_BUCKET_NAME" --prefix "nonexistent-prefix/" 2>/dev/null
    assert_success
    assert_output -p "No objects found"
}

# --- delete ---

@test "object-storage bucket delete: missing --name flag returns error" {
    run ionosctl object-storage bucket delete 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage bucket delete: without --force prompts and fails non-interactively" {
    run bash -c 'echo "n" | ionosctl object-storage bucket delete --name "some-bucket" 2>&1'
    assert_failure
}

@test "object-storage bucket delete: --force skips prompt on nonexistent bucket" {
    run ionosctl object-storage bucket delete --name "nonexistent-bucket-$(randStr 10 | tr '[:upper:]' '[:lower:]')" -f 2>&1
    assert_failure
}

@test "object-storage bucket delete: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket delete --name some-bucket -f \
        --config /dev/null 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

# --- output ---

@test "object-storage bucket get: --cols flag filters output columns" {
    run ionosctl object-storage bucket get --name "$TEST_BUCKET_NAME" --cols Name 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
    # CreationDate column should not appear in header
    refute_output -p "CREATIONDATE"
}
