#!/usr/bin/env bats

# tags: cloudapi-v6, datacenter, server, lan, nic, ipblock, volume

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup() {
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

#!/usr/bin/env bats

# tags: cloudapi-v6, datacenter, server, lan, nic, ipblock, volume

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="es/vit"

setup_file() {
    ssh_key_path=$(generate_ssh_key)
    ssh_pub_key=$(cat "${ssh_key_path}.pub")
}

teardown_file() {
    rm -f "${ssh_key_path}" "${ssh_key_path}.pub"
    echo "cleaning up datacenter $datacenter_id and ipblock $ipblock_id"
    ionosctl datacenter delete --datacenter-id "$datacenter_id" -f -w
    sleep 10
    retry_command ionosctl ipblock delete --ipblock_id "$ipblock_id" -f -w
}

@test "Create a server with internet access, attach a volume, and verify SSH connection" {
    ipblock_id=$(find_or_create_resource \
        "ionosctl ipblock list -F location=${location} -M 1 -o json 2> /dev/null | jq -r '.items[] | .properties.ips[0]'" \
        "ionosctl ipblock create --location ${location} --size 1 -o json 2> /dev/null -w -t 300 | jq -r '.id'")
    [ -n "$ipblock_id" ] || fail "ipblock_id is empty"
    assert_regex "$ipblock_id" "$uuid_v4_regex"
    echo "created/found ipblock $ipblock_id"

    datacenter_id=$(find_or_create_resource \
        "ionosctl datacenter list -F location=${location} -M 1 -o json 2> /dev/null | jq -r '.items[] | .id'" \
        "ionosctl datacenter create --name \"$(randStr 8)\" --location ${location} -o json 2> /dev/null -w -t 300 | jq -r '.id'")
    [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "created/found datacenter $datacenter_id"

    # Assuming create server is always a new operation
    run ionosctl server create --datacenter-id "$datacenter_id" --cpu-family INTEL_SKYLAKE --cores 1 --ram 4GB -o json 2> /dev/null -w -t 300
    assert_success
    server_id=$(echo "$output" | jq -r '.id')
    assert_regex "$server_id" "$uuid_v4_regex"
    echo "created server $server_id"

    # Assuming create LAN is always a new operation
    run ionosctl lan create --datacenter-id "$datacenter_id" --public -o json 2> /dev/null -w -t 300
    assert_success
    lan_id=$(echo "$output" | jq -r '.id')
    assert_regex "$lan_id" "^[0-9]+$"
    echo "created lan $lan_id"

    # Assuming create NIC is always a new operation
    run ionosctl nic create --datacenter-id "$datacenter_id" --server-id "$server_id" --lan-id "$lan_id" --ips "$(echo "$ipblock_id" | jq -r '.properties.ips[0]')" -o json 2> /dev/null -w -t 300
    assert_success
    nic_id=$(echo "$output" | jq -r '.id')
    assert_regex "$nic_id" "$uuid_v4_regex"
    echo "created nic $nic_id"

    # Assuming create volume is always a new operation
    run ionosctl volume create --datacenter-id "$datacenter_id" --size 50 --image-alias ubuntu:latest --ssh-key-paths "${ssh_key_path}.pub" -o json 2> /dev/null -w -t 300
    assert_success
    vol_id=$(echo "$output" | jq -r '.id')
    assert_regex "$vol_id" "$uuid_v4_regex"
    echo "created volume $vol_id"

    # Assuming attach volume to server is always a new operation
    run ionosctl server volume attach --server-id "$server_id" --volume-id "$vol_id" --datacenter-id "$datacenter_id" -w -t 300
    assert_success

    # SSH Connection Test
    # This test tries to connect to the server using the SSH key and checks if the connection is successful
    # TODO: Re-enable me when no longer flaky
    # run ssh -i "$ssh_key_path" -o StrictHostKeyChecking=no ubuntu@"$nic_ip" echo "SSH connection successful"
    # assert_success
}
