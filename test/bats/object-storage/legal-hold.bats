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
    export TEST_BUCKET_NAME="ionosctl-ci-lh-$(randStr 8 | tr '[:upper:]' '[:lower:]')"
    export TEST_KEY="test-legal-hold.txt"
    export TEST_FILE="/tmp/ionosctl-test-lh-$$.txt"

    echo "legal hold test content" > "$TEST_FILE"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION" --object-lock
    assert_success

    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --source "$TEST_FILE" 2>/dev/null
    assert_success

    echo "created test bucket with object-lock for legal-hold tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    rm -f "$TEST_FILE"
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        # SAFETY: Always remove legal hold before deleting
        run ionosctl object-storage object legal-hold put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --status OFF 2>/dev/null
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- legal-hold get: validation ---

@test "object-storage object legal-hold get: missing --name flag returns error" {
    run ionosctl object-storage object legal-hold get --key foo 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object legal-hold get: missing --key flag returns error" {
    run ionosctl object-storage object legal-hold get --name "$TEST_BUCKET_NAME" 2>&1
    assert_failure
    assert_output -p "requires at least"
}

# --- legal-hold put: validation ---

@test "object-storage object legal-hold put: missing --name flag returns error" {
    run ionosctl object-storage object legal-hold put --key foo --status ON 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object legal-hold put: missing --key flag returns error" {
    run ionosctl object-storage object legal-hold put --name "$TEST_BUCKET_NAME" --status ON 2>&1
    assert_failure
    assert_output -p "requires at least"
}

@test "object-storage object legal-hold put: missing --status flag returns error" {
    run ionosctl object-storage object legal-hold put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" 2>&1
    assert_failure
    assert_output -p "requires at least"
}

# --- legal-hold put/get: lifecycle ---

# SAFETY: Legal hold ON is applied and then immediately set to OFF in subsequent tests.
# Teardown also ensures legal hold is always removed before cleanup.
@test "object-storage object legal-hold put: set status ON" {
    run ionosctl object-storage object legal-hold put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --status ON 2>/dev/null
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage object legal-hold get: verify status is ON" {
    run ionosctl object-storage object legal-hold get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" 2>/dev/null
    assert_success
    assert_output -p "ON"
}

@test "object-storage object legal-hold put: set status OFF" {
    run ionosctl object-storage object legal-hold put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --status OFF 2>/dev/null
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage object legal-hold get: verify status is OFF" {
    run ionosctl object-storage object legal-hold get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" 2>/dev/null
    assert_success
    assert_output -p "OFF"
}
