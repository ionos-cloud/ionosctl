#!/usr/bin/env bats

# tags: natgateway

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

    run ionosctl compute group create --name "test-natgw-$(randStr 4)" \
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
    run ionosctl compute datacenter create --name "natgw-test-$(randStr 8)" --location "es/vit" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5
}

@test "Reserve IPBlock for NAT Gateway" {
    run ionosctl compute ipblock create --location "es/vit" --size 1 --name "bats-natgw-$(randStr 8)" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.properties.ips[0]' > /tmp/bats_test/ip
    echo "$output" | jq -r '.id' > /tmp/bats_test/ipblock_id
}

@test "Create LAN for NAT Gateway" {
    run ionosctl compute lan create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-natgw-lan-$(randStr 8)" \
     --public=false -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/lan_id
}

@test "Create NAT Gateway" {
    run ionosctl compute natgateway create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --name "bats-natgw-$(randStr 8)" --ips "$(cat /tmp/bats_test/ip)" \
     -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/natgateway_id
}

@test "NAT Gateway get --cols" {
    run ionosctl compute natgateway get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --natgateway-id "$(cat /tmp/bats_test/natgateway_id)" --cols NatGatewayId,Name --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/natgateway_id)"
}

@test "NAT Gateway list" {
    run ionosctl compute natgateway list --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/natgateway_id)"
}

@test "Add NAT Gateway rule" {
    run ionosctl compute natgateway rule create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --natgateway-id "$(cat /tmp/bats_test/natgateway_id)" --name "bats-rule-$(randStr 8)" \
     --ip "$(cat /tmp/bats_test/ip)" --protocol TCP --source-subnet "10.0.1.0/24" --target-subnet "10.0.2.0/24" \
     --port-range-start 1 --port-range-end 65534 \
     -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/rule_id
}

@test "List NAT Gateway rules" {
    run ionosctl compute natgateway rule list --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --natgateway-id "$(cat /tmp/bats_test/natgateway_id)" --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/rule_id)"
}

@test "Get NAT Gateway rule" {
    run ionosctl compute natgateway rule get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --natgateway-id "$(cat /tmp/bats_test/natgateway_id)" --rule-id "$(cat /tmp/bats_test/rule_id)" --no-headers
    assert_success
    assert_output -p "TCP"
}

@test "Delete NAT Gateway rule" {
    run ionosctl compute natgateway rule delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --natgateway-id "$(cat /tmp/bats_test/natgateway_id)" --rule-id "$(cat /tmp/bats_test/rule_id)" -f -w -t 600
    assert_success
}

@test "Delete NAT Gateway" {
    run ionosctl compute natgateway delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --natgateway-id "$(cat /tmp/bats_test/natgateway_id)" -f -w -t 600
    assert_success
}

@test "Delete LAN" {
    run ionosctl compute lan delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" -w -f -t 600
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
