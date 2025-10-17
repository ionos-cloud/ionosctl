#!/usr/bin/env bats

# tags: container-registry, cr

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

@test "Generate Token" {
    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

@test "Create Container Registry" {
    registry_name="cli-test-$(randStr 8 | tr '[:upper:]' '[:lower:]' | tr -cd 'a-z0-9')"
    echo "$registry_name" > /tmp/bats_test/registry_name

    run ionosctl container-registry registry create --name "$registry_name" --location "$location" -o json 2> /dev/null
    assert_success

    registry_id=$(echo "$output" | jq -r '.id')
    assert_regex "$registry_id" "$uuid_v4_regex"
    assert_output -p "\"name\": \"$registry_name\""
    assert_output -p "\"location\": \"$location\""

    echo "$registry_id" > /tmp/bats_test/registry_id
}

@test "List Container Registries" {
    registry_name=$(cat /tmp/bats_test/registry_name)

    run ionosctl container-registry registry list -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$registry_name\""
}

@test "Get Container Registry" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    registry_name=$(cat /tmp/bats_test/registry_name)

    run ionosctl container-registry registry get --registry-id "$registry_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"id\": \"$registry_id\""
    assert_output -p "\"name\": \"$registry_name\""
    assert_output -p "\"location\": \"$location\""
}

@test "Update Container Registry" {
    registry_id=$(cat /tmp/bats_test/registry_id)

    run ionosctl container-registry registry update --registry-id "$registry_id" \
        --garbage-collection-schedule-days Friday \
        --garbage-collection-schedule-time 01:23:00+00:00 -o json 2> /dev/null
    assert_success
    assert_output -p "\"time\": \"01:23:00+00:00\""
    assert_output -p "\"Friday\""
}

@test "List Container Registry Locations" {
    run ionosctl container-registry locations -o json 2> /dev/null
    assert_success
    assert_output -p "\"de/fra\""
}

@test "List Container Registry Repositories" {
    sleep 15
    registry_id=$(cat /tmp/bats_test/registry_id)

    run ionosctl container-registry repository list --registry-id "$registry_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"id\": \"repositories\"",
}

@test "Create Registry Token" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    token_name="registry-token-test-$(randStr 5)"
    echo "$token_name" > /tmp/bats_test/token_name

    run ionosctl container-registry token create --registry-id "$registry_id" --name "$token_name" -o json 2> /dev/null
    assert_success

    token_id=$(echo "$output" | jq -r '.id')
    assert_regex "$token_id" "$uuid_v4_regex"
    assert_output -p "\"name\": \"$token_name\""

    echo "$token_id" > /tmp/bats_test/token_id
}

@test "List Registry Tokens" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    token_name=$(cat /tmp/bats_test/token_name)

    run ionosctl container-registry token list --registry-id "$registry_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$token_name\""
}

@test "Get Registry Token" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    token_id=$(cat /tmp/bats_test/token_id)
    token_name=$(cat /tmp/bats_test/token_name)

    run ionosctl container-registry token get --registry-id "$registry_id" --token-id "$token_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"id\": \"$token_id\""
    assert_output -p "\"name\": \"$token_name\""
}

@test "List All registry tokens" {
    run ionosctl container-registry token list --all -o json 2> /dev/null
    assert_success
    token_id=$(cat /tmp/bats_test/token_id)
    assert_output -p "\"id\": \"$token_id\""
}

@test "Update Registry Token" {
    registry_id=$(cat /tmp/bats_test/registry_id)
    token_id=$(cat /tmp/bats_test/token_id)

    run ionosctl container-registry token update --registry-id "$registry_id" \
        --token-id "$token_id" --status disabled -o json 2> /dev/null
    assert_success
    assert_output -p "\"status\": \"disabled\""

    run ionosctl container-registry token scope add --registry-id "$registry_id" --name "repo-full-access" \
        --token-id "$token_id" --type repository --actions "*" -o json 2> /dev/null
    assert_success
    assert_output -p "\"actions\": [\n      \"*\"\n    ]"
    assert_output -p "\"name\": \"repo-full-access\""
}

#@test "Delete Registry Token" {
#    registry_id=$(cat /tmp/bats_test/registry_id)
#    token_id=$(cat /tmp/bats_test/token_id)
#
#    run ionosctl container-registry token delete --registry-id "$registry_id" --token-id "$token_id" -f 2> /dev/null
#    assert_success
#}
#
#@test "Delete Container Registry" {
#    registry_id=$(cat /tmp/bats_test/registry_id)
#
#    run ionosctl container-registry registry delete --registry-id "$registry_id" -f 2> /dev/null
#    assert_success
#}
#
#teardown_file() {
#    echo "Cleaning up token"
#    if [[ -f /tmp/bats_test/token ]]; then
#        run ionosctl token delete --token "$(cat /tmp/bats_test/token)" -f
#        unset IONOS_TOKEN
#    fi
#
#    ionosctl container-registry registry delete -af
#
#    echo "Cleaning up test directory"
#    rm -rf /tmp/bats_test
#}
