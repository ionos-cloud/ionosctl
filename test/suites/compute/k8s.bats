#!/usr/bin/env bats

# paths: commands/compute/k8s/*

load '../setup.bats'

location="es/vit"
uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'

setup_file() {
    mkdir -p /tmp/bats_test
}

@test "Create Datacenter" {
    run ionosctl compute datacenter create --name "CLI-Test-$(randStr 8)" --location "${location}" -o json
    assert_success
    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
}

@test "Create K8s Cluster" {
    run ionosctl compute k8s cluster create --name "CLI-Test-$(randStr 8)" -w -o json
    assert_success
    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Create K8s Nodepool" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    [ -n "$datacenter_id" ] || fail "Datacenter ID not found"
    [ -n "$cluster_id" ] || fail "Cluster ID not found"

    run ionosctl compute k8s nodepool create --name "CLI-Test-$(randStr 8)" --cluster-id "$cluster_id" --datacenter-id "$datacenter_id" --taints "dedicated=gpu:NoSchedule,team:NoExecute" -w -o json
    assert_success
    nodepool_id=$(echo "$output" | jq -r '.id')
    assert_regex "$nodepool_id" "$uuid_v4_regex"
    echo "$nodepool_id" > /tmp/bats_test/nodepool_id

    # taints from --taints must be reflected on the created nodepool
    assert_equal "$(echo "$output" | jq -r '.properties.taints[] | select(.key=="dedicated") | .value')" "gpu"
    assert_equal "$(echo "$output" | jq -r '.properties.taints[] | select(.key=="dedicated") | .effect')" "NoSchedule"
    assert_equal "$(echo "$output" | jq -r '.properties.taints[] | select(.key=="team") | .effect')" "NoExecute"
}

@test "Update K8s Nodepool taints" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)
    [ -n "$cluster_id" ] || fail "Cluster ID not found"
    [ -n "$nodepool_id" ] || fail "Nodepool ID not found"

    # --taints overwrites the existing set
    run ionosctl compute k8s nodepool update --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --taints "dedicated=cpu:PreferNoSchedule" -w -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.taints | length')" "1"
    assert_equal "$(echo "$output" | jq -r '.properties.taints[0].key')" "dedicated"
    assert_equal "$(echo "$output" | jq -r '.properties.taints[0].value')" "cpu"
    assert_equal "$(echo "$output" | jq -r '.properties.taints[0].effect')" "PreferNoSchedule"
}

@test "Update K8s Nodepool taints rejects invalid effect" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    run ionosctl compute k8s nodepool update --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --taints "dedicated=cpu:Nope" -o json
    assert_failure
    assert_output --partial "invalid taint effect"
}

@test "Get IP of K8s Node" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)
    [ -n "$cluster_id" ] || fail "Cluster ID not found"
    [ -n "$nodepool_id" ] || fail "Nodepool ID not found"

    run ionosctl compute k8s node list --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --cols PublicIP --no-headers
    assert_success
    node_ip=$(echo "$output" | tr -d '\n')
    [ -n "$node_ip" ] || fail "Node list did not return an IP address"
    assert_regex "$node_ip" "$ip_regex"
    echo "$node_ip" > /tmp/bats_test/node_ip
}

@test "IP is reachable" {
    node_ip=$(cat /tmp/bats_test/node_ip)
    [ -n "$node_ip" ] || fail "Node IP not found"

    run ssh -o StrictHostKeyChecking=no "$node_ip" exit
    assert_stderr --partial "Permission denied"
}

teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    echo "cleaning up datacenter $datacenter_id and k8s resources $cluster_id ; $nodepool_id"
    ionosctl compute k8s nodepool delete --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -f -w
    ionosctl compute k8s cluster delete --cluster-id "$cluster_id" -f -w
    ionosctl compute datacenter delete --datacenter_id "$datacenter_id" -f -w

    ionosctl compute k8s cluster delete -af

    rm -rf /tmp/bats_test
}
