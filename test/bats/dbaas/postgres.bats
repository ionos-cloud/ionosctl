#!/usr/bin/env bats

# tags: postgres, pg

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

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

@test "Create temporary sub-user with postgres permissions" {
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

@test "Get Postgres Api version" {
    gateway_name="cli-test-gateway-$(randStr 8)"
    run ionosctl apigateway gateway create --name "$gateway_name" -o json 2> /dev/null
    assert_success

    gateway_id=$(echo "$output" | jq -r '.id')

    assert_output -p "\"name\": \"$gateway_name\""
    assert_output -p "\"status\": \"PROVISIONING\""

    echo "created apigateway gateway $gateway_id ($gateway_name)"
    echo "$gateway_id" > /tmp/bats_test/gateway_id
    echo "$gateway_name" > /tmp/bats_test/gateway_name
}
