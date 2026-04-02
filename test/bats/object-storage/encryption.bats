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
    export TEST_BUCKET_NAME="ionosctl-ci-enc-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --region "$TEST_REGION"
    assert_success

    echo "created test bucket for encryption tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        # Clean up any leftover encryption config before deleting the bucket
        run ionosctl object-storage encryption delete --name "$TEST_BUCKET_NAME" -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- validation ---

@test "object-storage encryption get: missing --name flag returns error" {
    run ionosctl object-storage encryption get 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage encryption put: missing --name flag returns error" {
    run ionosctl object-storage encryption put 2>&1
    assert_failure
}

@test "object-storage encryption put: missing --json-properties returns error" {
    run ionosctl object-storage encryption put --name some-bucket 2>&1
    assert_failure
    assert_output -p "requires at least 2 options"
}

@test "object-storage encryption delete: missing --name flag returns error" {
    run ionosctl object-storage encryption delete 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage encryption get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage encryption get --name some-bucket \
        --config /dev/null 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

@test "object-storage encryption put: --json-properties-example prints example JSON" {
    run ionosctl object-storage encryption put --json-properties-example 2>/dev/null
    assert_success
    assert_output -p "Rules"
    assert_output -p "SSEAlgorithm"
    assert_output -p "AES256"
}

@test "object-storage encryption put: nonexistent file returns error" {
    run ionosctl object-storage encryption put --name "$TEST_BUCKET_NAME" --json-properties "/tmp/nonexistent-encryption-file.json" 2>&1
    assert_failure
}

# --- encryption: put, get, delete lifecycle ---

@test "object-storage encryption put: apply AES256 encryption" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Rules":[{"ApplyServerSideEncryptionByDefault":{"SSEAlgorithm":"AES256"}}]}
EOF
    run ionosctl object-storage encryption put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage encryption get: retrieve encryption configuration" {
    run ionosctl object-storage encryption get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "AES256"
}

@test "object-storage encryption delete: remove encryption configuration" {
    run ionosctl object-storage encryption delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
    assert_output -p "deleted successfully"
}

@test "object-storage encryption get: after delete returns error" {
    run ionosctl object-storage encryption get --name "$TEST_BUCKET_NAME" 2>&1
    assert_failure
}

# --- encryption: from pretty-printed JSON file ---

@test "object-storage encryption put: apply from pretty-printed JSON file" {
    local tmpfile
    tmpfile="$(mktemp)"
    cat > "$tmpfile" <<ENC
{
  "Rules": [
    {
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }
  ]
}
ENC
    run ionosctl object-storage encryption put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage encryption get: file-based config applied correctly" {
    run ionosctl object-storage encryption get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "AES256"
}

@test "object-storage encryption delete: cleanup file-based config" {
    run ionosctl object-storage encryption delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}
