#!/usr/bin/env bats

# tags: completer, completion, datacenter, server, k8s, location, lan

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

@test "Create temporary user with relevant permissions" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl compute user create --first-name "random-$(randStr 4)" --last-name "last-$(randStr 4)" \
     --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl compute group create --name "test-compl-$(randStr 4)" \
     --create-dc --create-nic --reserve-ip \
     -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    sleep 10

    run ionosctl compute group user add --user-id "$(cat /tmp/bats_test/user_id)" \
     --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success

    (
        unset IONOS_TOKEN
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        run ionosctl token generate --ttl 1h
        assert_success
        echo "$output" > /tmp/bats_test/token
    )
}

@test "Is temp user" {
    run ionosctl whoami
    assert_success
    assert_output "$(cat /tmp/bats_test/email)"
}

@test "Create Datacenter, Server, LAN for completion tests" {
    run ionosctl compute datacenter create --name "compl-test-$(randStr 8)" --location "es/vit" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5

    run ionosctl compute server create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-compl-$(randStr 8)" \
     --cores 1 --ram 1GB -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/server_id

    run ionosctl compute lan create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-compl-lan-$(randStr 8)" \
     --public -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/lan_id
}

@test "Location completion: suggests de/txl and es/vit" {
    run ionosctl __complete compute location list --location ""
    assert_success
    assert_output -p "de/txl"
    assert_output -p "es/vit"
}

@test "Datacenter ID completion: suggests created datacenter" {
    run ionosctl __complete compute server list --datacenter-id ""
    assert_success
    assert_output -p "$(cat /tmp/bats_test/datacenter_id)"
}

@test "Server ID completion: suggests created server" {
    run ionosctl __complete compute server get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id ""
    assert_success
    assert_output -p "$(cat /tmp/bats_test/server_id)"
}

@test "LAN ID completion: suggests created LAN" {
    run ionosctl __complete compute nic create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --lan-id ""
    assert_success
    assert_output -p "$(cat /tmp/bats_test/lan_id)"
}

@test "Image ID completion: returns non-empty" {
    run ionosctl __complete compute volume create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --image-id ""
    assert_success
    [ -n "$output" ]
}

@test "K8s version completion: returns version numbers" {
    run ionosctl __complete compute k8s cluster create --k8s-version ""
    assert_success
    [ -n "$output" ]
}

@test "Completion via hidden alias: server list --datacenter-id" {
    run ionosctl __complete server list --datacenter-id ""
    assert_success
    assert_output -p "$(cat /tmp/bats_test/datacenter_id)"
}

@test "Delete resources" {
    run ionosctl compute lan delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" -w -f -t 600
    assert_success

    run ionosctl compute server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" -w -t 600 -f
    assert_success

    run ionosctl compute datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f -w -t 600
    assert_success
}

teardown_file() {
    (
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
        ionosctl compute datacenter delete -af
    )

    ionosctl compute user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl compute group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
