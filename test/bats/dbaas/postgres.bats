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
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

@test "Create temporary sub-user with postgres permissions" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "test-user-$(randStr 4)" --last-name "test-last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl group create --name "test-group-$(randStr 4)" \
        -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    run ionosctl group user add --user-id "$(cat /tmp/bats_test/user_id)" \
        --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success

    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

@test "Create Postgres Cluster" {
    name="CLI-Test-$(randStr 8)"
    datacenter_id=$(ionosctl datacenter create  -w --name "$name" --location ${location} -o json 2> /dev/null | jq -r '.id')
    assert_success
    [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
    assert_output -p "\"name\": \"$name\""
    assert_output -p "\"location\": \"$location\""
    assert_regex "$datacenter_id" "$uuid_v4_regex"

    lan_id=$(ionosctl lan create -w --datacenter-id "${datacenter_id}" --name "$name" --public=false -o json 2> /dev/null | jq -r '.id')
    [ -n "$lan_id" ] || fail "lan_id is empty"
    assert_success
    assert_output -p "\"name\": \"$name\""

    sleep 10

    cluster_id=$(ionosctl dbaas postgres cluster create --name "$name" --datacenter-id "${datacenter_id}" \
     --lan-id 1 --cidr 192.168.1.127/24 --db-username "username$(randStr 6)" --db-password "pass$(randStr 6)" -o json 2> /dev/null | jq -r '.id')
    assert_success
    assert_output -p "\"name\": \"$name\""
    assert_regex "$cluster_id" "$uuid_v4_regex"

    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
    echo "$lan_id" > /tmp/bats_test/lan_id
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Get Postgres Cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Finding postgres cluster $cluster_id"

    run ionosctl dbaas postgres cluster get --cluster-id "$cluster_id" -o json 2> /dev/null
    assert_success

    sleep 30
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
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    name="test-dbname-$(randStr 6)"
    echo "$user_name" > /tmp/bats_test/user_name

    run ionosctl dbaas postgres database create --cluster-id "$cluster_id" --database "$name" --owner "test-dbowner-$(randStr 6)" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$name\""
}

@test "List Postgres Database" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    name="test-dbname-$(randStr 6)"

    run ionosctl dbaas postgres database list -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$name\""
}

@test "Delete Postgres Database" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    name="test-dbname-$(randStr 6)"
    echo "Listing users for postgres cluster $cluster_id"

    run ionosctl dbaas postgres database create --cluster-id "$cluster_id" --database "$name" --owner "test-dbowner-$(randStr 6)" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$name\""
}

@test "Delete Postgres User" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    user_name=$(cat /tmp/bats_test/user_name)

    echo "Deleting postgres user $user_name from cluster $cluster_id"
    run ionosctl dbaas postgres user delete --cluster-id "$cluster_id" --name "$user_name" -f 2> /dev/null
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
    run ionosctl token delete --token "$IONOS_TOKEN" -f
    unset IONOS_TOKEN

    echo "cleaning up test directory"
    rm -rf /tmp/bats_test
}
