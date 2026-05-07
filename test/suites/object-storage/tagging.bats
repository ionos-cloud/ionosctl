#!/usr/bin/env bats

# paths: commands/object-storage/*

load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-tag-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION"
    assert_success

    echo "created test bucket for tagging tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        run ionosctl object-storage bucket tagging delete --name "$TEST_BUCKET_NAME" -f
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- validation ---

@test "object-storage bucket tagging get: missing --name flag returns error" {
    run ionosctl object-storage bucket tagging get
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket tagging put: missing --name flag returns error" {
    run ionosctl object-storage bucket tagging put
    assert_failure
}

@test "object-storage bucket tagging put: missing --json-properties returns error" {
    run ionosctl object-storage bucket tagging put --name some-bucket
    assert_failure
}

@test "object-storage bucket tagging delete: missing --name flag returns error" {
    run ionosctl object-storage bucket tagging delete
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket tagging get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket tagging get --name some-bucket 
    assert_failure
    assert_stderr -p "object storage credentials not found"
}

@test "object-storage bucket tagging put: --json-properties-example prints example JSON" {
    run ionosctl object-storage bucket tagging put --json-properties-example
    assert_success
    assert_output -p "TagSet"
    assert_output -p "Environment"
}

@test "object-storage bucket tagging put: nonexistent file returns error" {
    run ionosctl object-storage bucket tagging put --name "$TEST_BUCKET_NAME" --json-properties "/tmp/nonexistent-tags.json"
    assert_failure
}

# --- lifecycle: put, get, delete ---

@test "object-storage bucket tagging put: apply tags" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"TagSet":[{"Key":"Environment","Value":"test"},{"Key":"Team","Value":"platform"}]}
EOF
    run ionosctl object-storage bucket tagging put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile"
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage bucket tagging get: retrieve tags" {
    run ionosctl object-storage bucket tagging get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "Environment"
    assert_output -p "test"
    assert_output -p "Team"
    assert_output -p "platform"
}

@test "object-storage bucket tagging get: json output" {
    run ionosctl object-storage bucket tagging get --name "$TEST_BUCKET_NAME" -o json
    assert_success
    echo "$output" | jq -e '.items[0].Key // .[0].Key' >/dev/null
    echo "$output" | jq -e '.items[0].Value // .[0].Value' >/dev/null
}

@test "object-storage bucket tagging delete: remove tags" {
    run ionosctl object-storage bucket tagging delete --name "$TEST_BUCKET_NAME" -f
    assert_success
    assert_output -p "deleted successfully"
}

@test "object-storage bucket tagging get: after delete returns error" {
    run ionosctl object-storage bucket tagging get --name "$TEST_BUCKET_NAME"
    assert_failure
}
