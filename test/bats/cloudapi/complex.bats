#!/usr/bin/env bats

setup() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

randStr() {
    size=$1
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c "$size"
}

generate_ssh_key() {
    ssh_key_path="/tmp/id_rsa_$(randStr 8)"
    ssh-keygen -t rsa -b 4096 -f "$ssh_key_path" -N "" >/dev/null 2>&1
    echo $ssh_key_path
}

encode_user_data() {
    echo "#!/bin/bash" > /tmp/user_data.sh
    echo "ip route add default via $1" >> /tmp/user_data.sh
    base64 /tmp/user_data.sh
}

@test "Create a server under a NAT" {
    # Generate SSH Key Pair
    ssh_key_path=$(generate_ssh_key)
    ssh_pub_key=$(cat "${ssh_key_path}.pub")


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

    sleep 10

    # Add the LAN to the NAT
    run ionosctl natgateway lan add --datacenter-id "$datacenter_id" --natgateway-id "$nat_id" --lan-id "$lan_id"
    [ "$status" -eq 0 ]
    echo "added lan $lan_id to nat $nat_id"

    sleep 60

    # Retrieve NAT Public IP and make sure it is the one of the ipblock
    run ionosctl natgateway list --datacenter-id "$datacenter_id" -M 1 --cols PublicIps --no-headers
    [ "$status" -eq 0 ]
    public_ip="$output"
    [[ $public_ip =~ $ip_regex ]]
    [[ $public_ip == $found_ip ]]

    # Retrieve Gateway IP with Mask
    run ionosctl natgateway lan list --datacenter-id "$datacenter_id" --natgateway-id "$nat_id" --cols GatewayIps --no-headers
    [ "$status" -eq 0 ]
    gateway_ip="$output"
    [[ $gateway_ip =~ $ip_regex ]]
    echo "Retrieved gateway IP: $gateway_ip"

    # TODO: This fails
    run ionosctl natgateway rule create \
        --datacenter-id "$datacenter_id" \
        --natgateway-id "$nat_id" \
        --ip "$found_ip" \
        --port-range-start 1 \
        --port-range-end 65534 \
        --protocol TCP \
        --source-subnet "0.0.0.0/0" \
        --target-subnet "10.0.0.0/24" \
        -w
    [ "$status" -eq 0 ]
    echo "NAT Gateway rule created"

    # Encode User Data
    user_data=$(encode_user_data $gateway_ip)

    # Create a volume with SSH Key and User Data
    run ionosctl volume create --datacenter-id "$datacenter_id" --size 50 --image-alias ubuntu:latest --ssh-key-paths "${ssh_key_path}.pub" --user-data "$user_data" -o json -w
    [ "$status" -eq 0 ]
    vol_id=$(echo "$output" | jq -r '.id')
    [[ $vol_id =~ $uuid_v4_regex ]]
    echo "created volume $vol_id"

    # Attach volume to server
    run ionosctl volume attach --server-id "$server_id" --volume-id "$vol_id" --datacenter-id "$datacenter_id" -w
    [ "$status" -eq 0 ]

    # TODO: get server_ip
    # SSH Connection Test
    # This test tries to connect to the server using the SSH key and checks if the connection is successful
    run ssh -i "$ssh_key_path" -o StrictHostKeyChecking=no ubuntu@"$server_ip" echo "SSH connection successful"
    [ "$status" -eq 0 ]
}
