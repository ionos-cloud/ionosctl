#!/usr/bin/env bats

# tags: in-memory-db, inmemorydb

BATS_LIBS_PATH="${LIBS_PATH:-..//libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/fra"

setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    mkdir -p /tmp/bats_test
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

@test "Create datacenter and LAN" {
    run ionosctl datacenter create --name "CLI-Test-InMemDb-$(randStr 6)" --location "${location}" -o json
    assert_success
    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id

    # wait for datacenter to be AVAILABLE
    sleep 60

    run ionosctl lan create --datacenter-id "$datacenter_id" --public=false -o json
    assert_success
    lan_id=$(echo "$output" | jq -r '.id')
    assert_regex "$lan_id" "$uuid_v4_regex"
    echo "$lan_id" > /tmp/bats_test/lan_id

    # wait until LAN state is AVAILABLE
    status=""
    for i in $(seq 1 30); do
      status=$(ionosctl lan get --datacenter-id "$datacenter_id" --lan-id "$lan_id" -o json | jq -r '.metadata.state')
      [ "$status" = "AVAILABLE" ] && break
      sleep 10
    done
    [ "$status" = "AVAILABLE" ] || fail "LAN did not become AVAILABLE"
}

@test "Create" {
    replica_name="cli-imdb-$(randStr 6 | tr '[:upper:]' '[:lower:]')"
    db_user="user$(randStr 4)"
    db_pass="pass$(randStr 6)"

    run ionosctl db in-memory-db replicaset create \
      --location "${location}" \
      --name "$replica_name" \
      --replicas 1 \
      --cores 1 \
      --ram 4GB \
      --user "$db_user" \
      --password "$db_pass" \
      --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
      --lan-id "$(cat /tmp/bats_test/lan_id)" \
      --cidr "192.168.1.70/24" \
      -o json
    assert_success

    replicaset_id=$(echo "$output" | jq -r '.id')
    assert_regex "$replicaset_id" "$uuid_v4_regex"
    user=$(echo "$output" | jq -r '.properties.credentials.username')
    assert_equal "$user" "$db_user"
    echo "$replicaset_id" > /tmp/bats_test/replicaset_id
    echo "$replica_name" > /tmp/bats_test/replicaset_name
}

@test "List replicasets" {
    run ionosctl db in-memory-db replicaset list --location "${location}" -o json
    assert_success
    assert_output -p "\"displayName\": \"$(cat /tmp/bats_test/replicaset_name)\""
}

@test "Get" {
    run ionosctl db in-memory-db replicaset get --location "${location}" --replica-set-id "$(cat /tmp/bats_test/replicaset_id)" -o json
    assert_success
    assert_output -p "\"displayName\": \"$(cat /tmp/bats_test/replicaset_name)\""
}

@test "List available snapshots" {
    run ionosctl db in-memory-db snapshot list --location "${location}" -o json
    assert_success

    snapshot_id=$(echo "$output" | jq -r '.items[0].id')
    if [ -z "$snapshot_id" ]; then
        echo "" > /tmp/bats_test/snapshot_id
        skip "No snapshots available in the location $location"

    fi

    assert_regex "$snapshot_id" "$uuid_v4_regex"
    echo "$snapshot_id" > /tmp/bats_test/snapshot_id

    datacenter_id=$(echo "$output" | jq -r '.items[0].metadata.datacenterId')
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id_2
}

@test "Create from snapshot" {
    snapshot_id=$(cat /tmp/bats_test/snapshot_id)
    if [ -z "$snapshot_id" ]; then
        skip "No snapshots available in the location $location"
    fi

    datacenter_id=$(cat /tmp/bats_test/datacenter_id_2)
    if [ -z "$datacenter_id" ]; then
        skip "No datacenter ID available for snapshot creation"
    fi

    # find a private lan ID in the datacenter
    run ionosctl lan list --datacenter-id "$datacenter_id" -F public=false -o json
    if [[ "$status" -ne 0 ]] && [[ "$output" == *"Resource does not exist"* ]]; then
        skip "skipping replica creation from snapshot due to flaky API response regarding datacenter $datacenter_id"
    fi
    assert_success
    lan_id=$(echo "$output" | jq -r '.items[0].id')
    if [ -z "$lan_id" ]; then
        skip "No private LAN found in the datacenter $datacenter_id"
    fi

    run ionosctl db in-memory-db replicaset create \
      --location "${location}" \
      --snapshot-id "$snapshot_id" \
      --name "cli-imdb-snapshot-$(randStr 4 | tr '[:upper:]' '[:lower:]')" \
      --replicas 1 \
      --cores 1 \
      --ram 4GB \
      --user "snapshot_user" \
      --password "snapshot_pass$(randStr 2)" \
      --datacenter-id "$datacenter_id" \
      --lan-id "$(cat /tmp/bats_test/lan_id)" \
      --cidr "192.168.1.70/24" \
      -o json

    if [[ "$status" -ne 0 ]] && [[ "$output" == *"Resource does not exist"* ]]; then
        skip "skipping replica creation from snapshot due to flaky API response regarding datacenter $datacenter_id"
    fi
    assert_success

    replicaset_id=$(echo "$output" | jq -r '.id')
    assert_regex "$replicaset_id" "$uuid_v4_regex"
    echo "$replicaset_id" > /tmp/bats_test/replicaset_id_2
    replicaset_name=$(echo "$output" | jq -r '.properties.displayName')
    echo "$replicaset_name" > /tmp/bats_test/replicaset_name_2
}

@test "List and expect both replicasets" {
    if [ ! -f /tmp/bats_test/replicaset_name_2 ]; then
        skip "No second replicaset created from snapshot"
    fi

    run ionosctl db in-memory-db replicaset list --location "${location}" -o json
    assert_success
    assert_output -p "\"displayName\": \"$(cat /tmp/bats_test/replicaset_name)\""
    assert_output -p "\"displayName\": \"$(cat /tmp/bats_test/replicaset_name_2)\""
}

@test "Update" {
    skip "replicasets cannot be updated at the moment due to an API bug"

    run ionosctl db in-memory-db replicaset update \
      --location "${location}" \
      --replica-set-id "$(cat /tmp/bats_test/replicaset_id)" \
      --user "updated_user$(randStr 4)" \
      --ram 6GB \
      -o json
    assert_success
    assert_output -p "\"username\": \"updated_user"
    assert_output -p "\"ram\": 6"
}

@test "Delete" {
    run ionosctl db in-memory-db replicaset delete \
      --location "${location}" \
      --replica-set-id "$(cat /tmp/bats_test/replicaset_id)" \
      -f
    assert_success

    if [ -f /tmp/bats_test/replicaset_id_2 ]; then
        run ionosctl db in-memory-db replicaset delete \
          --location "${location}" \
          --replica-set-id "$(cat /tmp/bats_test/replicaset_id_2)" \
          -f
        assert_success
    fi
}

teardown_file() {
    ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f

    ionosctl token delete --token "$IONOS_TOKEN" --force

    ionosctl db in-memory-db replicaset delete --location "${location}" -af

    rm -rf /tmp/bats_test
}
