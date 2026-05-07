#!/usr/bin/env bats

# paths: commands/object-storage/*

load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-pab-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION"
    assert_success

    echo "created test bucket for public-access-block tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        run ionosctl object-storage bucket public-access-block delete --name "$TEST_BUCKET_NAME" -f
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- validation ---

@test "object-storage bucket public-access-block get: missing --name flag returns error" {
    run ionosctl object-storage bucket public-access-block get
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket public-access-block put: missing --name flag returns error" {
    run ionosctl object-storage bucket public-access-block put
    assert_failure
}

@test "object-storage bucket public-access-block put: missing --json-properties returns error" {
    run ionosctl object-storage bucket public-access-block put --name some-bucket
    assert_failure
}

@test "object-storage bucket public-access-block delete: missing --name flag returns error" {
    run ionosctl object-storage bucket public-access-block delete
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket public-access-block get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket public-access-block get --name some-bucket
    assert_failure
    assert_stderr -p "object storage credentials not found"
}

@test "object-storage bucket public-access-block put: --json-properties-example prints example JSON" {
    run ionosctl object-storage bucket public-access-block put --json-properties-example
    assert_success
    assert_output -p "BlockPublicAcls"
    assert_output -p "RestrictPublicBuckets"
}

@test "object-storage bucket public-access-block put: nonexistent file returns error" {
    run ionosctl object-storage bucket public-access-block put --name "$TEST_BUCKET_NAME" --json-properties "/tmp/nonexistent-pab.json"
    assert_failure
}

# --- lifecycle: put, get, delete ---

@test "object-storage bucket public-access-block put: apply configuration" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"BlockPublicAcls":true,"IgnorePublicAcls":true,"BlockPublicPolicy":true,"RestrictPublicBuckets":true}
EOF
    run ionosctl object-storage bucket public-access-block put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile"
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage bucket public-access-block get: retrieve configuration" {
    run ionosctl object-storage bucket public-access-block get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "true"
}

@test "object-storage bucket public-access-block get: json output" {
    run ionosctl object-storage bucket public-access-block get --name "$TEST_BUCKET_NAME" -o json
    assert_success
    echo "$output" | jq -e '.BlockPublicAcls' >/dev/null
    echo "$output" | jq -e '.RestrictPublicBuckets' >/dev/null
}

@test "object-storage bucket public-access-block delete: remove configuration" {
    run ionosctl object-storage bucket public-access-block delete --name "$TEST_BUCKET_NAME" -f
    assert_success
    assert_output -p "deleted successfully"
}

@test "object-storage bucket public-access-block get: after delete returns error" {
    run ionosctl object-storage bucket public-access-block get --name "$TEST_BUCKET_NAME"
    assert_failure
}
