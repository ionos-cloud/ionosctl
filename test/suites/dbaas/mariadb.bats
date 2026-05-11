#!/usr/bin/env bats

# paths: commands/dbaas/mariadb/*

load '../setup.bats'

location="de/txl"

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}


@test "Generate Token" {
    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

@test "Create Datacenter" {
    run ionosctl datacenter create --name "CLI-Test-$(randStr 8)" --location ${location} -o json --wait
    assert_success

    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "created datacenter $datacenter_id"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
}

@test "Create LAN" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)

    run ionosctl lan create --datacenter-id ${datacenter_id} --public=false -o json --wait
    assert_success

    lan_id=$(echo "$output" | jq -r '.id')
    echo "$lan_id" > /tmp/bats_test/lan_id
}

@test "Create MariaDB Cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)

    sleep 10

    run ionosctl dbaas mariadb cluster create --name "CLI-Test-$(randStr 6)" --version 10.6 --user testuser1234 \
       --password "$(randStr 12)" --datacenter-id ${datacenter_id} --lan-id ${lan_id} --cidr 192.168.1.127/24 -o json -w
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created mariadb cluster $cluster_id"
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Get MariaDB Cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    # Get cluster by ID
    run ionosctl dbaas mariadb cluster get --cluster-id "$cluster_id" -o json -w
    assert_success
    cluster_name=$(echo "$output" | jq -r '.properties.displayName')
    assert_output -p "\"displayName\": \"$cluster_name\""
    echo "$cluster_name" > /tmp/bats_test/cluster_name
}

@test "List MariaDB Backups" {
    skip "Skipping temporarily because flaky test failures"

    cluster_id=$(cat /tmp/bats_test/cluster_id)

    sleep 30

    # List all backups
    run ionosctl dbaas mariadb backup list
    assert_success

    # List backups for specific cluster
    run ionosctl dbaas mariadb backup list --cluster-id "${cluster_id}"
    assert_success
}

@test "List MariaDB Clusters" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    cluster_name=$(cat /tmp/bats_test/cluster_name)

    # List clusters (JSON output)
    run ionosctl dbaas mariadb cluster list -o json
    assert_success
    assert_output -p "\"displayName\": \"$cluster_name\""

    # List clusters (Column output)
    run ionosctl dbaas mariadb cluster list --cols ClusterId --no-headers
    assert_success
    assert_output -p "$cluster_id"
}

@test "Update MariaDB cluster maintenance day" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    sleep 10

    run ionosctl dbaas mariadb cluster update --cluster-id "${cluster_id}" \
      --maintenance-day Wednesday --maintenance-time 12:00:00 -o json -w
    assert_success

    new_day=$(echo "$output" | jq -r '.properties.maintenanceWindow.dayOfTheWeek')
    assert_equal "$new_day" "Wednesday"
    new_time=$(echo "$output" | jq -r '.properties.maintenanceWindow.time')
    assert_equal "$new_time" "12:00:00"
}

@test "Verify MariaDB Cluster DNS Resolution" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    # Extract the DNS and CIDR from the JSON output
    clusters_json=$(ionosctl dbaas mariadb cluster get --cluster-id "${cluster_id}" -o json -w)
    dns_name=$(echo "$clusters_json" | jq -r '.properties.dnsName')
    cidr=$(echo "$clusters_json" | jq -r '.properties.connections[0].cidr')

    echo "Looking up DNS name: $dns_name expecting an IP in CIDR: $cidr"
    ip=$(dig +short "$dns_name")
    echo "dig result: $ip"

    run ip_in_cidr "$ip" "$cidr"
    assert_success
}

@test "Delete MariaDB Cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas mariadb cluster delete --cluster-id "$cluster_id" -o json -w -f
    assert_success
}


teardown_file() {
    ionosctl dbaas mariadb cluster delete -af
    sleep 120

    ionosctl datacenter delete -af
    ionosctl token delete --token "$(cat /tmp/bats_test/token)" -f

    rm -rf /tmp/bats_test
}
