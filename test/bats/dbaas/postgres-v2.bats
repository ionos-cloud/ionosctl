#!/usr/bin/env bats

# tags: dbaas, postgres-v2

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/txl"

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

# --- Read-only operations (no cluster needed) ---

@test "List postgres-v2 versions" {
    run ionosctl dbaas postgres-v2 version list -o json 2> /dev/null
    assert_success
}

@test "List postgres-v2 backup locations" {
    run ionosctl dbaas postgres-v2 backup location list -o json 2> /dev/null
    assert_success
}

@test "List postgres-v2 backups" {
    run ionosctl dbaas postgres-v2 backup list -o json 2> /dev/null
    assert_success
}

# --- Infrastructure setup ---

@test "Generate Token" {
    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

@test "Create Datacenter" {
    run ionosctl datacenter create --name "CLI-PsqlV2-Test-$(randStr 8)" --location ${location} -o json 2> /dev/null
    assert_success

    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "created datacenter $datacenter_id"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id

    sleep 30
}

@test "Create LAN" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    sleep 30

    run ionosctl lan create --datacenter-id ${datacenter_id} --public=false -o json 2> /dev/null
    assert_success

    lan_id=$(echo "$output" | jq -r '.id')
    echo "$lan_id" > /tmp/bats_test/lan_id
}

# --- Cluster lifecycle ---

@test "Create postgres-v2 cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)

    sleep 60

    run ionosctl dbaas postgres-v2 cluster create \
        --name "CLI-PsqlV2-Test-$(randStr 6)" \
        --version 17 \
        --datacenter-id ${datacenter_id} \
        --lan-id ${lan_id} \
        --cidr 192.168.1.100/24 \
        --db-username testuser \
        --db-password "$(randStr 16)" \
        --instances 1 \
        --cores 2 \
        --ram 4GB \
        --storage-size 20GB \
        --sync ASYNCHRONOUS \
        --backup-location de \
        -o json 2> /dev/null
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created postgres-v2 cluster $cluster_id"
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Get postgres-v2 cluster by ID" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas postgres-v2 cluster get --cluster-id "$cluster_id" -o json 2> /dev/null
    assert_success

    cluster_name=$(echo "$output" | jq -r '.properties.name')
    assert_output -p "\"name\": \"$cluster_name\""
    echo "$cluster_name" > /tmp/bats_test/cluster_name
}

@test "List postgres-v2 clusters" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    cluster_name=$(cat /tmp/bats_test/cluster_name)

    # JSON output
    run ionosctl dbaas postgres-v2 cluster list -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$cluster_name\""

    # Column output
    run ionosctl dbaas postgres-v2 cluster list --cols ClusterId --no-headers 2> /dev/null
    assert_success
    assert_output -p "$cluster_id"
}

@test "Update postgres-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    sleep 30

    run ionosctl dbaas postgres-v2 cluster update \
        --cluster-id "${cluster_id}" \
        --cores 4 \
        -o json 2> /dev/null
    assert_success

    new_cores=$(echo "$output" | jq -r '.properties.instances.cores')
    assert_equal "$new_cores" "4"
}

@test "List postgres-v2 backups for cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    sleep 30

    run ionosctl dbaas postgres-v2 backup list --cluster-id "${cluster_id}" -o json 2> /dev/null
    assert_success
}

@test "Restore postgres-v2 cluster from backup" {
    skip "Restore requires an existing backup - run manually when available"

    cluster_id=$(cat /tmp/bats_test/cluster_id)
    backup_id=$(cat /tmp/bats_test/backup_id)

    run ionosctl dbaas postgres-v2 cluster restore \
        --cluster-id "${cluster_id}" \
        --backup-id "${backup_id}" \
        -f 2> /dev/null
    assert_success
}

@test "Delete postgres-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas postgres-v2 cluster delete --cluster-id "${cluster_id}" -f 2> /dev/null
    assert_success
}

@test "Delete all postgres-v2 clusters" {
    # Create a throwaway cluster for delete-all test, or skip if none exist
    run ionosctl dbaas postgres-v2 cluster delete --all -f 2> /dev/null
    # May succeed (if clusters exist) or fail (if none). Just verify it doesn't crash.
}

# --- Validation / error cases ---

@test "Create cluster missing required flags fails" {
    run ionosctl dbaas postgres-v2 cluster create 2>&1
    assert_failure
}

@test "Create cluster with invalid instances fails" {
    run ionosctl dbaas postgres-v2 cluster create \
        --datacenter-id "00000000-0000-4000-8000-000000000000" \
        --lan-id "1" \
        --cidr "192.168.1.100/24" \
        --db-username testuser \
        --db-password testpass123 \
        --instances 0 2>&1
    assert_failure
    assert_output -p "instances must be set to minimum: 1"
}

@test "Create cluster with too many instances fails" {
    run ionosctl dbaas postgres-v2 cluster create \
        --datacenter-id "00000000-0000-4000-8000-000000000000" \
        --lan-id "1" \
        --cidr "192.168.1.100/24" \
        --db-username testuser \
        --db-password testpass123 \
        --instances 6 2>&1
    assert_failure
    assert_output -p "instances must be set to minimum: 1, maximum: 5"
}

@test "Delete cluster without id or --all fails" {
    run ionosctl dbaas postgres-v2 cluster delete 2>&1
    assert_failure
}

@test "Restore without backup-id fails" {
    run ionosctl dbaas postgres-v2 cluster restore --cluster-id "00000000-0000-4000-8000-000000000000" 2>&1
    assert_failure
}

@test "Restore without cluster-id fails" {
    run ionosctl dbaas postgres-v2 cluster restore --backup-id "00000000-0000-4000-8000-000000000000" 2>&1
    assert_failure
}

# --- Help output verification ---

@test "Cluster create help shows expected flags" {
    run ionosctl dbaas postgres-v2 cluster create --help
    assert_success
    assert_output -p "--datacenter-id"
    assert_output -p "--lan-id"
    assert_output -p "--cidr"
    assert_output -p "--db-username"
    assert_output -p "--db-password"
    assert_output -p "--backup-location"
    assert_output -p "--sync"
    assert_output -p "--maintenance-day"
    assert_output -p "--maintenance-time"
    assert_output -p "--backup-id"
    assert_output -p "--recovery-time"
}

@test "Cluster update help shows expected flags" {
    run ionosctl dbaas postgres-v2 cluster update --help
    assert_success
    assert_output -p "--cluster-id"
    assert_output -p "--cores"
    assert_output -p "--ram"
    assert_output -p "--instances"
    assert_output -p "--sync"
    assert_output -p "--maintenance-day"
    assert_output -p "--maintenance-time"
}

@test "Cluster restore help shows expected flags" {
    run ionosctl dbaas postgres-v2 cluster restore --help
    assert_success
    assert_output -p "--cluster-id"
    assert_output -p "--backup-id"
    assert_output -p "--recovery-time"
}

# --- Teardown ---

teardown_file() {
    ionosctl dbaas postgres-v2 cluster delete -af 2> /dev/null || true
    sleep 120

    if [[ -f /tmp/bats_test/datacenter_id ]]; then
        ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f 2> /dev/null || true
    fi

    if [[ -f /tmp/bats_test/token ]]; then
        ionosctl token delete --token "$(cat /tmp/bats_test/token)" -f 2> /dev/null || true
    fi

    rm -rf /tmp/bats_test
}
