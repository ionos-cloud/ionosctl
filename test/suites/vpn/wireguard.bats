#!/usr/bin/env bats

# paths: commands/vpn/wireguard/*

load '../setup.bats'

location="de/txl"

setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    mkdir -p /tmp/bats_test
}

@test "Create datacenter, LAN, and IP block" {
    datacenter_id=$(ionosctl datacenter create --name "CLI-Test-$(randStr 8)" --location "${location}" -w -o json | jq -r '.id')
    [ -n "$datacenter_id" ] || fail "Failed to create datacenter"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id

    lan_id=$(ionosctl lan create --datacenter-id "$datacenter_id" --public=false -w -o json | jq -r '.id')
    [ -n "$lan_id" ] || fail "Failed to create LAN"
    echo "$lan_id" > /tmp/bats_test/lan_id

    ipblock_id=$(ionosctl ipblock create --location "$location" --size 1 -w -o json | jq -r '.id')
    [ -n "$ipblock_id" ] || fail "Failed to create IP block"
    echo "$ipblock_id" > /tmp/bats_test/ipblock_id

    ipblock_ip=$(ionosctl ipblock get --ipblock-id "$ipblock_id" -o json | jq -r '.properties.ips[0]')
    [ -n "$ipblock_ip" ] || fail "Failed to retrieve IP block IP"
    echo "$ipblock_ip" > /tmp/bats_test/ipblock_ip
}

@test "Create Wireguard Gateway" {
    # generate private key
    run ionosctl vpn wireguard gateway create --location "${location}" --name "cli-test-$(randStr 6)" \
      --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --lan-id 1 --connection-ip 10.7.222.239/24 \
      --gateway-ip "$(cat /tmp/bats_test/ipblock_ip)" --interface-ip 10.7.222.97/24 --private-key "$(openssl rand -base64 32)" \
      -o json
    assert_success

    gateway_id=$(echo "$output" | jq -r '.id')
    assert_regex "$gateway_id" "$uuid_v4_regex"
    echo "$gateway_id" > /tmp/bats_test/gateway_id

    # check if the gateway is created
    run ionosctl vpn wireguard gateway get --location "${location}" --gateway-id "$gateway_id" -o json
    assert_success
    assert_equal "$gateway_id" "$(echo "$output" | jq -r '.id')"

    # check if the gateway is in the list
    run ionosctl vpn wireguard gateway list --location "${location}" -o json -M 1
    assert_success
    assert_equal "$(echo "$output" | jq -r '.items | length')" "1"
}

@test "Update name of Wireguard Gateway with cols flag" {
    new_name="cli-test-$(randStr 6)"

    run ionosctl vpn wireguard gateway update --location "${location}" --gateway-id "$(cat /tmp/bats_test/gateway_id)" \
       --private-key "$(openssl rand -base64 32)" --name "$new_name" --cols ID --no-headers
    assert_success
    assert_output "$(cat /tmp/bats_test/gateway_id)"

    run ionosctl vpn wireguard gateway get --location "${location}" --gateway-id "$(cat /tmp/bats_test/gateway_id)" --cols name --no-headers
    assert_success
    assert_output "$new_name"

    # Not using no-headers shows the header and the value
    run ionosctl vpn wireguard gateway get --location "${location}" --gateway-id "$(cat /tmp/bats_test/gateway_id)" --cols name
    assert_success
    assert_output -p "Name"
    assert_output -p "$new_name"
}

@test "Can use --private-key-path flag for updating private key" {
    new_key=$(openssl rand -base64 32)
    echo "$new_key" > /tmp/bats_test/new_key

    run ionosctl vpn wireguard gateway update --location "${location}" --gateway-id "$(cat /tmp/bats_test/gateway_id)" \
      --private-key-path /tmp/bats_test/new_key -f
    assert_success

    rm /tmp/bats_test/new_key
}

@test "Create Wireguard Peer" {
    run ionosctl vpn wireguard peer create --location "${location}" --name "cli-test-$(randStr 6)" \
      --gateway-id "$(cat /tmp/bats_test/gateway_id)" --public-key "$(openssl rand -base64 32)" \
      --ips "::/0" --host "$(cat /tmp/bats_test/ipblock_ip)" -o json
    assert_success

    peer_id=$(echo "$output" | jq -r '.id')
    assert_regex "$peer_id" "$uuid_v4_regex"
    echo "$peer_id" > /tmp/bats_test/peer_id

    run ionosctl vpn wireguard peer get --location "${location}" --gateway-id "$(cat /tmp/bats_test/gateway_id)" \
      --peer-id "$peer_id" -o json
    assert_success
    assert_equal "$peer_id" "$(echo "$output" | jq -r '.id')"

    run ionosctl vpn wireguard peer list --location "${location}" \
      --gateway-id "$(cat /tmp/bats_test/gateway_id)" -o json -M 1
    assert_success
    assert_equal "$(echo "$output" | jq -r '.items | length')" "1"
}

@test "Update name of Wireguard Peer with cols flag" {
    new_name="cli-test-$(randStr 6)"

    run ionosctl vpn wireguard peer update --location "${location}" --gateway-id "$(cat /tmp/bats_test/gateway_id)" \
      --peer-id "$(cat /tmp/bats_test/peer_id)" --name "$new_name" --cols ID --no-headers
    assert_success
    assert_output "$(cat /tmp/bats_test/peer_id)"

    run ionosctl vpn wireguard peer get --location "${location}" --gateway-id "$(cat /tmp/bats_test/gateway_id)" \
      --peer-id "$(cat /tmp/bats_test/peer_id)" --cols name --no-headers
    assert_success
    assert_output "$new_name"

    # Not using no-headers shows the header and the value
    run ionosctl vpn wireguard peer get --location "${location}" \
      --gateway-id "$(cat /tmp/bats_test/gateway_id)" --peer-id "$(cat /tmp/bats_test/peer_id)" --cols name
    assert_success
    assert_output -p "Name"
    assert_output -p "$new_name"
}

@test "Delete Wireguard Peer" {
    run ionosctl vpn wireguard peer delete --location "${location}" \
      --gateway-id "$(cat /tmp/bats_test/gateway_id)" --peer-id "$(cat /tmp/bats_test/peer_id)" -f
    assert_success

    run ionosctl vpn wireguard peer list --location "${location}" \
      --gateway-id "$(cat /tmp/bats_test/gateway_id)" -o json -M 1
    assert_success
    assert_equal "$(echo "$output" | jq -r '.items | length')" "0"
}

teardown_file() {
    ionosctl vpn wireguard gateway delete -af -w --location "${location}"
    sleep 10

    ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f -w
    ionosctl ipblock delete --ipblock-id "$(cat /tmp/bats_test/ipblock_id)" -f -w

    ionosctl token delete --token "$IONOS_TOKEN"

    rm -rf /tmp/bats_test
}
