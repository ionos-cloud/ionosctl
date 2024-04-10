#!/usr/bin/env bats

# tags: k8s, k8s_cluster, k8s_nodepool, k8s_node

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="es/vit"

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

@test "Create K8s Nodepool" {
    datacenter_id=$(find_or_create_resource \
        "ionosctl datacenter list -M 1 -F location=${location},state=available -o json 2> /dev/null | jq -r '.items[] | .id'" \
        "ionosctl datacenter create --name \"CLI-Test-$(randStr 8)\" --location ${location} -o json 2> /dev/null | jq -r '.id'")
    [ -n "$datacenter_id" ] || fail "$datacenter_id is empty"
    assert_regex "$datacenter_id" "$uuid_v4_regex"

    cluster_id=$(find_or_create_resource \
        "ionosctl k8s cluster list -F public=true,state=available -M 1 -o json 2> /dev/null | jq -r '.items[] | .id'" \
        "ionosctl k8s cluster create --name \"CLI-Test-$(randStr 8)\" -o json 2> /dev/null | jq -r '.id'")
    [ -n "$cluster_id" ] || fail "$cluster_id is empty"
    assert_regex "cluster_id" "$uuid_v4_regex"

    sleep 120

    echo "Trying to create k8s nodepool in cluster $cluster_id and datacenter $datacenter_id"
    run ionosctl k8s nodepool create --name "CLI-Test-$(randStr 8)" --cluster-id "$cluster_id" --datacenter-id "$datacenter_id" -o json 2> /dev/null -W -t 600
    assert_success
    nodepool_id=$(echo "$output" | jq -r '.id')
    assert_regex "nodepool_id" "$uuid_v4_regex"
    echo "created k8s nodepool $nodepool_id"

    run ionosctl k8s node list --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --cols PublicIP --no-headers
    assert_success
    node_ip=$(echo "$output")
    assert_regex "nodepool_id" "$ip_regex"
    echo "Found IP"

    run ssh -o StrictHostKeyChecking=no "$node_ip"
    assert_output --partial 'Permission denied (publickey)'
}

teardown_file() {
    echo "cleaning up datacenter $datacenter_id and k8s resources $cluster_id ; $nodepool_id"
    retry_command run ionosctl k8s nodepool delete --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -f -w -t 1200
    sleep 30
    retry_command run ionosctl k8s cluster delete --cluster-id "$cluster_id" -f -w -t 1200
    sleep 30
    retry_command run ionosctl datacenter delete --datacenter_id "$datacenter_id" -f -w -t 1200
}
