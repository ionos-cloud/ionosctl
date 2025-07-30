#!/usr/bin/env bats

# tags: monitoring, pipeline
# file names only or directories

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

@test "Create temporary sub-user with Monitoring permissions" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "test-user-$(randStr 4)" --last-name "test-last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl group create --name "test-group-$(randStr 4)" \
        -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    run ionosctl group user add --user-id "$(cat /tmp/bats_test/user_id)" \
        --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success

    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

@test "Create Monitoring Pipeline" {
    pipeline_name="cli-test-pipeline-$(randStr 1)"
    run ionosctl monitoring pipeline create --name "$pipeline_name" -o json 2> /dev/null
    assert_success

    pipeline_id=$(echo "$output" | jq -r '.id')

    assert_output -p "\"name\": \"$pipeline_name\""
    assert_output -p "\"status\": \"PROVISIONING\""

    echo "created monitoring pipeline $pipeline_id ($pipeline_name)"
    echo "$pipeline_id" > /tmp/bats_test/pipeline_id
    echo "$pipeline_name" > /tmp/bats_test/pipeline_name
}

@test "List Monitoring Pipelines" {
    pipeline_name=$(cat /tmp/bats_test/pipeline_name)
    # List Pipeline (JSON output)
    run ionosctl monitoring pipeline list -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$pipeline_name\""

    # List Pipeline (Column output)
    run ionosctl monitoring pipeline list --cols name --no-headers
    assert_success
    assert_output -p "$pipeline_name"
}

@test "Get Monitoring Pipelines" {
    pipeline_id=$(cat /tmp/bats_test/pipeline_id)
    pipeline_name=$(cat /tmp/bats_test/pipeline_name)

    # Get Pipeline (JSON output)
    run ionosctl monitoring pipeline get --pipeline-id "$pipeline_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$pipeline_name\""
    assert_output -p "\"status\": \"PROVISIONING\""
}

@test "Create Monitoring Key" {
    sleep 100
    pipeline_id=$(cat /tmp/bats_test/pipeline_id)
    run ionosctl monitoring key create --pipeline-id "$pipeline_id" -f
    assert_success
    assert_output -p "key"
}

@test "Delete Monitoring Pipeline" {
    pipeline_id=$(cat /tmp/bats_test/pipeline_id)
    run ionosctl monitoring pipeline delete --pipeline-id "$pipeline_id" -f
    assert_success
}

teardown_file() {
    (
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        ionosctl monitoring pipeline delete -af
    )

    ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
