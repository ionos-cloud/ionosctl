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

    sleep 120

    echo "Trying to create mariadb cluster in datacenter $datacenter_id"
    run ionosctl db mariadb cluster create --name "CLI-Test-$(randStr 8)" --version 10.6 \
      --user $(randStr 8) --password $(randStr 8) json 2> /dev/null -W -t 600
    assert_success
    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created mariadb cluster $cluster_id"
}

teardown_file() {
    echo "cleaning up datacenter $datacenter_id and mariadb cluster $cluster_id"
    retry_command run ionosctl k8s nodepool delete --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -f -w -t 1200
    sleep 30
    retry_command run ionosctl k8s cluster delete --cluster-id "$cluster_id" -f -w -t 1200
    sleep 30
    retry_command run ionosctl datacenter delete --datacenter_id "$datacenter_id" -f -w -t 1200
}
