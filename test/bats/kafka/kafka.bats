#!/usr/bin/env bats

# tags: kafka

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    mkdir -p /tmp/bats_test

    run ionosctl datacenter create --location "de/fra" --name "cli-test-$(randStr 6)" -o json 2> /dev/null
    assert_success

    run ionosctl datacenter create --name \"cli-test-$(randStr 8)\" --location de/fra -o json 2> /dev/null
    datacenter_id=$(echo "$output" | jq -r '.id')
    [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
    assert_regex "$datacenter_id" "$uuid_v4_regex"
    echo "$datacenter_id" > /tmp/bats_test/datacenter_id

    retry_until "ionosctl datacenter get --datacenter-id ${datacenter_id} -o json 2> /dev/null | jq -r '.metadata.state'" \
        "[[ \$output == \"AVAILABLE\" ]]" 10 60

    run ionosctl lan create --datacenter-id ${datacenter_id} --public=false -o json 2> /dev/null
    lan_id=$(echo "$output" | jq -r '.id')
    [ -n "$lan_id" ] || fail "lan_id is empty"

    retry_until "ionosctl lan get --datacenter-id ${datacenter_id} --lan-id ${lan_id} -o json 2> /dev/null | jq -r '.metadata.state'" \
        "[[ \$output == \"AVAILABLE\" ]]" 10 60

    echo "$lan_id" > /tmp/bats_test/lan_id
}

@test "Create Kafka Cluster" {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)

    run ionosctl kafka cluster create --name "cli-test-$(randStr 6)" --location "de/fra" --datacenter-id "${datacenter_id}" \
        --lan-id "${lan_id}" --size XS --version 3.7.0 --broker-addresses 192.168.0.100/24,192.168.0.101/24,192.168.0.102/24 -o json 2> /dev/null
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"

    retry_until "ionosctl kafka cluster get --location \"de/fra\" --cluster-id ${cluster_id} -o json 2> /dev/null | jq -r '.metadata.state'" \
        "[[ \$output == \"AVAILABLE\" ]]" 10 120

    echo "created kafka cluster $cluster_id"
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Get and list Kafka Clusters" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl kafka cluster get --location "de/fra" --cluster-id "${cluster_id}" -o json 2> /dev/null
    assert_success

    run ionosctl kafka cluster list --location "de/fra" -o json 2> /dev/null | jq -r '.items[] | select(.id == "'${cluster_id}'") | .id'
    assert_success
}

@test "Create Kafka Topic" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    run ionosctl kafka topic create --location "de/fra" --cluster-id "${cluster_id}" --name "cli-test-$(randStr 6)" -o json 2> /dev/null
    assert_success

    topic_id=$(echo "$output" | jq -r '.id')
    assert_regex "$topic_id" "$uuid_v4_regex"

    retry_until "ionosctl kafka topic get --location \"de/fra\" --cluster-id ${cluster_id} --topic-id ${topic_id} -o json 2> /dev/null | jq -r '.metadata.state'" \
        "[[ \$output == \"AVAILABLE\" ]]" 10 60

    echo "$topic_id" > /tmp/bats_test/topic_id
}

@test "Get and list Kafka Cluster Topics" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    topic_id=$(cat /tmp/bats_test/topic_id)

    run ionosctl kafka topic get --location "de/fra" --cluster-id "${cluster_id}" --topic-id "${topic_id}" -o json 2> /dev/null
    assert_success

    run ionosctl kafka topic list --location "de/fra" --cluster-id "${cluster_id}" -o json 2> /dev/null | jq -r '.items[] | select(.id == "'${topic_id}'") | .id'
    assert_success
}

# executed last in a file
teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    lan_id=$(cat /tmp/bats_test/lan_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    topic_id=$(cat /tmp/bats_test/topic_id)

    run ionosctl kafka topic delete --location "de/fra" --cluster-id "${cluster_id}" --topic-id "${topic_id}" -f
    assert_success

    run ionosctl kafka cluster delete --location "de/fra" --cluster-id "${cluster_id}" -f
    assert_success
    retry_until "ionosctl kafka cluster get --location \"de/fra\" --cluster-id ${cluster_id} -o json 2> /dev/null | jq -r '.metadata.state'" \
        "[[ \$output != \"DESTROYING\" ]]" 10 120

    run ionosctl datacenter delete --datacenter-id "${datacenter_id}" -f
    assert_success

    run ionosctl token delete --token "${IONOS_TOKEN}" -f
    unset IONOS_TOKEN
    rm -rf /tmp/bats_test
}
