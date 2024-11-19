#!/usr/bin/env bats

# tags: vpn, wireguard, peer

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/txl"

setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    mkdir -p /tmp/bats_test

    # if openssl doesn't exist just alias it to a random 32 b64 string "w7X5hj3jgj+DhS5pZmR7K2wU1zL3Oluy6lDfW5JcHtU="
    if ! command -v openssl &> /dev/null
    then
        alias openssl="echo w7X5hj3jgj+DhS5pZmR7K2wU1zL3Oluy6lDfW5JcHtU="
    fi
}

#@test "Ensure ipblock, and connection datacenter & lan" {
#      datacenter_id=$(find_or_create_resource \
#          "ionosctl datacenter list -M 1 -F location=${location},state=available -o json 2> /dev/null | jq -r '.items[] | .id'" \
#          "ionosctl datacenter create --name \"CLI-Test-$(randStr 8)\" --location ${location} -o json 2> /dev/null | jq -r '.id'")
#      [ -n "$datacenter_id" ] || [ "$datacenter_id" = "null" ] || fail "datacenter_id is empty"
#      assert_regex "$datacenter_id" "$uuid_v4_regex"
#      echo "$datacenter_id" > /tmp/bats_test/datacenter_id
#
#      lan_id=$(find_or_create_resource \
#          "ionosctl lan list -M 1 --datacenter-id ${datacenter_id} -F public=false-o json 2> /dev/null | jq -r '.items[] | .id'" \
#          "sleep 30 && ionosctl lan create --datacenter-id ${datacenter_id} --public=false -o json 2> /dev/null | jq -r '.id'")
#      [ -n "$lan_id" ] || [ "$lan_id" = "null" ] || fail "lan_id is empty"
#      echo "$lan_id" > /tmp/bats_test/lan_id
#
#      ipblock_id=$(find_or_create_resource \
#          "ionosctl ipblock list -M 1 -F location=${location},size=1 -o json 2> /dev/null | jq -r '.items[] | .id'" \
#          "ionosctl ipblock create --location ${location} --size 1 -o json 2> /dev/null | jq -r '.id'")
#      # same chefck above but also if null string
#      [ -n "$ipblock_id" ] || [ "$ipblock_id" = "null" ] || fail "ipblock_id is empty"
#
#      echo "$ipblock_id" > /tmp/bats_test/ipblock_id
#
#      sleep 30
#
#      ipblock_ip=$(ionosctl ipblock get --ipblock-id "$ipblock_id" -o json 2> /dev/null | jq -r '.properties.ips[0]')
#      [ -n "$ipblock_ip" ] || [ "$ipblock_ip" = "null" ] || fail "ipblock_ip is empty"
#      echo "$ipblock_ip" > /tmp/bats_test/ipblock_ip
#}

@test "Ensure ipblock, and connection datacenter & lan" {
    # place holder
    echo "2ecec754-f7f9-418a-8dfc-2824e10ac0ac" > /tmp/bats_test/datacenter_id
    echo "2" > /tmp/bats_test/lan_id
    echo "cb622374-217f-49f8-98bf-831e8c9d2e8a" > /tmp/bats_test/ipblock_id
    echo "85.215.151.14" > /tmp/bats_test/ipblock_ip
}

@test "Create Wireguard Gateway" {
    # generate private key
    run ionosctl vpn wireguard gateway create --name "cli-test-$(randStr 6)" \
      --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --lan-id 1 --connection-ip 10.7.222.239/24 \
      --gateway-ip 87.106.32.247 --interface-ip 10.7.222.97/24 --private-key "$(openssl rand -base64 32)" \
      -o json 2> /dev/null
    assert_success

    gateway_id=$(echo "$output" | jq -r '.id')
    assert_regex "$gateway_id" "$uuid_v4_regex"
    echo "$gateway_id" > /tmp/bats_test/gateway_id

    # check if the gateway is created
    run ionosctl vpn wireguard gateway get --gateway-id "$gateway_id" -o json 2> /dev/null
    assert_success
    assert_equal "$gateway_id" "$(echo "$output" | jq -r '.id')"

    # check if the gateway is in the list
    run ionosctl vpn wireguard gateway list -o json -M 1 2> /dev/null
    assert_success
    assert_equal "$(echo "$output" | jq -r '.items | length')" "1"
}

@test "Update name of Wireguard Gateway with cols flag" {
    new_name="cli-test-$(randStr 6)"

    run ionosctl vpn wireguard gateway update --gateway-id "$(cat /tmp/bats_test/gateway_id)" --name "$new_name" --cols ID --no-headers 2> /dev/null
    assert_success
    assert_equal "$(echo "$output" | jq -r '.id')" "$(cat /tmp/bats_test/gateway_id)"

    run ionosctl vpn wireguard gateway get --gateway-id "$(cat /tmp/bats_test/gateway_id)" --cols name --no-headers 2> /dev/null
    assert_success
    assert_output "$new_name"

    # Not using no-headers shows the header and the value
    run ionosctl vpn wireguard gateway get --gateway-id "$(cat /tmp/bats_test/gateway_id)" --cols name 2> /dev/null
    assert_success
    assert_output "Name"
    assert_output "$new_name"
}

@test "Create Wireguard Peer" {
    run ionosctl vpn wireguard peer create --name "cli-test-$(randStr 6)" \
      --gateway-id "$(cat /tmp/bats_test/gateway_id)" --public-key "$(openssl rand -base64 32)" \
      --ips "::/0" --host "$(cat /tmp/bats_test/ipblock_ip)" -o json 2> /dev/null
    assert_success

    peer_id=$(echo "$output" | jq -r '.id')
    assert_regex "$peer_id" "$uuid_v4_regex"
    echo "$peer_id" > /tmp/bats_test/peer_id

    run ionosctl vpn wireguard peer get --gateway-id "$(cat /tmp/bats_test/gateway_id)" \
      --peer-id "$peer_id" -o json 2> /dev/null
    assert_success
    assert_equal "$peer_id" "$(echo "$output" | jq -r '.id')"

    run ionosctl vpn wireguard peer list --gateway-id "$(cat /tmp/bats_test/gateway_id)" -o json -M 1 2> /dev/null
    assert_success
    assert_equal "$(echo "$output" | jq -r '.items | length')" "1"
}

@test "Update name of Wireguard Peer with cols flag" {
    new_name="cli-test-$(randStr 6)"

    run ionosctl vpn wireguard peer update --gateway-id "$(cat /tmp/bats_test/gateway_id)" --peer-id "$(cat /tmp/bats_test/peer_id)" --name "$new_name" --cols ID --no-headers 2> /dev/null
    assert_success
    assert_equal "$(echo "$output" | jq -r '.id')" "$(cat /tmp/bats_test/peer_id)"

    run ionosctl vpn wireguard peer get --gateway-id "$(cat /tmp/bats_test/gateway_id)" --peer-id "$(cat /tmp/bats_test/peer_id)" --cols name --no-headers 2> /dev/null
    assert_success
    assert_output "$new_name"

    # Not using no-headers shows the header and the value
    run ionosctl vpn wireguard peer get --gateway-id "$(cat /tmp/bats_test/gateway_id)" --peer-id "$(cat /tmp/bats_test/peer_id)" --cols name 2> /dev/null
    assert_success
    assert_output "Name"
    assert_output "$new_name"
}

@test "Delete Wireguard Peer" {
    run ionosctl vpn wireguard peer delete --gateway-id "$(cat /tmp/bats_test/gateway_id)" --peer-id "$(cat /tmp/bats_test/peer_id)" -f
    assert_success

    run ionosctl vpn wireguard peer list --gateway-id "$(cat /tmp/bats_test/gateway_id)" -o json -M 1 2> /dev/null
    assert_success
    assert_equal "$(echo "$output" | jq -r '.items | length')" "0"
}

teardown_file() {
    ionosctl vpn wireguard gateway delete -af

    rm -rf /tmp/bats_test
}
