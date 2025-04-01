#!/usr/bin/env bats

# tags: k8s, k8s_cluster, k8s_nodepool, k8s_node

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="es/vit"
uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'

setup_file() {
    mkdir -p /tmp/bats_test
}

@test "Create Datacenter" {
    run ionosctl datacenter create --name "CLI-Test-$(randStr 8)" --location "${location}" -o json 2> /dev/null
    assert_success
    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
}

@test "Create K8s Cluster" {
    run ionosctl k8s cluster create --name "CLI-Test-$(randStr 8)" -w -W -o json 2> /dev/null
    assert_success
    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "$cluster_id" > /tmp/bats_test/cluster_id

    sleep 60
}

@test "Create K8s Nodepool" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    [ -n "$datacenter_id" ] || fail "Datacenter ID not found"
    [ -n "$cluster_id" ] || fail "Cluster ID not found"

    run ionosctl k8s nodepool create --name "CLI-Test-$(randStr 8)" --cluster-id "$cluster_id" --datacenter-id "$datacenter_id" -W -t 600 -o json 2> /dev/null
    assert_success
    nodepool_id=$(echo "$output" | jq -r '.id')
    assert_regex "$nodepool_id" "$uuid_v4_regex"
    echo "$nodepool_id" > /tmp/bats_test/nodepool_id
}

@test "Get IP of K8s Node" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)
    [ -n "$cluster_id" ] || fail "Cluster ID not found"
    [ -n "$nodepool_id" ] || fail "Nodepool ID not found"

    run ionosctl k8s node list --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --cols PublicIP --no-headers 2> /dev/null
    assert_success
    node_ip=$(echo "$output" | tr -d '\n')
    [ -n "$node_ip" ] || fail "Node list did not return an IP address"
    assert_regex "$node_ip" "$ip_regex"
    echo "$node_ip" > /tmp/bats_test/node_ip
}

@test "IP is reachable" {
    node_ip=$(cat /tmp/bats_test/node_ip)
    [ -n "$node_ip" ] || fail "Node IP not found"

    run ssh -o StrictHostKeyChecking=no "$node_ip" exit 2> /dev/null
    assert_output --partial "Permission denied"
}

teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    echo "cleaning up datacenter $datacenter_id and k8s resources $cluster_id ; $nodepool_id"
    ionosctl k8s nodepool delete --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -f
    sleep 300
    ionosctl k8s cluster delete --cluster-id "$cluster_id" -f
    sleep 30
    ionosctl datacenter delete --datacenter_id "$datacenter_id" -f -w

    ionosctl k8s cluster delete -af

    rm -rf /tmp/bats_test
}
