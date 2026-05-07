#!/usr/bin/env bats

# paths: commands/object-storage/*

load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-ol-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION" --object-lock
    assert_success

    echo "created test bucket with object-lock for object-lock tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all --bypass-governance-retention -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- object-lock get ---

@test "object-storage bucket object-lock get: missing --name flag returns error" {
    run ionosctl object-storage bucket object-lock get
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket object-lock get: returns config for object-lock enabled bucket" {
    run ionosctl object-storage bucket object-lock get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "Enabled"
}

@test "object-storage bucket object-lock get: json output contains expected fields" {
    run ionosctl object-storage bucket object-lock get --name "$TEST_BUCKET_NAME" -o json
    assert_success
    echo "$output" | jq -e '.ObjectLockEnabled' >/dev/null
}

# --- object-lock put ---

@test "object-storage bucket object-lock put: missing --name flag returns error" {
    run ionosctl object-storage bucket object-lock put --mode GOVERNANCE --days 1
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage bucket object-lock put: missing --mode flag returns error" {
    run ionosctl object-storage bucket object-lock put --name "$TEST_BUCKET_NAME" --days 1
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage bucket object-lock put: missing --days and --years returns error" {
    run ionosctl object-storage bucket object-lock put --name "$TEST_BUCKET_NAME" --mode GOVERNANCE
    assert_failure
    assert_stderr -p "at least one of"
}

# SAFETY: Uses GOVERNANCE mode only (never COMPLIANCE). Default retention only
# affects NEW objects placed after the config, not existing ones or the bucket itself.
@test "object-storage bucket object-lock put: apply GOVERNANCE config with 1 day retention" {
    run ionosctl object-storage bucket object-lock put --name "$TEST_BUCKET_NAME" --mode GOVERNANCE --days 1
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage bucket object-lock put: verify config after put" {
    run ionosctl object-storage bucket object-lock get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "GOVERNANCE"
}
