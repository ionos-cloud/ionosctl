#!/usr/bin/env bats

# tags: k8s, k8s_cluster, k8s_nodepool, k8s_node

load '../libs/bats-assert/load'
load '../libs/bats-support/load'
load '../setup.bats'

setup() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

teardown() {
    echo "cleaning up datacenter $datacenter_id and k8s resources $cluster_id ; $nodepool_id"
    retry_command run ionosctl k8s nodepool delete --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -f -w -t 1200
    sleep 30
    retry_command run ionosctl k8s cluster delete --cluster-id "$cluster_id" -f -w -t 1200
    sleep 30
    retry_command run ionosctl datacenter delete --datacenter_id "$datacenter_id" -f -w -t 1200
}

@test "Create K8s cluster and datacenter" {
    run ionosctl k8s cluster create --name "$(randStr 8)" -o json -w -t 300
    assert_success
    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created k8s cluster $cluster_id"

    run ionosctl datacenter create --name "$(randStr 8)" --location us/las -o json -w -t 300
    assert_success
    datacenter_id=$(echo "$output" | jq -r '.id')
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "created datacenter $datacenter_id"
}

@test "Create public K8s Nodepool and assert that the exposed IP of the node is reachable" {
    run ionosctl k8s nodepool create --name "$(randStr 8)" --cluster-id "$cluster_id" --datacenter-id "$datacenter_id" -o json -W -t 600
    assert_success
    nodepool_id=$(echo "$output" | jq -r '.id')
    assert_regex "nodepool_id" "$uuid_v4_regex"
    echo "created k8s nodepool $nodepool_id"

    run ionosctl k8s node list --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --cols PublicIP --no-headers
    assert_success
    # TODO: Perhaps we could add a setup util which repeats the last command with '-o json' and verifies jq filtered output matches this output,
    # TODO: as well as verifying there are no columns (or other invalid things such as status reports) in the output
    node_ip=$(echo "$output")
    assert_regex "nodepool_id" "$ip_regex"
    echo "Found IP"

    run ssh -o StrictHostKeyChecking=no "$node_ip"
    assert_output --partial 'Permission denied (publickey)'
}
