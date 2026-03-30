#!/usr/bin/env bats

# tags: object-storage

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'


setup_file() {
    # S3 credentials required for object storage e2e tests
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        skip "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests"
    fi

    export TEST_BUCKET_NAME="${IONOS_S3_TEST_BUCKET:-}" # pre-existing bucket name for read-only tests
    if [[ -z "$TEST_BUCKET_NAME" ]]; then
        skip "IONOS_S3_TEST_BUCKET must be set to the name of an existing bucket for e2e tests"
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
}

@test "object-storage bucket get: missing --name flag returns error" {
    run ionosctl object-storage bucket get 2>&1
    assert_failure
    assert_output -p "flag(s) not set: name"
}

@test "object-storage bucket get: missing S3 credentials returns error" {
    # Temporarily unset credentials to test the error path
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket get --name some-bucket \
        --config /dev/null 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

@test "object-storage bucket get: retrieve existing bucket" {
    run ionosctl object-storage bucket get --name "$TEST_BUCKET_NAME" --region "$TEST_REGION" 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
}

@test "object-storage bucket get: bucket not found returns error" {
    run ionosctl object-storage bucket get --name "nonexistent-bucket-$(randStr 10)" --region "$TEST_REGION" 2>&1
    assert_failure
    assert_output -p "not found"
}

@test "object-storage bucket get: json output contains expected fields" {
    run ionosctl object-storage bucket get --name "$TEST_BUCKET_NAME" --region "$TEST_REGION" -o json 2>/dev/null
    assert_success

    echo "$output" | jq -e '.Name' >/dev/null
    echo "$output" | jq -e '.CreationDate' >/dev/null
}

@test "object-storage bucket get: --cols flag filters output columns" {
    run ionosctl object-storage bucket get --name "$TEST_BUCKET_NAME" --region "$TEST_REGION" --cols Name 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
    # CreationDate column should not appear in header
    refute_output -p "CREATIONDATE"
}