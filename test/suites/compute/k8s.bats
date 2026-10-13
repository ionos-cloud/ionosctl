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

    run ionosctl compute k8s nodepool create --name "CLI-Test-$(randStr 8)" --cluster-id "$cluster_id" --datacenter-id "$datacenter_id" \
        --taints "dedicated=gpu:NoSchedule,spot:PreferNoSchedule" -w -o json
    assert_success
    nodepool_id=$(echo "$output" | jq -r '.id')
    assert_regex "$nodepool_id" "$uuid_v4_regex"
    echo "$nodepool_id" > /tmp/bats_test/nodepool_id

    assert_equal "$(echo "$output" | jq -c '[.properties.taints[] | [.key, .value, .effect]]')" \
        '[["dedicated","gpu","NoSchedule"],["spot",null,"PreferNoSchedule"]]'
}

@test "K8s Nodepool get shows Taints column" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    run ionosctl compute k8s nodepool get --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --cols Taints --no-headers
    assert_success
    assert_output "dedicated=gpu:NoSchedule,spot:PreferNoSchedule"
}

@test "K8s Nodepool update without --taints keeps existing taints" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    run ionosctl compute k8s nodepool update --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --labels bats=taints -w -o json
    assert_success

    run ionosctl compute k8s nodepool get --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -o json
    assert_success
    assert_equal "$(echo "$output" | jq -r '.properties.labels.bats')" "taints"
    assert_equal "$(echo "$output" | jq -c '[.properties.taints[] | .key] | sort')" '["dedicated","spot"]'
}

@test "K8s Nodepool update --taints replaces taints" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    run ionosctl compute k8s nodepool update --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" \
        --taints "example.com/team=bats:NoExecute" -w -o json
    assert_success

    run ionosctl compute k8s nodepool get --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -o json
    assert_success
    assert_equal "$(echo "$output" | jq -c '[.properties.taints[] | [.key, .value, .effect]]')" \
        '[["example.com/team","bats","NoExecute"]]'
}

@test "K8s Nodes carry the node pool taints" {
    command -v kubectl >/dev/null 2>&1 || skip "kubectl not found"
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl compute k8s kubeconfig get --cluster-id "$cluster_id"
    assert_success
    echo "$output" > /tmp/bats_test/kubeconfig

    # Node taints can lag behind the node pool update
    for _ in $(seq 1 30); do
        run kubectl --kubeconfig /tmp/bats_test/kubeconfig get nodes \
            -o jsonpath='{range .items[*].spec.taints[*]}{.key}={.value}:{.effect}{"\n"}{end}'
        [[ "$output" == *"example.com/team=bats:NoExecute"* ]] && break
        sleep 10
    done
    assert_success
    assert_output --partial "example.com/team=bats:NoExecute"
    refute_output --partial "dedicated=gpu"
}

@test "K8s Nodepool update --taints \"\" removes all taints" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    run ionosctl compute k8s nodepool update --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --taints "" -w -o json
    assert_success

    run ionosctl compute k8s nodepool get --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -o json
    assert_success
    assert_equal "$(echo "$output" | jq -c '.properties.taints // [] | length')" "0"
}

@test "K8s Nodepool rejects an invalid taint before calling the API" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    run ionosctl compute k8s nodepool update --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --taints "dedicated=gpu:Never"
    assert_failure
    assert_stderr --partial "effect must be one of NoSchedule, NoExecute, PreferNoSchedule"

    run ionosctl compute k8s nodepool update --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" --taints "dedicated=gpu"
    assert_failure
    assert_stderr --partial "expected KEY[=VALUE]:EFFECT"
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

# Polls a get command until it reports 404. Other errors (5xx, timeouts) keep polling.
wait_until_not_found() {
    local out
    for _ in $(seq 1 60); do
        if ! out=$("$@" 2>&1) && [[ "$out" == *"404 Not Found"* ]]; then
            return 0
        fi
        sleep 20
    done
    echo "gave up waiting for deletion: $*" >&2
    return 1
}

teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    nodepool_id=$(cat /tmp/bats_test/nodepool_id)

    echo "cleaning up datacenter $datacenter_id and k8s resources $cluster_id ; $nodepool_id"
    ionosctl compute k8s nodepool delete --cluster-id "$cluster_id" --nodepool-id "$nodepool_id" -f -w
    # nodepool delete -w returns while the pool is still DESTROYING; cluster and DC deletes 409 until it is gone
    wait_until_not_found ionosctl compute k8s nodepool get --cluster-id "$cluster_id" --nodepool-id "$nodepool_id"
    ionosctl compute k8s cluster delete --cluster-id "$cluster_id" -f -w
    wait_until_not_found ionosctl compute k8s cluster get --cluster-id "$cluster_id"
    ionosctl compute datacenter delete --datacenter-id "$datacenter_id" -f -w

    ionosctl compute k8s cluster delete -af

    rm -rf /tmp/bats_test
}
