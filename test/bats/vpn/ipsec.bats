#!/usr/bin/env bats

# tags: vpn, ipsec, tunnel

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/txl"

setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    mkdir -p /tmp/bats_test
}

@test "Create datacenter, LAN, and IP block" {
    datacenter_id=$(ionosctl datacenter create --name "CLI-Test-$(randStr 8)" --location "${location}" -o json 2> /dev/null | jq -r '.id')
    [ -n "$datacenter_id" ] || fail "Failed to create datacenter"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id

    sleep 60

    lan_id=$(ionosctl lan create --datacenter-id "$datacenter_id" --public=false -o json 2> /dev/null | jq -r '.id')
    [ -n "$lan_id" ] || fail "Failed to create LAN"
    echo "$lan_id" > /tmp/bats_test/lan_id

    ipblock_id=$(ionosctl ipblock create --location "$location" --size 1 -o json 2> /dev/null | jq -r '.id')
    [ -n "$ipblock_id" ] || fail "Failed to create IP block"
    echo "$ipblock_id" > /tmp/bats_test/ipblock_id

    ipblock_ip=$(ionosctl ipblock get --ipblock-id "$ipblock_id" -o json 2> /dev/null | jq -r '.properties.ips[0]')
    [ -n "$ipblock_ip" ] || fail "Failed to retrieve IP block IP"
    echo "$ipblock_ip" > /tmp/bats_test/ipblock_ip

    lan_status=""
    i=0
    while [ "$lan_status" != "AVAILABLE" ] && [ $i -lt 30 ]; do
        lan_status=$(ionosctl lan get --lan-id "$lan_id" --datacenter-id "$datacenter_id" -o json 2> /dev/null | jq -r '.metadata.state')
        sleep 10
        i=$((i+1))
    done
    [ "$lan_status" = "AVAILABLE" ] || fail "LAN is not available"

}

@test "Create IPSec Gateway" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)
    ipblock_ip=$(cat /tmp/bats_test/ipblock_ip)

    run ionosctl vpn ipsec gateway create --name "cli-test-$(randStr 6)" \
      --datacenter-id "$datacenter_id" --lan-id "$lan_id" --connection-ip "10.7.222.239/24" --gateway-ip "$ipblock_ip" \
      -o json 2> /dev/null
    assert_success

    gateway_id=$(echo "$output" | jq -r '.id')
    [ -n "$gateway_id" ] || fail "Failed to create IPSec Gateway"
    echo "$gateway_id" > /tmp/bats_test/ipsec_gateway_id

    # Verify the gateway was created
    run ionosctl vpn ipsec gateway get --gateway-id "$gateway_id" -o json 2> /dev/null
    assert_success
    assert_equal "$gateway_id" "$(echo "$output" | jq -r '.id')"
}

@test "Create IPSec Tunnel" {
    gateway_id=$(cat /tmp/bats_test/ipsec_gateway_id)

    run ionosctl vpn ipsec tunnel create --gateway-id "$gateway_id" --name "cli-test-tunnel" \
      --host "192.168.1.1" --auth-method "PSK" --psk-key "$(openssl rand -base64 32)" \
      --ike-diffie-hellman-group "19-ECP256" --ike-encryption-algorithm "AES256" --ike-integrity-algorithm "SHA256" --ike-lifetime 86400 \
      --esp-diffie-hellman-group "19-ECP256" --esp-encryption-algorithm "AES256" --esp-integrity-algorithm "SHA256" --esp-lifetime 3600 \
      --cloud-network-cidrs "10.0.0.0/16" --peer-network-cidrs "192.168.0.0/16" -o json 2> /dev/null
    assert_success

    tunnel_id=$(echo "$output" | jq -r '.id')
    [ -n "$tunnel_id" ] || fail "Failed to create IPSec Tunnel"
    echo "$tunnel_id" > /tmp/bats_test/ipsec_tunnel_id

    # Verify the tunnel was created
    run ionosctl vpn ipsec tunnel get --gateway-id "$gateway_id" --tunnel-id "$tunnel_id" -o json 2> /dev/null
    assert_success
    assert_equal "$tunnel_id" "$(echo "$output" | jq -r '.id')"
}

@test "Create IPSec Tunnel with JSON properties file" {
    gateway_id=$(cat /tmp/bats_test/ipsec_gateway_id)

    json=$(cat <<EOF
{
  "properties": {
    "name": "$name",
    "description": "Allows local subnet X to connect to virtual network Y.",
    "remoteHost": "vpn.mycompany.com",
    "auth": {
      "method": "PSK",
      "psk": {
        "key": "X2wosbaw74M8hQGbK3jCCaEusR6CCFRa"
      }
    },
    "ike": {
      "diffieHellmanGroup": "16-MODP4096",
      "encryptionAlgorithm": "AES256",
      "integrityAlgorithm": "SHA256",
      "lifetime": 86400
    },
    "esp": {
      "diffieHellmanGroup": "16-MODP4096",
      "encryptionAlgorithm": "AES256",
      "integrityAlgorithm": "SHA256",
      "lifetime": 3600
    },
    "cloudNetworkCIDRs": [
      "192.168.1.100/24"
    ],
    "peerNetworkCIDRs": [
      "1.2.3.4/32"
    ]
  }
}
EOF
    )
    printf "%s\n" "$json" > /tmp/bats_test/tunnel_properties.json

    run ionosctl vpn ipsec tunnel create --gateway-id "$gateway_id" --json-properties /tmp/bats_test/tunnel_properties.json -o json
    assert_success

    tunnel_id=$(echo "$output" | jq -r '.id')
    [ -n "$tunnel_id" ] || fail "Failed to create IPSec Tunnel with JSON properties"
    echo "$tunnel_id" > /tmp/bats_test/ipsec_json_tunnel_id
}

teardown_file() {
    ionosctl vpn ipsec gateway delete -af

    ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f
    ionosctl ipblock delete --ipblock-id "$(cat /tmp/bats_test/ipblock_id)" -f

    ionosctl token delete --token "$IONOS_TOKEN"

#    rm -rf /tmp/bats_test
}
