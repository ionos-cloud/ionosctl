#!/usr/bin/env bats

# tags: postgres, pg

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/fra"

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

@test "Generate Token" {
    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

@test "Create Datacenter" {
    run ionosctl datacenter create --name "CLI-Test-$(randStr 8)" --location ${location} -w -o json 2> /dev/null
    assert_success

    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "created datacenter $datacenter_id"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
}

@test "Create LAN" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)

    sleep 30

    run ionosctl lan create --datacenter-id ${datacenter_id} --public=false -w -o json 2> /dev/null
    assert_success

    lan_id=$(echo "$output" | jq -r '.id')
    assert_regex "$lan_id" "$uuid_v4_regex"
    echo "created lan $lan_id"
    echo "$lan_id" > /tmp/bats_test/lan_id
}

@test "Create Postgres Cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)

    run ionosctl dbaas postgres cluster create --datacenter-id "$datacenter_id" --lan-id "$lan_id" \
      --cidr 192.168.1.127/24 --db-username testuser1234 --db-password "$(randStr 12)" -W -o json 2> /dev/null
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')

    assert_output -p "\"postgresVersion\": \"15\""
    assert_regex "$cluster_id" "$uuid_v4_regex"

    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Get Postgres Cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Finding postgres cluster $cluster_id"

    run ionosctl dbaas postgres cluster get --cluster-id "$cluster_id" -o json 2> /dev/null
    assert_output -p "\"id\": \"$cluster_id\""
    assert_output -p "\"postgresVersion\": \"15\""
    assert_success
}

@test "Create Postgres User" {
    postgres_user="CLI-User-$(randStr 8)"
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Creating user for postgres cluster $cluster_id"

    run ionosctl dbaas postgres user create --cluster-id "$cluster_id" --user "$postgres_user" --password "$(randStr 10)" -o json 2> /dev/null
    assert_success
    assert_output -p "\"username\": \"$postgres_user\""

    user_name=$(echo "$output" | jq -r '.properties.username')
    echo "created postgres user $user_name"
    echo "$user_name" > /tmp/bats_test/user_name
}

@test "List Postgres Users" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    username=$(cat /tmp/bats_test/user_name)
    echo "Listing users for postgres cluster $cluster_id"

    run ionosctl dbaas postgres user list --cluster-id "$cluster_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"username\": \"$username\""
}

@test "Create Postgres Database" {
    user=$(cat /tmp/bats_test/user_name)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    name="test-dbname-$(randStr 6)"
    echo $name > /tmp/bats_test/db_name

    run ionosctl dbaas postgres database create --cluster-id "$cluster_id" --database "$name" --owner "$user" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$name\""
}

@test "List Postgres Database" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    name=$(cat /tmp/bats_test/db_name)

    run ionosctl dbaas postgres database list -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$name\""
}

@test "Delete Postgres Database" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    name=$(cat /tmp/bats_test/db_name)

    run ionosctl dbaas postgres database delete --cluster-id "$cluster_id" --database "$name" -f 2> /dev/null
    assert_success
}

@test "Delete Postgres User" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    user_name=$(cat /tmp/bats_test/user_name)

    echo "Deleting postgres user $user_name from cluster $cluster_id"
    run ionosctl dbaas postgres user delete --cluster-id "$cluster_id" --user "$user_name" -f 2> /dev/null
    assert_success
}

teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    user_name=$(cat /tmp/bats_test/user_name)

    run ionosctl dbaas postgres user delete --cluster-id "$cluster_id" -af
    run ionosctl dbaas postgres cluster delete -af
    run ionosctl datacenter delete --datacenter-id "$datacenter_id" -f 2> /dev/null

    echo "cleaning up token"
    run ionosctl token delete --token "$(cat /tmp/bats_test/token)" -f
    unset IONOS_TOKEN

    echo "cleaning up test directory"
    rm -rf /tmp/bats_test
}
