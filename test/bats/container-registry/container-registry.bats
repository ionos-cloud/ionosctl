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

@test "Create Container-registry Registry" {
    registry_name="cli-test-registry$(randStr 8)"
    run ionosctl apigateway gateway create --name "$gateway_name" -o json 2> /dev/null
    assert_success

    gateway_id=$(echo "$output" | jq -r '.id')

    assert_output -p "\"name\": \"$gateway_name\""
    assert_output -p "\"status\": \"PROVISIONING\""

    sleep 20

    echo "created apigateway gateway $gateway_id ($gateway_name)"
    echo "$gateway_id" > /tmp/bats_test/gateway_id
    echo "$gateway_name" > /tmp/bats_test/gateway_name
}