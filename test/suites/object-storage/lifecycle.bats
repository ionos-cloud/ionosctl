#!/usr/bin/env bats

# paths: commands/object-storage/*

load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-lc-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION"
    assert_success

    echo "created test bucket for lifecycle tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        run ionosctl object-storage bucket lifecycle delete --name "$TEST_BUCKET_NAME" -f
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- validation ---

@test "object-storage bucket lifecycle get: missing --name flag returns error" {
    run ionosctl object-storage bucket lifecycle get
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket lifecycle put: missing --name flag returns error" {
    run ionosctl object-storage bucket lifecycle put
    assert_failure
}

@test "object-storage bucket lifecycle put: missing --json-properties returns error" {
    run ionosctl object-storage bucket lifecycle put --name some-bucket
    assert_failure
}

@test "object-storage bucket lifecycle delete: missing --name flag returns error" {
    run ionosctl object-storage bucket lifecycle delete
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket lifecycle get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket lifecycle get --name some-bucket
    assert_failure
    assert_stderr -p "object storage credentials not found"
}

@test "object-storage bucket lifecycle put: --json-properties-example prints example JSON" {
    run ionosctl object-storage bucket lifecycle put --json-properties-example
    assert_success
    assert_output -p "Rules"
    assert_output -p "Expiration"
    assert_output -p "Days"
}

@test "object-storage bucket lifecycle put: nonexistent file returns error" {
    run ionosctl object-storage bucket lifecycle put --name "$TEST_BUCKET_NAME" --json-properties "/tmp/nonexistent-lifecycle.json"
    assert_failure
}

# --- lifecycle: put, get, delete ---
# SAFETY: Rules use Days: 90 — they will NOT trigger during the test run.
# No versioning, no object-lock, no retention.

@test "object-storage bucket lifecycle put: apply lifecycle rules" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Rules":[{"ID":"expire-test","Prefix":"logs/","Status":"Enabled","Expiration":{"Days":90}},{"ID":"abort-uploads","Prefix":"","Status":"Enabled","AbortIncompleteMultipartUpload":{"DaysAfterInitiation":7}}]}
EOF
    run ionosctl object-storage bucket lifecycle put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile"
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage bucket lifecycle get: retrieve lifecycle rules" {
    run ionosctl object-storage bucket lifecycle get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "expire-test"
    assert_output -p "Enabled"
}

@test "object-storage bucket lifecycle get: json output" {
    run ionosctl object-storage bucket lifecycle get --name "$TEST_BUCKET_NAME" -o json
    assert_success
    echo "$output" | jq -e '.items[0].ID' >/dev/null
}

@test "object-storage bucket lifecycle delete: remove lifecycle configuration" {
    run ionosctl object-storage bucket lifecycle delete --name "$TEST_BUCKET_NAME" -f
    assert_success
    assert_output -p "deleted successfully"
}

@test "object-storage bucket lifecycle get: after delete returns error" {
    run ionosctl object-storage bucket lifecycle get --name "$TEST_BUCKET_NAME"
    assert_failure
}
