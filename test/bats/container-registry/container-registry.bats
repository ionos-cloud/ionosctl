#!/usr/bin/env bats

# tags: container-registry, cr, contreg, cont-reg,

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/fra"

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

@test "Create temporary sub-user with ApiGateway permissions" {
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

@test "Create Container-registry Registry" {
    registry_name="registry-cli-test"
    run ionosctl container-registry registry create --name "$registry_name" --location "$location" -o json 2> /dev/null
    assert_success

    registry_id=$(echo "$output" | jq -r '.id')

    assert_output -p "\"name\": \"$registry_name\""
    assert_output -p "\"location\": \"$location\""

    echo "$registry_id" > /tmp/bats_test/registry_id
}

@test "List Container-registry registry" {
    run ionosctl container-registry registry list -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"registry-cli-test\""
}

@test "Get Container-registry registry" {
    registry_id=$(cat /tmp/bats_test/registry_id)

    run ionosctl container-registry registry get --registry-id "$registry_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"registry-cli-test\""
    assert_output -p "\"location\": \"$location\""
}

@test "Update Container-registry registry" {
    registry_id=$(cat /tmp/bats_test/registry_id)

    run ionosctl container-registry registry update --registry-id "$registry_id" --garbage-collection-schedule-days Friday --garbage-collection-schedule-time 01:23:00+00:00 -o json 2> /dev/null
    assert_success
    assert_output -p "\"time\": \"01:23:00+00:00\""
}

@test "Create Container-registry token" {
    registry_id=$(cat /tmp/bats_test/registry_id)

    run ionosctl container-registry token create --registry-id "$registry_id" --name registry-token-test -o json 2> /dev/null
    assert_success

    token_id=$(echo "$output" | jq -r '.id')

    assert_output -p "\"name\": \"registry-token-test\""
    echo "$token_id" > /tmp/bats_test/token_id
}

@test "List Container-registry token" {
    registry_id=$(cat /tmp/bats_test/registry_id)

    run ionosctl container-registry token list --registry-id "$registry_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"registry-token-test\""
}

@test "Get Container-registry token" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    token_id=$(cat /tmp/bats_test/token_id)

    run ionosctl container-registry token get --token-id "$token_id" --registry-id "$registry_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"registry-token-test\""
}

@test "Update Container-registry token" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    token_id=$(cat /tmp/bats_test/token_id)

    run ionosctl container-registry token update --registry-id "$registry_id" --token-id "$token_id" --status disabled -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"registry-token-test\""
    assert_output -p "\"status\": \"disabled\""
}

@test "Delete Container-registry token" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    token_id=$(cat /tmp/bats_test/token_id)

    run ionosctl container-registry token delete --registry-id "$registry_id" --token-id "$token_id" -f
    assert_success
}

@test "Delete Container-registry registry" {
    registry_id=$(cat /tmp/bats_test/registry_id)

    run ionosctl container-registry registry delete --registry-id "$registry_id" -f
    assert_success
}

teardown_file() {
    (
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        ionosctl apigateway gateway delete -af
    )

    ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}