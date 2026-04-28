#!/usr/bin/env bats

# tags: firewallrule

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

    run ionosctl compute group create --name "test-fw-$(randStr 4)" \
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
    run ionosctl compute datacenter create --name "fw-test-$(randStr 8)" --location "es/vit" -w -t 600 -o json 2> /dev/null
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

@test "Create LAN" {
    run ionosctl compute lan create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-fw-lan-$(randStr 8)" \
     --public -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/lan_id
}

@test "Create NIC" {
    run ionosctl compute nic create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" --name "bats-fw-nic-$(randStr 8)" --firewall-active -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/nic_id
    sleep 5
}

@test "Create firewall rule (TCP port 22)" {
    run ionosctl compute firewallrule create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" \
     --name "allow-ssh-$(randStr 4)" --protocol TCP --port-range-start 22 --port-range-end 22 \
     --source-ip "0.0.0.0" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/fw_rule_id
}

@test "Firewall rule get --cols" {
    run ionosctl compute firewallrule get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" \
     --firewallrule-id "$(cat /tmp/bats_test/fw_rule_id)" --cols FirewallRuleId,Name,Protocol --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/fw_rule_id)"
    assert_output -p "TCP"
}

@test "Firewall rule list" {
    run ionosctl compute firewallrule list --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/fw_rule_id)"
}

@test "Create second firewall rule (ICMP)" {
    run ionosctl compute firewallrule create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" \
     --name "allow-ping-$(randStr 4)" --protocol ICMP --icmp-type 8 --icmp-code 0 \
     -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/fw_rule_id_2
}

@test "Update firewall rule (change port to 443)" {
    run ionosctl compute firewallrule update --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" \
     --firewallrule-id "$(cat /tmp/bats_test/fw_rule_id)" \
     --port-range-start 443 --port-range-end 443 -w -t 600 -o json 2> /dev/null
    assert_success

    run ionosctl compute firewallrule get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" \
     --firewallrule-id "$(cat /tmp/bats_test/fw_rule_id)" --cols PortRangeStart --no-headers
    assert_success
    assert_output "443"
}

@test "Delete firewall rules" {
    run ionosctl compute firewallrule delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" \
     --firewallrule-id "$(cat /tmp/bats_test/fw_rule_id)" -f -w -t 600
    assert_success

    run ionosctl compute firewallrule delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" \
     --firewallrule-id "$(cat /tmp/bats_test/fw_rule_id_2)" -f -w -t 600
    assert_success
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
        ionosctl compute datacenter delete -af
    )

    ionosctl compute user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl compute group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
