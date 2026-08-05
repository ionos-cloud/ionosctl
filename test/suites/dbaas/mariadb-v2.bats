#!/usr/bin/env bats

# paths: commands/dbaas/mariadb-v2/*

load '../setup.bats'

location="de/txl"

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}


# --- Auth setup ---

@test "Generate Token" {
    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

# --- Read-only operations (no cluster needed) ---

@test "List mariadb-v2 versions" {
    run ionosctl dbaas mariadb-v2 version list --location ${location} -o json
    assert_success
}

@test "List mariadb-v2 backup locations" {
    run ionosctl dbaas mariadb-v2 backup location list --location ${location} -o json
    assert_success
}

@test "List mariadb-v2 backups" {
    run ionosctl dbaas mariadb-v2 backup list -o json
    assert_success
}

# --- Infrastructure setup ---

@test "Create Datacenter" {
    run ionosctl datacenter create --name "CLI-MariaV2-Test-$(randStr 8)" --location ${location} -w -o json
    assert_success

    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "created datacenter $datacenter_id"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
}

@test "Create LAN" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)

    run ionosctl lan create --datacenter-id ${datacenter_id} --public=false -w -o json
    assert_success

    lan_id=$(echo "$output" | jq -r '.id')
    echo "$lan_id" > /tmp/bats_test/lan_id
}

# --- Cluster lifecycle ---

@test "Create mariadb-v2 cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)

    # NOTE: intentionally no global -w here. The MariaDB v3 API returns a relative
    # href (/clusters/{id}) that omits the /v2 path segment the real endpoint needs,
    # so the generic --wait poller builds a 404 URL and would block until timeout.
    # We poll for AVAILABLE explicitly in the next test instead.
    run ionosctl dbaas mariadb-v2 cluster create \
        --location ${location} \
        --name "CLI-MariaV2-Test-$(randStr 6)" \
        --version 10.11 \
        --datacenter-id ${datacenter_id} \
        --lan-id ${lan_id} \
        --cidr 192.168.1.100/24 \
        --user testuser \
        --password "$(randStr 13)A1@" \
        --database testdb \
        --instances 1 \
        --cores 1 \
        --ram 4GB \
        --storage-size 10GB \
        --backup-location eu-central-4 \
        --backup-retention-days 7 \
        -o json
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created mariadb-v2 cluster $cluster_id"
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Wait for mariadb-v2 cluster to become AVAILABLE" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    # Poll up to ~40 minutes (mariadb provisioning is slow).
    for _ in $(seq 1 80); do
        state=$(ionosctl dbaas mariadb-v2 cluster get --cluster-id "${cluster_id}" --location ${location} -o json 2>/dev/null | jq -r '.metadata.state // empty')
        if [[ "$state" == "AVAILABLE" ]]; then
            return 0
        fi
        if [[ "$state" == "FAILED" ]]; then
            fail "cluster entered FAILED state"
        fi
        sleep 30
    done
    fail "cluster did not become AVAILABLE in time (last state: ${state:-unknown})"
}

@test "Get mariadb-v2 cluster by ID" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas mariadb-v2 cluster get --cluster-id "$cluster_id" --location ${location} -o json
    assert_success

    cluster_name=$(echo "$output" | jq -r '.properties.name')
    assert_output -p "\"name\": \"$cluster_name\""
    echo "$cluster_name" > /tmp/bats_test/cluster_name
}

@test "List mariadb-v2 clusters" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    cluster_name=$(cat /tmp/bats_test/cluster_name)

    # JSON output
    run ionosctl dbaas mariadb-v2 cluster list -o json
    assert_success
    assert_output -p "\"name\": \"$cluster_name\""

    # Column output
    run ionosctl dbaas mariadb-v2 cluster list --cols ClusterId --no-headers
    assert_success
    assert_output -p "$cluster_id"
}

@test "Update mariadb-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    state=$(ionosctl dbaas mariadb-v2 cluster get --cluster-id "${cluster_id}" --location ${location} -o json | jq -r '.metadata.state // empty')
    if [[ "$state" != "AVAILABLE" ]]; then
        skip "Cluster not AVAILABLE (state: ${state:-unknown})"
    fi

    # The API never returns the password on GET, so a PUT must re-supply --password.
    # PUT returns the new desired properties immediately (state goes UPDATING), so
    # no --wait is needed to assert the change took.
    run ionosctl dbaas mariadb-v2 cluster update \
        --cluster-id "${cluster_id}" \
        --location ${location} \
        --cores 2 \
        --password "$(randStr 13)A1@" \
        -o json
    assert_success

    new_cores=$(echo "$output" | jq -r '.properties.instances.cores')
    assert_equal "$new_cores" "2"
}

@test "List mariadb-v2 backups for cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas mariadb-v2 backup list --cluster-id "${cluster_id}" --location ${location} -o json
    assert_success

    backup_id=$(echo "$output" | jq -r '.items[0].id // empty')
    if [[ -n "$backup_id" ]]; then
        echo "$backup_id" > /tmp/bats_test/backup_id
    fi
}

@test "Restore mariadb-v2 cluster in place" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    if [[ ! -f /tmp/bats_test/backup_id ]]; then
        skip "No backups available yet for this cluster"
    fi

    # The prior update left the cluster UPDATING; a PUT-based restore needs it back
    # to AVAILABLE first (the command itself guards on this). Wait, bounded.
    for _ in $(seq 1 40); do
        state=$(ionosctl dbaas mariadb-v2 cluster get --cluster-id "${cluster_id}" --location ${location} -o json 2>/dev/null | jq -r '.metadata.state // empty')
        [[ "$state" == "AVAILABLE" ]] && break
        sleep 30
    done
    if [[ "$state" != "AVAILABLE" ]]; then
        skip "Cluster not AVAILABLE after update (state: ${state:-unknown})"
    fi

    # In-place point-in-time restore to the latest recoverable point.
    run ionosctl dbaas mariadb-v2 cluster restore \
        --cluster-id "${cluster_id}" \
        --location ${location} \
        --recovery-time now \
        --password "$(randStr 13)A1@" \
        -f
    assert_success
}

@test "Delete mariadb-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas mariadb-v2 cluster delete --cluster-id "${cluster_id}" --location ${location} -f
    assert_success
}

# --- Validation / error cases ---

@test "Create cluster missing required flags fails" {
    run ionosctl dbaas mariadb-v2 cluster create 2>&1
    assert_failure
}

@test "Create with out-of-range backup-retention-days fails" {
    run ionosctl dbaas mariadb-v2 cluster create \
        --location ${location} \
        --datacenter-id "00000000-0000-4000-8000-000000000000" \
        --lan-id 1 --cidr 192.168.1.100/24 \
        --user testuser --password "passw0rd123" --database testdb \
        --version 10.11 --backup-retention-days 999
    assert_failure
    assert_stderr -p "backup-retention-days must be between 1 and 365"
}

@test "Create with --recovery-time but no --backup-id fails" {
    run ionosctl dbaas mariadb-v2 cluster create \
        --location ${location} \
        --datacenter-id "00000000-0000-4000-8000-000000000000" \
        --lan-id 1 --cidr 192.168.1.100/24 \
        --user testuser --password "passw0rd123" --database testdb \
        --version 10.11 --recovery-time now
    assert_failure
    assert_stderr -p "--recovery-time requires --backup-id"
}

@test "Delete cluster without id or --all fails" {
    run ionosctl dbaas mariadb-v2 cluster delete 2>&1
    assert_failure
}

@test "Delete with --name but no --all fails" {
    run ionosctl dbaas mariadb-v2 cluster delete --name foo 2>&1
    assert_failure
}

@test "Update without --password fails" {
    run ionosctl dbaas mariadb-v2 cluster update --cluster-id "00000000-0000-4000-8000-000000000000" --location ${location} 2>&1
    assert_failure
}

# --- Help output verification ---

@test "Cluster create help shows expected flags" {
    run ionosctl dbaas mariadb-v2 cluster create --help
    assert_success
    assert_output -p "--datacenter-id"
    assert_output -p "--lan-id"
    assert_output -p "--cidr"
    assert_output -p "--user"
    assert_output -p "--password"
    assert_output -p "--database"
    assert_output -p "--backup-location"
    assert_output -p "--backup-retention-days"
    assert_output -p "--storage-size"
    assert_output -p "--maintenance-day"
    assert_output -p "--maintenance-time"
    assert_output -p "--backup-id"
    assert_output -p "--recovery-time"
}

@test "Cluster update help shows expected flags" {
    run ionosctl dbaas mariadb-v2 cluster update --help
    assert_success
    assert_output -p "--cluster-id"
    assert_output -p "--password"
    assert_output -p "--cores"
    assert_output -p "--ram"
    assert_output -p "--instances"
    assert_output -p "--backup-retention-days"
    assert_output -p "--maintenance-day"
    assert_output -p "--maintenance-time"
}

@test "Cluster restore help shows expected flags" {
    run ionosctl dbaas mariadb-v2 cluster restore --help
    assert_success
    assert_output -p "--cluster-id"
    assert_output -p "--recovery-time"
    assert_output -p "--password"
}

# --- Teardown ---

teardown_file() {
    # Delete any leftover cluster (the lifecycle test may have already deleted it).
    if [[ -f /tmp/bats_test/cluster_id ]]; then
        ionosctl dbaas mariadb-v2 cluster delete --cluster-id "$(cat /tmp/bats_test/cluster_id)" --location de/txl -f || true
    fi
    ionosctl dbaas mariadb-v2 cluster delete -af || true

    # The cluster's private LAN is delete-protected until the cluster is fully gone,
    # which blocks the datacenter delete. Wait for the cluster to disappear first.
    # No --wait here: the MariaDB v3 relative-href issue makes the generic -w poller
    # unreliable for this API, so we poll cluster listing directly.
    if [[ -f /tmp/bats_test/cluster_id ]]; then
        cluster_id=$(cat /tmp/bats_test/cluster_id)
        for _ in $(seq 1 40); do
            if ! ionosctl dbaas mariadb-v2 cluster list --location de/txl --cols ClusterId --no-headers 2>/dev/null | grep -q "$cluster_id"; then
                break
            fi
            sleep 30
        done
    fi

    if [[ -f /tmp/bats_test/datacenter_id ]]; then
        ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f -w || true
    fi

    if [[ -f /tmp/bats_test/token ]]; then
        ionosctl token delete --token "$(cat /tmp/bats_test/token)" -f || true
    fi

    rm -rf /tmp/bats_test
}
