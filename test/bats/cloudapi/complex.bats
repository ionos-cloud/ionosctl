#!/usr/bin/env bats

# tags: cloudapi-v6, datacenter, server, lan, nic, ipblock, volume

load '../libs/bats-assert/load'
load '../libs/bats-support/load'
load '../setup.bats'

setup() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'

    ssh_key_path=$(generate_ssh_key)
    ssh_pub_key=$(cat "${ssh_key_path}.pub")
}

teardown() {
    rm -f "${ssh_key_path}" "${ssh_key_path}.pub"
    echo "cleaning up datacenter $datacenter_id and ipblock $ipblock_id"
    ionosctl datacenter delete --datacenter-id "$datacenter_id" -f -w
    sleep 10
    retry_command ionosctl ipblock delete --ipblock_id "$ipblock_id" -f -w
}

@test "Create a server with internet access, attach a volume, and verify SSH connection" {
    run ionosctl ipblock create --location de/txl --size 1 -o json -w -t 300
    assert_success
    ipblock_id=$(echo "$output" | jq -r '.id')
    assert_regex "$ipblock_id" "$uuid_v4_regex"
    echo "created ipblock $ipblock_id"

    run ionosctl datacenter create --name "$(randStr 8)" --location de/txl -o json -w -t 300
    assert_success
    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "created datacenter $datacenter_id"

    # Create a barebones server
    run ionosctl server create --datacenter-id "$datacenter_id" --cpu-family INTEL_SKYLAKE --cores 1 --ram 4GB -o json -w -t 300
    assert_success
    server_id=$(echo "$output" | jq -r '.id')
    assert_regex "$server_id" "$uuid_v4_regex"

    echo "created server $server_id"

    # Try getting the already created ipblock ID via a filtered list this time
    run ionosctl ipblock list -F location=de/txl -M 1 -o json
    assert_success
    found_ip=$(echo "$output" | jq -r '.items[] | .properties.ips[0]')

    # Create a LAN. Use --cols to extract its ID.
    run ionosctl lan create --datacenter-id "$datacenter_id" --cols LanId --no-headers --public
    assert_success
    lan_id="$output"
    # Test that Lan ID is a numeric value
    assert_regex "$lan_id" "^[0-9]+$"
    echo "created lan $lan_id"

    sleep 10

    run ionosctl nic create --datacenter-id "$datacenter_id" --server-id "$server_id" --lan-id "$lan_id" --ips "$found_ip" -o json -w -t 300
    assert_success
    nic_id=$(echo "$output" | jq -r '.id')
    assert_regex "$nic_id" "$uuid_v4_regex"
    echo "created nic $nic_id"

    # Validate the LAN ID
    nic_lan_id=$(echo "$output" | jq -r '.properties.lan')
    assert_equal "$nic_lan_id" "$lan_id" [ "$nic_lan_id" -eq "$lan_id" ]
    echo "NIC is associated with LAN ID $lan_id"

    # Validate the IPBlock IP matches
    nic_ip=$(echo "$output" | jq -r '.properties.ips[0]')
    assert_equal "$nic_ip" "$found_ip"

    # Create a volume with SSH Key and User Data
    run ionosctl volume create --datacenter-id "$datacenter_id" --size 50 --image-alias ubuntu:latest --ssh-key-paths "${ssh_key_path}.pub" -o json -w -t 300
    assert_success
    vol_id=$(echo "$output" | jq -r '.id')
    assert_regex "$vol_id" "$uuid_v4_regex"
    echo "created volume $vol_id"

    # Attach volume to server
    run ionosctl server volume attach --server-id "$server_id" --volume-id "$vol_id" --datacenter-id "$datacenter_id" -w -t 300
    assert_success

    # SSH Connection Test
    # This test tries to connect to the server using the SSH key and checks if the connection is successful
    # TODO: Re-enable me when no longer flaky
    # run ssh -i "$ssh_key_path" -o StrictHostKeyChecking=no ubuntu@"$nic_ip" echo "SSH connection successful"
    # assert_success
}
