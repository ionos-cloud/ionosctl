#!/usr/bin/env bats

# tags: in-memory-db, inmemorydb

BATS_LIBS_PATH="${LIBS_PATH:-..//libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="es/vit"

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

@test "List available In-Memory DB snapshots" {
    run ionosctl db in-memory-db snapshot list --location "${location}" -o json
    assert_success

    snapshot_id=$(echo "$output" | jq -r '.[0].id')
    assert_regex "$snapshot_id" "$uuid_v4_regex"
    echo "$snapshot_id" > /tmp/bats_test/snapshot_id
}

@test "Create In-Memory DB ReplicaSet" {
    replica_name="cli-imdb-$(randStr 6 | tr '[:upper:]' '[:lower:]')"
    db_user="user$(randStr 4)"
    db_pass="pass$(randStr 6)"

    run ionosctl db in-memory-db replicaset create \
      --location "${location}" \
      --name "$replica_name" \
      --replicas 1 \
      --cores 1 \
      --ram 4 \
      --user "$db_user" \
      --password "$db_pass" \
      --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
      --lan-id "$(cat /tmp/bats_test/lan_id)" \
      --cidr "192.168.1.70/24" \
      --snapshot-id "$snapshot_id" \
      -o json
    assert_success

    replicaset_id=$(echo "$output" | jq -r '.id')
    assert_regex "$replicaset_id" "$uuid_v4_regex"
    echo "$replicaset_id" > /tmp/bats_test/replicaset_id
    echo "$replica_name" > /tmp/bats_test/replicaset_name
}

@test "List In-Memory DB ReplicaSets" {
    run ionosctl db in-memory-db replicaset list --location "${location}" -o json
    assert_success
    assert_output -p "\"name\": \"$(cat /tmp/bats_test/replicaset_name)\""
}

@test "Get In-Memory DB ReplicaSet details" {
    run ionosctl db in-memory-db replicaset get --location "${location}" --replicaset "$(cat /tmp/bats_test/replicaset_id)" -o json
    assert_success
    assert_output -p "\"name\": \"$(cat /tmp/bats_test/replicaset_name)\""
}

@test "Create In-Memory DB ReplicaSet from snapshot" {
    snapshot_id=$(cat /tmp/bats_test/snapshot_id)
    run ionosctl db in-memory-db replicaset create \
      --location "${location}" \
      --snapshot-id "$snapshot_id" \
      --name "cli-imdb-snapshot-$(randStr 6 | tr '[:upper:]' '[:lower:]')" \
      --replicas 1 \
      --cores 1 \
      --ram 4 \
      --user "snapshot_user$(randStr 4)" \
      --password "snapshot_pass$(randStr 6)" \
      --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
      --lan-id "$(cat /tmp/bats_test/lan_id)" \
      --cidr "192.168.1.70/24" \
      -o json
    assert_success

    replicaset_id=$(echo "$output" | jq -r '.id')
    assert_regex "$replicaset_id" "$uuid_v4_regex"
    echo "$replicaset_id" > /tmp/bats_test/replicaset_id
    replicaset_name=$(echo "$output" | jq -r '.properties.displayName')
    assert_not_empty "$replicaset_name"
    echo "$replicaset_name" > /tmp/bats_test/replicaset_name
}

@test "List In-Memory DB ReplicaSets" {
    run ionosctl db in-memory-db replicaset list --location "${location}" -o json
    assert_success
    assert_output -p "\"name\": \"$(cat /tmp/bats_test/replicaset_name)\""
}

@test "Update In-Memory DB ReplicaSet" {
    run ionosctl db in-memory-db replicaset update \
      --location "${location}" \
      --replicaset "$(cat /tmp/bats_test/replicaset_id)" \
      --description "Updated via BATS" \
      -o json
    assert_success
    assert_output -p "\"description\": \"Updated via BATS\""
}

@test "Delete In-Memory DB ReplicaSet" {
    run ionosctl db in-memory-db replicaset delete \
      --location "${location}" \
      --replicaset "$(cat /tmp/bats_test/replicaset_id)" \
      -f
    assert_success
}

teardown_file() {
    ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f

    ionosctl token delete --token "$IONOS_TOKEN" --force

    ionosctl db in-memory-db replicaset delete --location "${location}" -af

    rm -rf /tmp/bats_test
}
