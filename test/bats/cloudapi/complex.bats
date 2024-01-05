#!/usr/bin/env bats

load 'libs/bats-support/load'
load 'libs/bats-assert/load'

setup() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

randStr() {
    size=$1
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c "$size"
}

@test "Create a server under a NAT" {
    run ionosctl ipblock create --location de/txl -o json -w
    [ "$status" -eq 0 ]
    ipblock_id=$(echo "$output" | jq -r '.id')
    [[ $ipblock_id =~ $uuid_v4_regex ]]
    echo "created ipblock $ipblock_id"

    run ionosctl datacenter create --name "$(randStr 8)" --location de/txl -o json -w
    [ "$status" -eq 0 ]
    datacenter_id=$(echo "$output" | jq -r '.id')
    [[ $datacenter_id =~ $uuid_v4_regex ]]
    echo "created datacenter $datacenter_id"

    # Create a barebones server
    run ionosctl server create --datacenter-id "$datacenter_id" --cpu-family INTEL_SKYLAKE --cores 1 --ram 1GB -o json -w
    [ "$status" -eq 0 ]
    server_id=$(echo "$output" | jq -r '.id')
    [[ $server_id =~ $uuid_v4_regex ]]
    echo "created server $server_id"

    # Try getting the already created ipblock ID via a filtered list this time
    run ionosctl ipblock list -F location=de/txl -M 1 -o json
    [ "$status" -eq 0 ]
    found_ip=$(echo "$output" | jq -r '.items[] | .properties.ips[0]')
    echo "found ip $found_ip"

    # Create a NAT
    run ionosctl natgateway create --datacenter-id "$datacenter_id" --ips "$found_ip" -o json -w --timeout 500
    [ "$status" -eq 0 ]
    nat_id=$(echo "$output" | jq -r '.id')
    [[ $nat_id =~ $uuid_v4_regex ]]
    echo "created nat $nat_id"

    # Create a LAN. Use --cols to extract its ID.
    run ionosctl lan create --datacenter-id "$datacenter_id" --cols LanId --no-headers
    [ "$status" -eq 0 ]
    lan_id="$output"
    # Test that Lan ID is a numeric value
    [[ $lan_id =~ ^[0-9]+$ ]]
    echo "created lan $lan_id"

    sleep 60

    # Add the LAN to the NAT
    run ionosctl natgateway lan add --datacenter-id "$datacenter_id" --natgateway-id "$nat_id" --lan-id "$lan_id"
    [ "$status" -eq 0 ]
    echo "added lan $lan_id to nat $nat_id"

    # Volume Password
    vol_pass=$(randStr 16)
    # Create a volume
    run ionosctl volume create --datacenter-id "$datacenter_id" --size 50 --image-alias ubuntu:latest -p "$vol_pass" -o json -w
    [ "$status" -eq 0 ]
    vol_id=$(echo "$output" | jq -r '.id')
    [[ $vol_id =~ $uuid_v4_regex ]]
    echo "created volume $vol_id"

    # Attach volume to server
    run ionosctl volume attach --server-id "$server_id" --volume-id "$vol_id" --datacenter-id "$datacenter_id" -w
    [ "$status" -eq 0 ]
}

teardown() {
    echo "cleaning up..."
    if [ -n "$datacenter_id" ]; then
        run ionosctl datacenter delete -f --datacenter-id $datacenter_id
        assert_success
    fi

    if [ -n "$ipblock_id" ]; then
        run ionosctl ipblock delete -f --ipblock-id $ipblock_id
        assert_success
    fi
}
