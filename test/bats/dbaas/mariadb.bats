#!/usr/bin/env bats

# tags: maria

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/txl"

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

@test "Create MariaDB Cluster" {
    datacenter_id=$(find_or_create_resource \
        "ionosctl datacenter list -M 1 -F location=${location},state=available -o json 2> /dev/null | jq -r '.items[] | .id'" \
        "ionosctl datacenter create --name \"CLI-Test-$(randStr 8)\" --location ${location} -o json 2> /dev/null | jq -r '.id'")
    [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
    assert_regex "$datacenter_id" "$uuid_v4_regex"

    lan_id=$(find_or_create_resource \
        "ionosctl lan list -M 1 --datacenter-id ${datacenter_id} -F public=false-o json 2> /dev/null | jq -r '.items[] | .id'" \
        "sleep 30 && ionosctl lan create --datacenter-id ${datacenter_id} --public=false -o json 2> /dev/null | jq -r '.id'")
    [ -n "$lan_id" ] || fail "lan_id is empty"

    sleep 60

    echo "Trying to create mariadb cluster in datacenter $datacenter_id"
    run ionosctl db mariadb cluster create --name "CLI-Test-$(randStr 6)" --version 10.6 \
      --user $(randStr 12) --password $(randStr 12) --datacenter-id ${datacenter_id} --lan-id 1 --cidr 192.168.1.127/24 -o json 2> /dev/null
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created mariadb cluster $cluster_id"
}

@test "Find MariaDB Cluster" {
    sleep 10

    run ionosctl db mariadb cluster list -n CLI --cols version --no-headers 2> /dev/null
    assert_output "10.6"
}

teardown_file() {
    echo "cleaning up datacenter $datacenter_id and mariadb cluster $cluster_id"
    retry_command run ionosctl dbaas mariadb cluster delete --cluster-id "$cluster_id" -f
    sleep 120
    retry_command run ionosctl datacenter delete --datacenter_id "$datacenter_id" -f -w -t 1200
}
