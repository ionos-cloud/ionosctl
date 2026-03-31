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
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" --region "$TEST_REGION" -f
    fi
}

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
    run ionosctl object-storage bucket get --name "$TEST_BUCKET_NAME" --region "$TEST_REGION" 2>/dev/null
    assert_success
    assert_output -p "$TEST_BUCKET_NAME"
}

@test "object-storage bucket get: bucket not found returns error" {
    run ionosctl object-storage bucket get --name "nonexistent-bucket-$(randStr 10 | tr '[:upper:]' '[:lower:]')" --region "$TEST_REGION" 2>&1
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
