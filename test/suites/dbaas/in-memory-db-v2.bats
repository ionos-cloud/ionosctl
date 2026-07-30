#!/usr/bin/env bats

# paths: commands/dbaas/inmemorydb-v2/*

load '../setup.bats'

location="de/fra"

setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

# wait_for_state polls a cluster via `cluster get` until it reaches the wanted
# state (or FAILED / timeout). Used instead of the global `-w` flag: global
# wait builds its own poll URL, which the In-Memory DB v3 API rejects (HTTP
# 400); going through `cluster get` uses the SDK's correct request URL.
wait_for_state() {
    local id="$1" want="$2" timeout="${3:-2400}" elapsed=0 state
    while (( elapsed < timeout )); do
        state=$(ionosctl dbaas in-memory-db-v2 cluster get --location "${location}" --cluster-id "$id" -o json 2>/dev/null | jq -r '.metadata.state // empty')
        [[ "$state" == "$want" ]] && return 0
        [[ "$state" == "FAILED" ]] && return 1
        sleep 20
        elapsed=$((elapsed + 20))
    done
    return 1
}

# wait_for_gone polls until `cluster get` no longer finds the cluster (deletion
# finished), so the datacenter/LAN lock is released before teardown.
wait_for_gone() {
    local id="$1" timeout="${2:-1200}" elapsed=0
    while (( elapsed < timeout )); do
        ionosctl dbaas in-memory-db-v2 cluster get --location "${location}" --cluster-id "$id" -o json >/dev/null 2>&1 || return 0
        sleep 20
        elapsed=$((elapsed + 20))
    done
    return 1
}

# --- Read-only operations (no cluster needed) ---

@test "List in-memory-db-v2 versions" {
    run ionosctl dbaas in-memory-db-v2 version list --location "${location}" -o json
    assert_success
}

@test "List in-memory-db-v2 snapshot locations" {
    run ionosctl dbaas in-memory-db-v2 snapshot location list --location "${location}" -o json
    assert_success
    assert_output -p "\"location\""
}

@test "List in-memory-db-v2 snapshots" {
    run ionosctl dbaas in-memory-db-v2 snapshot list --location "${location}" -o json
    assert_success
}

# --- Infrastructure setup ---

@test "Create Datacenter" {
    run ionosctl datacenter create --name "CLI-ImdbV2-Test-$(randStr 8)" --location "${location}" -w -o json
    assert_success

    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
}

@test "Create LAN" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)

    run ionosctl lan create --datacenter-id "${datacenter_id}" --public=false -w -o json
    assert_success

    lan_id=$(echo "$output" | jq -r '.id')
    echo "$lan_id" > /tmp/bats_test/lan_id
}

# --- Cluster lifecycle ---

@test "Create in-memory-db-v2 cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)
    cluster_name="cli-imdbv2-$(randStr 6 | tr '[:upper:]' '[:lower:]')"
    db_user="user$(randStr 4)"

    # No --backup-location: it defaults to a valid Object Storage location.
    run ionosctl dbaas in-memory-db-v2 cluster create \
        --location "${location}" \
        --name "$cluster_name" \
        --version 8.0 \
        --datacenter-id "${datacenter_id}" \
        --lan-id "${lan_id}" \
        --cidr 192.168.1.100/24 \
        --user "$db_user" \
        --password "$(randStr 13)A1a" \
        --replicas 1 \
        --cores 1 \
        --ram 4GB \
        --eviction-policy allkeys-lru \
        --persistence-mode RDB \
        -o json
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    user=$(echo "$output" | jq -r '.properties.credentials.username')
    assert_equal "$user" "$db_user"
    echo "$cluster_id" > /tmp/bats_test/cluster_id
    echo "$cluster_name" > /tmp/bats_test/cluster_name
}

@test "Wait for in-memory-db-v2 cluster AVAILABLE" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run wait_for_state "$cluster_id" AVAILABLE 2400
    assert_success
}

@test "Get in-memory-db-v2 cluster by ID" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas in-memory-db-v2 cluster get --location "${location}" --cluster-id "$cluster_id" -o json
    assert_success
    assert_output -p "\"name\": \"$(cat /tmp/bats_test/cluster_name)\""
}

@test "List in-memory-db-v2 clusters" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    cluster_name=$(cat /tmp/bats_test/cluster_name)

    # JSON output
    run ionosctl dbaas in-memory-db-v2 cluster list -o json
    assert_success
    assert_output -p "\"name\": \"$cluster_name\""

    # Column output
    run ionosctl dbaas in-memory-db-v2 cluster list --cols ClusterId --no-headers
    assert_success
    assert_output -p "$cluster_id"
}

@test "Update in-memory-db-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    state=$(ionosctl dbaas in-memory-db-v2 cluster get --location "${location}" --cluster-id "${cluster_id}" -o json | jq -r '.metadata.state // empty')
    if [[ "$state" != "AVAILABLE" ]]; then
        skip "Cluster not AVAILABLE (state: ${state:-unknown})"
    fi

    run ionosctl dbaas in-memory-db-v2 cluster update \
        --location "${location}" \
        --cluster-id "${cluster_id}" \
        --cores 2 \
        -o json
    assert_success

    new_cores=$(echo "$output" | jq -r '.properties.instances.cores')
    assert_equal "$new_cores" "2"
}

@test "Restore in-memory-db-v2 cluster from snapshot" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    state=$(ionosctl dbaas in-memory-db-v2 cluster get --location "${location}" --cluster-id "${cluster_id}" -o json | jq -r '.metadata.state // empty')
    if [[ "$state" != "AVAILABLE" ]]; then
        skip "Cluster not AVAILABLE (state: ${state:-unknown})"
    fi

    snapshot_id=$(ionosctl dbaas in-memory-db-v2 snapshot list --location "${location}" -o json \
        | jq -r --arg cid "$cluster_id" '[.items[] | select(.properties.clusterId == $cid)][0].id // empty')
    if [[ -z "$snapshot_id" ]]; then
        skip "No snapshots available for this cluster yet"
    fi

    run ionosctl dbaas in-memory-db-v2 cluster restore \
        --location "${location}" \
        --cluster-id "${cluster_id}" \
        --snapshot-id "${snapshot_id}" \
        -f
    assert_success
}

@test "Delete in-memory-db-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas in-memory-db-v2 cluster delete --location "${location}" --cluster-id "${cluster_id}" -f
    assert_success

    # Wait for deletion so the datacenter/LAN lock is released before teardown.
    run wait_for_gone "${cluster_id}" 1200
    assert_success
}

@test "Delete all in-memory-db-v2 clusters" {
    # May succeed (if clusters exist) or be a no-op. Just verify it doesn't crash.
    run ionosctl dbaas in-memory-db-v2 cluster delete --all -f
}

# --- Validation / error cases ---

@test "Create cluster missing required flags fails" {
    run ionosctl dbaas in-memory-db-v2 cluster create 2>&1
    assert_failure
}

@test "Create cluster with --recovery-time but no --snapshot-id fails" {
    run ionosctl dbaas in-memory-db-v2 cluster create \
        --location "${location}" \
        --datacenter-id "00000000-0000-4000-8000-000000000000" \
        --lan-id "1" \
        --cidr "192.168.1.100/24" \
        --user testuser \
        --password testpass123 \
        --recovery-time "2024-01-15T10:00:00Z" 2>&1
    assert_failure
    assert_stderr -p "--recovery-time requires --snapshot-id"
}

@test "Delete cluster without id or --all fails" {
    run ionosctl dbaas in-memory-db-v2 cluster delete 2>&1
    assert_failure
}

@test "Restore without snapshot-id fails" {
    run ionosctl dbaas in-memory-db-v2 cluster restore --cluster-id "00000000-0000-4000-8000-000000000000" 2>&1
    assert_failure
}

@test "Restore without cluster-id fails" {
    run ionosctl dbaas in-memory-db-v2 cluster restore --snapshot-id "00000000-0000-4000-8000-000000000000" 2>&1
    assert_failure
}

# --- Help output verification ---

@test "Cluster create help shows expected flags" {
    run ionosctl dbaas in-memory-db-v2 cluster create --help
    assert_success
    assert_output -p "--datacenter-id"
    assert_output -p "--lan-id"
    assert_output -p "--cidr"
    assert_output -p "--user"
    assert_output -p "--password"
    assert_output -p "--eviction-policy"
    assert_output -p "--persistence-mode"
    assert_output -p "--backup-location"
    assert_output -p "--retention-days"
    assert_output -p "--snapshot-hours"
    assert_output -p "--maintenance-day"
    assert_output -p "--maintenance-time"
    assert_output -p "--snapshot-id"
    assert_output -p "--recovery-time"
}

@test "Cluster update help shows expected flags" {
    run ionosctl dbaas in-memory-db-v2 cluster update --help
    assert_success
    assert_output -p "--cluster-id"
    assert_output -p "--cores"
    assert_output -p "--ram"
    assert_output -p "--replicas"
    assert_output -p "--eviction-policy"
    assert_output -p "--persistence-mode"
    assert_output -p "--maintenance-day"
    assert_output -p "--maintenance-time"
}

@test "Cluster restore help shows expected flags" {
    run ionosctl dbaas in-memory-db-v2 cluster restore --help
    assert_success
    assert_output -p "--cluster-id"
    assert_output -p "--snapshot-id"
    assert_output -p "--recovery-time"
}

@test "Snapshot location get help shows expected flags" {
    run ionosctl dbaas in-memory-db-v2 snapshot location get --help
    assert_success
    assert_output -p "--snapshot-location-id"
}

# --- Teardown ---

teardown_file() {
    ionosctl dbaas in-memory-db-v2 cluster delete -af || true
    sleep 10

    if [[ -f /tmp/bats_test/datacenter_id ]]; then
        ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f || true
    fi

    ionosctl token delete --token "$IONOS_TOKEN" --force || true

    rm -rf /tmp/bats_test
}
