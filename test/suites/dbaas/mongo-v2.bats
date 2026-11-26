#!/usr/bin/env bats

# paths: commands/dbaas/mongo-v2/*

load '../setup.bats'

location="gb/bhx"

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    export uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

cluster_state() {
    ionosctl dbaas mongo-v2 cluster get --cluster-id "$1" --location ${location} --cols State --no-headers 2>&1 | tr -d ' '
}

wait_for_cluster_available() {
    for _ in $(seq 1 "${2:-40}"); do
        [[ "$(cluster_state "$1")" == "AVAILABLE" ]] && return 0
        sleep 15
    done
    return 1
}

wait_for_cluster_gone() {
    for _ in $(seq 1 40); do
        ionosctl dbaas mongo-v2 cluster get --cluster-id "$1" --location ${location} 2>&1 | grep -q "404" && return 0
        sleep 15
    done
    return 1
}

# --- Auth setup ---

@test "Generate Token" {
    run ionosctl token generate --ttl 2h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

# --- Read-only operations (no cluster needed) ---

@test "List mongo-v2 versions" {
    run ionosctl dbaas mongo-v2 version list --location ${location} -o json
    assert_success
    assert_output -p "7.0"
}

@test "List and get mongo-v2 BUSINESS templates" {
    run ionosctl dbaas mongo-v2 template list --location ${location} --edition BUSINESS -o json
    assert_success

    template_id=$(echo "$output" | jq -r '.items | sort_by(.properties.cores, .properties.ram) | .[0].id')
    assert_regex "$template_id" "$uuid_v4_regex"
    echo "$template_id" > /tmp/bats_test/template_id

    run ionosctl dbaas mongo-v2 template get --template-id "$template_id" -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.edition')" "BUSINESS"
}

@test "List mongo-v2 snapshot locations" {
    run ionosctl dbaas mongo-v2 snapshot location list --location ${location} -o json
    assert_success
    assert_output -p "eu-central-3"
}

@test "List mongo-v2 snapshots" {
    run ionosctl dbaas mongo-v2 snapshot list -o json
    assert_success
}

# --- Infrastructure setup ---

@test "Create Datacenter" {
    run ionosctl datacenter create --name "CLI-MongoV2-Test-$(randStr 8)" --location ${location} -w -o json
    assert_success

    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
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

# Created without -w: it provisions while the ENTERPRISE tests run.
@test "Create mongo-v2 BUSINESS cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)
    business_name="CLI-MongoV2-Test-$(randStr 6)"
    echo "$business_name" > /tmp/bats_test/business_cluster_name

    run ionosctl dbaas mongo-v2 cluster create \
        --location ${location} \
        --edition BUSINESS \
        --template-id "$(cat /tmp/bats_test/template_id)" \
        --instances 1 \
        --name "$business_name" \
        --version 7.0 \
        --datacenter-id ${datacenter_id} \
        --lan-id ${lan_id} \
        --cidr 192.168.1.110/24 \
        --user testuser \
        --password "$(randStr 14)Aa1" \
        --database testdb \
        -o json
    assert_success
    assert_output_not_empty

    business_id=$(echo "$output" | jq -r '.id')
    assert_regex "$business_id" "$uuid_v4_regex"
    echo "$business_id" > /tmp/bats_test/business_cluster_id
}

@test "Create mongo-v2 ENTERPRISE cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)
    cluster_name="CLI-MongoV2-Test-$(randStr 6)"
    echo "$cluster_name" > /tmp/bats_test/cluster_name

    run ionosctl dbaas mongo-v2 cluster create \
        --location ${location} \
        --edition ENTERPRISE \
        --name "$cluster_name" \
        --version 7.0 \
        --datacenter-id ${datacenter_id} \
        --lan-id ${lan_id} \
        --cidr 192.168.1.100/24 \
        --user testuser \
        --password "$(randStr 14)Aa1" \
        --database testdb \
        --instances 1 \
        --cores 1 \
        --ram 8GB \
        --storage-size 10GB \
        --backup-location eu-central-3 \
        --retention-days 7 \
        -w --timeout 2400 -o json
    assert_success
    assert_output_not_empty

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Get mongo-v2 cluster by ID" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas mongo-v2 cluster get --cluster-id "$cluster_id" --location ${location} -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.name')" "$(cat /tmp/bats_test/cluster_name)"
    assert_equal "$(echo "$output" | jq -r '.properties.edition')" "ENTERPRISE"
    assert_equal "$(echo "$output" | jq -r '.metadata.state')" "AVAILABLE"
}

@test "List mongo-v2 clusters" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    cluster_name=$(cat /tmp/bats_test/cluster_name)

    run ionosctl dbaas mongo-v2 cluster list -o json
    assert_success
    assert_output -p "\"name\": \"$cluster_name\""

    run ionosctl dbaas mongo-v2 cluster list --location ${location} --cols ClusterId --no-headers
    assert_success
    assert_output -p "$cluster_id"

    run ionosctl dbaas mongo-v2 cluster list --location ${location} --name "$cluster_name" --cols ClusterId --no-headers
    assert_success
    assert_output -p "$cluster_id"

    business_id=$(cat /tmp/bats_test/business_cluster_id)

    run ionosctl dbaas mongo-v2 cluster list --location ${location} --edition ENTERPRISE --cols ClusterId --no-headers
    assert_success
    assert_output -p "$cluster_id"
    refute_output -p "$business_id"

    run ionosctl dbaas mongo-v2 cluster list --location ${location} --edition BUSINESS --cols ClusterId --no-headers
    assert_success
    assert_output -p "$business_id"
    refute_output -p "$cluster_id"
}

@test "Update mongo-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    new_name="CLI-MongoV2-Test-$(randStr 6)"

    run ionosctl dbaas mongo-v2 cluster update \
        --cluster-id "$cluster_id" \
        --location ${location} \
        --name "$new_name" \
        --description "updated by bats" \
        --retention-days 15 \
        -w --timeout 1200 -o json
    assert_success
    assert_output_not_empty
    echo "$new_name" > /tmp/bats_test/cluster_name

    assert_equal "$(echo "$output" | jq -r '.properties.name')" "$new_name"
    assert_equal "$(echo "$output" | jq -r '.properties.snapshot.retentionDays')" "15"
}

@test "List mongo-v2 snapshots for cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    for _ in $(seq 1 20); do
        run ionosctl dbaas mongo-v2 snapshot list --cluster-id "$cluster_id" --location ${location} -o json
        assert_success
        snapshot_id=$(echo "$output" | jq -r '.items[0].id // empty')
        [[ -n "$snapshot_id" ]] && break
        sleep 15
    done

    if [[ -n "$snapshot_id" ]]; then
        echo "$snapshot_id" > /tmp/bats_test/snapshot_id
    fi
}

@test "Get mongo-v2 snapshot" {
    if [[ ! -f /tmp/bats_test/snapshot_id ]]; then
        skip "No snapshot available yet for this cluster"
    fi

    run ionosctl dbaas mongo-v2 snapshot get --snapshot-id "$(cat /tmp/bats_test/snapshot_id)" --location ${location} -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.clusterId')" "$(cat /tmp/bats_test/cluster_id)"
}

@test "Restore mongo-v2 cluster in place" {
    if [[ ! -f /tmp/bats_test/snapshot_id ]]; then
        skip "No snapshot available yet for this cluster"
    fi
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas mongo-v2 cluster restore --cluster-id "$cluster_id" --location ${location} --recovery-time now -f
    assert_success
    assert_output -p "successfully restored"

    sleep 10
    run wait_for_cluster_available "$cluster_id"
    assert_success
}

@test "Delete mongo-v2 cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl dbaas mongo-v2 cluster delete --cluster-id "$cluster_id" --location ${location} -f
    assert_success

    run wait_for_cluster_gone "$cluster_id"
    assert_success
    rm -f /tmp/bats_test/cluster_id /tmp/bats_test/cluster_name
}

@test "Wait for mongo-v2 BUSINESS cluster" {
    run wait_for_cluster_available "$(cat /tmp/bats_test/business_cluster_id)" 80
    assert_success
}

@test "Get mongo-v2 BUSINESS cluster by ID" {
    business_id=$(cat /tmp/bats_test/business_cluster_id)

    run ionosctl dbaas mongo-v2 cluster get --cluster-id "$business_id" --location ${location} -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.name')" "$(cat /tmp/bats_test/business_cluster_name)"
    assert_equal "$(echo "$output" | jq -r '.properties.edition')" "BUSINESS"
    assert_equal "$(echo "$output" | jq -r '.properties.instanceTemplateId')" "$(cat /tmp/bats_test/template_id)"
    assert_equal "$(echo "$output" | jq -r '.properties.instanceCount')" "1"
}

@test "Update mongo-v2 BUSINESS cluster" {
    business_id=$(cat /tmp/bats_test/business_cluster_id)
    new_name="CLI-MongoV2-Test-$(randStr 6)"

    run ionosctl dbaas mongo-v2 cluster update \
        --cluster-id "$business_id" \
        --location ${location} \
        --name "$new_name" \
        --description "updated by bats" \
        --retention-days 15 \
        -w --timeout 1200 -o json
    assert_success
    assert_output_not_empty
    echo "$new_name" > /tmp/bats_test/business_cluster_name

    assert_equal "$(echo "$output" | jq -r '.properties.name')" "$new_name"
    assert_equal "$(echo "$output" | jq -r '.properties.snapshot.retentionDays')" "15"
}

@test "BUSINESS cluster refuses ENTERPRISE-only flags" {
    run ionosctl dbaas mongo-v2 cluster update --cluster-id "$(cat /tmp/bats_test/business_cluster_id)" --location ${location} --cores 2
    assert_failure
    assert_stderr -p "cannot be used on a BUSINESS cluster"
}

@test "Delete mongo-v2 BUSINESS cluster" {
    business_id=$(cat /tmp/bats_test/business_cluster_id)

    run ionosctl dbaas mongo-v2 cluster delete --cluster-id "$business_id" --location ${location} -f
    assert_success

    run wait_for_cluster_gone "$business_id"
    assert_success
    rm -f /tmp/bats_test/business_cluster_id /tmp/bats_test/business_cluster_name
}

# --- Validation / error cases ---

base_create_flags=(
    --location ${location}
    --datacenter-id "00000000-0000-4000-8000-000000000000"
    --lan-id 1 --cidr 192.168.1.100/24
    --user testuser --password "passw0rd123" --database testdb
)

@test "Create cluster missing required flags fails" {
    run ionosctl dbaas mongo-v2 cluster create 2>&1
    assert_failure
}

@test "Create BUSINESS cluster without --template-id fails" {
    run ionosctl dbaas mongo-v2 cluster create "${base_create_flags[@]}" --edition BUSINESS
    assert_failure
    assert_stderr -p "--template-id is required for the BUSINESS edition"
}

@test "Create ENTERPRISE cluster with --template-id fails" {
    run ionosctl dbaas mongo-v2 cluster create "${base_create_flags[@]}" --edition ENTERPRISE \
        --template-id "00000000-0000-4000-8000-000000000000"
    assert_failure
    assert_stderr -p "--template-id is only valid for the BUSINESS edition"
}

@test "Create sharded cluster with a single replica per shard fails" {
    run ionosctl dbaas mongo-v2 cluster create \
        --location ${location} \
        --datacenter-id "00000000-0000-4000-8000-000000000000" \
        --lan-id 1 --cidr 192.168.1.100/24,192.168.1.101/24,192.168.1.102/24 \
        --user testuser --password "passw0rd123" --database testdb \
        --edition ENTERPRISE --shards 2
    assert_failure
    assert_stderr -p "--shards needs more than one replica per shard"
}

@test "Create with fewer --cidr addresses than instances fails" {
    run ionosctl dbaas mongo-v2 cluster create "${base_create_flags[@]}" --edition ENTERPRISE --instances 3
    assert_failure
    assert_stderr -p "--cidr takes one address per instance"
}

@test "Create with out-of-range --retention-days fails" {
    run ionosctl dbaas mongo-v2 cluster create "${base_create_flags[@]}" --edition ENTERPRISE --retention-days 999
    assert_failure
    assert_stderr -p "--retention-days must be one of"
}

@test "Create with --recovery-time but no --snapshot-id fails" {
    run ionosctl dbaas mongo-v2 cluster create "${base_create_flags[@]}" --edition ENTERPRISE --recovery-time now
    assert_failure
    assert_stderr -p "--recovery-time requires --snapshot-id"
}

@test "Update with --template-id fails" {
    run ionosctl dbaas mongo-v2 cluster update --cluster-id "00000000-0000-4000-8000-000000000000" --location ${location} \
        --template-id "00000000-0000-4000-8000-000000000000"
    assert_failure
    assert_stderr -p "unknown flag: --template-id"
}

@test "Update --user without --password fails" {
    run ionosctl dbaas mongo-v2 cluster update --cluster-id "00000000-0000-4000-8000-000000000000" --location ${location} \
        --user newuser
    assert_failure
    assert_stderr -p "--password is required with --user"
}

@test "Restore with an unparseable --recovery-time fails" {
    run ionosctl dbaas mongo-v2 cluster restore --cluster-id "00000000-0000-4000-8000-000000000000" --location ${location} \
        --recovery-time yesterday -f
    assert_failure
    assert_stderr -p "could not parse --recovery-time"
}

@test "Delete cluster without id or --all fails" {
    run ionosctl dbaas mongo-v2 cluster delete 2>&1
    assert_failure
}

@test "Delete with --name but no --all fails" {
    run ionosctl dbaas mongo-v2 cluster delete --name foo 2>&1
    assert_failure
}

# --- Help output verification ---

@test "Cluster create help shows expected flags" {
    run ionosctl dbaas mongo-v2 cluster create --help
    assert_success
    for flag in --edition --template-id --instances --cores --ram --storage-size --shards \
        --datacenter-id --lan-id --cidr --user --password --database --backup-location --retention-days \
        --maintenance-day --maintenance-time --snapshot-id --recovery-time; do
        assert_output -p -- "$flag"
    done
}

@test "Cluster update help shows expected flags" {
    run ionosctl dbaas mongo-v2 cluster update --help
    assert_success
    for flag in --cluster-id --instances --cores --ram --storage-size --shards --version --retention-days \
        --maintenance-day --maintenance-time --password; do
        assert_output -p -- "$flag"
    done
    refute_output -p -- "--template-id"
    refute_output -p -- "--backup-location"
}

@test "Cluster restore help shows expected flags" {
    run ionosctl dbaas mongo-v2 cluster restore --help
    assert_success
    assert_output -p -- "--cluster-id"
    assert_output -p -- "--recovery-time"
}

# --- Teardown ---

teardown_file() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi

    for prefix in cluster business_cluster; do
        if [[ -f /tmp/bats_test/${prefix}_id ]]; then
            ionosctl dbaas mongo-v2 cluster delete --cluster-id "$(cat /tmp/bats_test/${prefix}_id)" --location ${location} -f || true
        fi
        if [[ -f /tmp/bats_test/${prefix}_name ]]; then
            ionosctl dbaas mongo-v2 cluster delete --all --name "$(cat /tmp/bats_test/${prefix}_name)" --location ${location} -f || true
        fi
    done
    for prefix in cluster business_cluster; do
        if [[ -f /tmp/bats_test/${prefix}_id ]]; then
            wait_for_cluster_gone "$(cat /tmp/bats_test/${prefix}_id)" || true
        fi
    done

    if [[ -f /tmp/bats_test/datacenter_id ]]; then
        dc=$(cat /tmp/bats_test/datacenter_id)
        for _ in $(seq 1 40); do
            out=$(ionosctl datacenter delete --datacenter-id "$dc" -f 2>&1)
            if ! echo "$out" | grep -qiE "protected|403"; then
                break
            fi
            sleep 30
        done
    fi

    if [[ -f /tmp/bats_test/token ]]; then
        ionosctl token delete --token "$(cat /tmp/bats_test/token)" -f || true
    fi

    rm -rf /tmp/bats_test
}
