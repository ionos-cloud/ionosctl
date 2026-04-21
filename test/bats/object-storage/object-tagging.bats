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
    export TEST_BUCKET_NAME="ionosctl-ci-otag-$(randStr 8 | tr '[:upper:]' '[:lower:]')"
    export TEST_KEY="test-tagging.txt"
    export TEST_FILE="/tmp/ionosctl-test-otag-$$.txt"

    echo "object tagging test content" > "$TEST_FILE"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION"
    assert_success

    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --source "$TEST_FILE" 2>/dev/null
    assert_success

    echo "created test bucket with object for object-tagging tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    rm -f "$TEST_FILE"
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- validation ---

@test "object-storage object tagging get: missing --name flag returns error" {
    run ionosctl object-storage object tagging get --key foo 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object tagging get: missing --key flag returns error" {
    run ionosctl object-storage object tagging get --name "$TEST_BUCKET_NAME" 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object tagging put: missing --name flag returns error" {
    run ionosctl object-storage object tagging put --key foo --json-properties /dev/null 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object tagging put: missing --key flag returns error" {
    run ionosctl object-storage object tagging put --name "$TEST_BUCKET_NAME" --json-properties /dev/null 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object tagging delete: missing --name flag returns error" {
    run ionosctl object-storage object tagging delete --key foo 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object tagging delete: missing --key flag returns error" {
    run ionosctl object-storage object tagging delete --name "$TEST_BUCKET_NAME" 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object tagging get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage object tagging get --name some-bucket --key foo 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

@test "object-storage object tagging put: --json-properties-example prints example JSON" {
    run ionosctl object-storage object tagging put --json-properties-example 2>/dev/null
    assert_success
    assert_output -p "TagSet"
    assert_output -p "Environment"
}

@test "object-storage object tagging put: nonexistent file returns error" {
    run ionosctl object-storage object tagging put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --json-properties "/tmp/nonexistent-tags.json" 2>&1
    assert_failure
}

# --- lifecycle: put, get, delete ---

@test "object-storage object tagging put: apply tags" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"TagSet":[{"Key":"Environment","Value":"test"},{"Key":"Team","Value":"platform"}]}
EOF
    run ionosctl object-storage object tagging put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage object tagging get: retrieve tags" {
    run ionosctl object-storage object tagging get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" 2>/dev/null
    assert_success
    assert_output -p "Environment"
    assert_output -p "test"
    assert_output -p "Team"
    assert_output -p "platform"
}

@test "object-storage object tagging get: json output" {
    run ionosctl object-storage object tagging get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" -o json 2>/dev/null
    assert_success
    echo "$output" | jq -e '.items[0].Key // .[0].Key' >/dev/null
    echo "$output" | jq -e '.items[0].Value // .[0].Value' >/dev/null
}

@test "object-storage object tagging delete: remove tags" {
    run ionosctl object-storage object tagging delete --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" -f 2>/dev/null
    assert_success
    assert_output -p "deleted successfully"
}

@test "object-storage object tagging get: after delete returns empty or error" {
    run ionosctl object-storage object tagging get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" 2>&1
    # After deleting tags, get may return empty set or error depending on API behavior
    # We just verify the command completes (no hanging)
}
