#!/usr/bin/env bats

# tags: mongo, mongodb

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/txl"

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
    export IONOS_TOKEN=$(ionosctl token generate --ttl 6h)

    mkdir -p /tmp/bats_test
}

@test "Create MongoDB Cluster" {
    datacenter_id=$(ionosctl datacenter create --name "CLI-Test-$(randStr 8)" --location ${location} -o json 2> /dev/null | jq -r '.id')
    [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
    assert_regex "$datacenter_id" "$uuid_v4_regex"

    retry_until "ionosctl datacenter get --datacenter-id $datacenter_id -o json 2> /dev/null | jq -r '.metadata.state'" \
        "[[ \$output == \"AVAILABLE\" ]]" 10 60

    lan_id=$(ionosctl lan create -w --datacenter-id "${datacenter_id}" --public=false -o json 2> /dev/null | jq -r '.id')
    [ -n "$lan_id" ] || fail "lan_id is empty"

    echo "Trying to create MongoDB cluster in datacenter $datacenter_id"
    run ionosctl db mongo cluster create --name "CLI-Test-$(randStr 6)" --edition playground \
        --datacenter-id "${datacenter_id}" --lan-id 1 --cidr 192.168.1.127/24 -o json 2> /dev/null
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created mongodb cluster $cluster_id"

    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Find MongoDB Cluster and wait until it is available" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Finding mongodb cluster $cluster_id"

    run ionosctl db mongo cluster get --cluster-id "$cluster_id" -o json 2> /dev/null
    assert_success

}

@test "Create MongoDB User" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Creating user for mongodb cluster $cluster_id"

    run ionosctl db mongo user create --cluster-id "$cluster_id" --name "CLI-Test-$(randStr 8)" --password "$(randStr 10)" \
        --roles db=read -o json 2> /dev/null
    assert_success

    user_name=$(echo "$output" | jq -r '.properties.name')
    echo "created mongodb user $user_name"
    echo "$user_name" > /tmp/bats_test/user_name
}

@test "List MongoDB Users" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Listing users for mongodb cluster $cluster_id"

    run ionosctl db mongo user list --cluster-id "$cluster_id" -o json 2> /dev/null
    assert_success
}

teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    user_name=$(cat /tmp/bats_test/user_name)

    echo "cleaning up mongodb cluster $cluster_id"
    run ionosctl db mongo user delete --cluster-id "$cluster_id" --name "$user_name" -f
    retry_until "ionosctl db mongo user get --cluster-id $cluster_id --user $user_name --database db 2> /dev/null" \
        "[[ \$output == \"\" ]]" 10 60

    run ionosctl dbaas mongo cluster delete --cluster-id "$cluster_id" -f
    retry_until "ionosctl dbaas mongo cluster get --cluster-id $cluster_id 2> /dev/null" \
        "[ -z \$output ]" 120 60

    echo "cleaning up datacenter $datacenter_id"
    retry_until "ionosctl datacenter delete --datacenter-id $datacenter_id -f 2> /dev/null" \
        "[ \$? -eq 0 ]" 10 60

    echo "cleaning up token"
    run ionosctl token delete --token "$IONOS_TOKEN" -f
    unset IONOS_TOKEN

    echo "cleaning up test directory"
    rm -rf /tmp/bats_test
}
