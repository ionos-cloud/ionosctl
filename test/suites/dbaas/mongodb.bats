#!/usr/bin/env bats

# paths: commands/dbaas/mongo/*

load '../setup.bats'

location="de/txl"

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
    export IONOS_TOKEN=$(ionosctl token generate --ttl 6h)

    mkdir -p /tmp/bats_test
}

@test "Create MongoDB Cluster" {
    run ionosctl datacenter create -w --name "CLI-Test-$(randStr 8)" --location ${location} -o json
    assert_success
    datacenter_id=$(echo "$output" | jq -r '.id')
    [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
    assert_regex "$datacenter_id" "$uuid_v4_regex"

    run ionosctl lan create -w --datacenter-id "${datacenter_id}" --public=false -o json
    assert_success
    lan_id=$(echo "$output" | jq -r '.id')
    [ -n "$lan_id" ] || fail "lan_id is empty"

    echo "Trying to create MongoDB cluster in datacenter $datacenter_id"
    run ionosctl db mongo cluster create --name "CLI-Test-$(randStr 6)" --edition playground \
        --datacenter-id "${datacenter_id}" --lan-id 1 --cidr 192.168.1.127/24 -w --timeout 2400 -o json
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created mongodb cluster $cluster_id"

    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Get MongoDB Cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Finding mongodb cluster $cluster_id"

    run ionosctl db mongo cluster get --cluster-id "$cluster_id" -o json
    assert_success
}

@test "Create MongoDB User" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Creating user for mongodb cluster $cluster_id"

    run ionosctl db mongo user create --cluster-id "$cluster_id" --name "CLI-Test-$(randStr 8)" --password "$(randStr 10)" \
        --roles db=read -o json
    assert_success

    user_name=$(echo "$output" | jq -r '.properties.username')
    echo "created mongodb user $user_name"
    echo "$user_name" > /tmp/bats_test/user_name
}

@test "List MongoDB Users" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Listing users for mongodb cluster $cluster_id"

    run ionosctl db mongo user list --cluster-id "$cluster_id" -o json
    assert_success
}

@test "Delete MongoDB User" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    user_name=$(cat /tmp/bats_test/user_name)

    echo "Deleting mongodb user $user_name from cluster $cluster_id"
    run ionosctl db mongo user delete --cluster-id "$cluster_id" --name "$user_name" -f
    assert_success
}

@test "Change MongoDB Cluster Name" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Patching mongodb cluster $cluster_id"

    new_name="CLI-Test-$(randStr 6)"

    run ionosctl db mongo cluster update --cluster-id "$cluster_id" --name "$new_name" -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.displayName')" "$new_name"
}

@test "Change MongoDB Cluster Version to 7.0" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Patching mongodb cluster $cluster_id to version 7.0"

    run ionosctl db mongo cluster update --cluster-id "$cluster_id" --version 7.0 -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.mongoDBVersion')" "7.0"
}

@test "Delete MongoDB Cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    echo "Deleting mongodb cluster $cluster_id"
    run ionosctl dbaas mongo cluster delete --cluster-id "$cluster_id" -f
    assert_success
}

teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    user_name=$(cat /tmp/bats_test/user_name)

    run ionosctl db mongo user delete --cluster-id "$cluster_id" -af
    run ionosctl dbaas mongo cluster delete -af -w
    sleep 10
    run ionosctl datacenter delete --datacenter-id "$datacenter_id" -f -w

    echo "cleaning up token"
    run ionosctl token delete --token "$IONOS_TOKEN" -f
    unset IONOS_TOKEN

    echo "cleaning up test directory"
    rm -rf /tmp/bats_test
}
