#!/usr/bin/env bats

# paths: commands/object-storage/*

load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-ret-$(randStr 8 | tr '[:upper:]' '[:lower:]')"
    export TEST_KEY="test-retention.txt"
    export TEST_FILE="/tmp/ionosctl-test-ret-$$.txt"

    # Compute retain-until-date as tomorrow in RFC 3339 format
    export RETAIN_UNTIL_DATE="$(date -u -d '+1 day' '+%Y-%m-%dT00:00:00Z')"

    echo "retention test content" > "$TEST_FILE"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION" --object-lock
    assert_success

    run ionosctl object-storage object put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" --source "$TEST_FILE"
    assert_success

    echo "created test bucket with object-lock for retention tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    rm -f "$TEST_FILE"
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        # SAFETY: Use --bypass-governance-retention to delete objects with GOVERNANCE retention
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all --bypass-governance-retention -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- retention get: validation ---

@test "object-storage object retention get: missing --name flag returns error" {
    run ionosctl object-storage object retention get --key foo
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object retention get: missing --key flag returns error" {
    run ionosctl object-storage object retention get --name "$TEST_BUCKET_NAME"
    assert_failure
    assert_stderr -p "requires at least"
}

# --- retention put: validation ---

@test "object-storage object retention put: missing --name flag returns error" {
    run ionosctl object-storage object retention put --key foo --mode GOVERNANCE --retain-until-date "$RETAIN_UNTIL_DATE"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object retention put: missing --key flag returns error" {
    run ionosctl object-storage object retention put --name "$TEST_BUCKET_NAME" --mode GOVERNANCE --retain-until-date "$RETAIN_UNTIL_DATE"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object retention put: missing --mode flag returns error" {
    run ionosctl object-storage object retention put --name "$TEST_BUCKET_NAME" --key foo --retain-until-date "$RETAIN_UNTIL_DATE"
    assert_failure
    assert_stderr -p "requires at least"
}

@test "object-storage object retention put: missing --retain-until-date flag returns error" {
    run ionosctl object-storage object retention put --name "$TEST_BUCKET_NAME" --key foo --mode GOVERNANCE
    assert_failure
    assert_stderr -p "requires at least"
}

# --- retention put/get: lifecycle ---

# SAFETY: Uses GOVERNANCE mode ONLY (never COMPLIANCE).
# Teardown uses --bypass-governance-retention to clean up.
@test "object-storage object retention put: apply GOVERNANCE retention" {
    run ionosctl object-storage object retention put --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" \
        --mode GOVERNANCE --retain-until-date "$RETAIN_UNTIL_DATE"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage object retention get: verify retention after put" {
    run ionosctl object-storage object retention get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY"
    assert_success
    assert_output -p "GOVERNANCE"
}

@test "object-storage object retention get: json output contains expected fields" {
    run ionosctl object-storage object retention get --name "$TEST_BUCKET_NAME" --key "$TEST_KEY" -o json
    assert_success
    echo "$output" | jq -e '.Mode' >/dev/null
    echo "$output" | jq -e '.RetainUntilDate' >/dev/null
}
