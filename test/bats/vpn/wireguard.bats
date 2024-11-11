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
}

@test "Create connection datacenter & lan" {
      datacenter_id=$(find_or_create_resource \
          "ionosctl datacenter list -M 1 -F location=${location},state=available -o json 2> /dev/null | jq -r '.items[] | .id'" \
          "ionosctl datacenter create --name \"CLI-Test-$(randStr 8)\" --location ${location} -o json 2> /dev/null | jq -r '.id'")
      [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
      assert_regex "$datacenter_id" "$uuid_v4_regex"

      lan_id=$(find_or_create_resource \
          "ionosctl lan list -M 1 --datacenter-id ${datacenter_id} -F public=false-o json 2> /dev/null | jq -r '.items[] | .id'" \
          "sleep 30 && ionosctl lan create --datacenter-id ${datacenter_id} --public=false -o json 2> /dev/null | jq -r '.id'")
      [ -n "$lan_id" ] || fail "lan_id is empty"

      sleep 60
}

@test "Create Wireguard Gateway" {
    # generate private key
    key=$(openssl rand -base64 32)

    run ionosctl vpn wireguard gateway create --name "cli-test-$(randStr 6)" \
      --datacenter-id 2ecec754-f7f9-418a-8dfc-2824e10ac0ac --lan-id 1 --connection-ip 10.7.222.239/24 \
      --gateway-ip 87.106.32.247 --interface-ip 10.7.222.97/24 --private-key "$key" \
      -o json 2> /dev/null
    assert_success

    id=$(echo "$output" | jq -r '.id')
    assert_regex "id" "$uuid_v4_regex"
}
