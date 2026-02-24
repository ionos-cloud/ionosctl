#!/usr/bin/env bats

# tags: nic, lan, ipblock, ipconsumer

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
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

    run ionosctl compute group create --name "test-net-$(randStr 4)" \
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

@test "Create Datacenter" {
    run ionosctl compute datacenter create --name "net-test-$(randStr 8)" --location "es/vit" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5
}

@test "Create Server" {
    run ionosctl compute server create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --cores 1 --ram 1GB -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/server_id
}

@test "Create public LAN" {
    run ionosctl compute lan create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --public -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/lan_id
}

@test "Reserve IPBlock" {
    run ionosctl compute ipblock create --location "es/vit" --size 1 --name "bats-test-$(randStr 8)" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.properties.ips[0]' > /tmp/bats_test/ip
    echo "$output" | jq -r '.id' > /tmp/bats_test/ipblock_id
}

@test "Create NIC with LAN and IP" {
    run ionosctl compute nic create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" --name "bats-test-$(randStr 8)" --ips "$(cat /tmp/bats_test/ip)" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/nic_id
    sleep 5
}

@test "NIC get --cols NicId,Name,LanId" {
    run ionosctl compute nic get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" --cols NicId,Name,Lan --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/nic_id)"
}

@test "Creating a NIC with a non-existent LAN ID will create a LAN" {
    run ionosctl compute nic create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" \
     --lan-id 123 -w -t 600 -o json 2> /dev/null
    assert_success
    sleep 5

    run ionosctl compute lan get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --lan-id 123 --no-headers --cols Public
    assert_success
    assert_output "false"
}

@test "LAN list, verify public=true" {
    run ionosctl compute lan list --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -o json 2> /dev/null
    assert_success
    # Verify our original LAN is public
    lan_id="$(cat /tmp/bats_test/lan_id)"
    public=$(echo "$output" | jq -r --arg id "$lan_id" '.items[] | select(.id == $id) | .properties.public')
    assert_equal "$public" "true"
}

@test "IPBlock list, verify size" {
    run ionosctl compute ipblock get --ipblock-id "$(cat /tmp/bats_test/ipblock_id)" --cols Size --no-headers
    assert_success
    assert_output "1"
}

@test "IP consumer list" {
    run ionosctl compute ipconsumer list --ipblock-id "$(cat /tmp/bats_test/ipblock_id)" --no-headers
    assert_success
    # The IP is assigned to a NIC, so there should be at least one consumer
    assert_output -p "$(cat /tmp/bats_test/ip)"
}

@test "Delete NIC" {
    run ionosctl compute nic delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" -w -f -t 600
    assert_success
}

@test "Delete LAN" {
    run ionosctl compute lan delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" -w -f -t 600
    assert_success
}

@test "Delete Server" {
    run ionosctl compute server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" -w -t 600 -f
    assert_success
}

@test "Delete Datacenter" {
    run ionosctl compute datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f -w -t 600
    assert_success
}

teardown_file() {
    (
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
        ionosctl compute ipblock delete -af
        ionosctl compute datacenter delete -af
    )

    ionosctl compute user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl compute group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
