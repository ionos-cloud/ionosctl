#!/usr/bin/env bats

# tags: maria, mariadb

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="de/txl"

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
    export IONOS_TOKEN=$(ionosctl token generate)

    mkdir -p /tmp/bats_test
}

@test "Ensure no previous MariaDB Clusters" {
    run ionosctl db mariadb cluster list 2> /dev/null
    [ "$status" -eq 0 ] || fail "Failed to list mariadb clusters"

    if [ "$output" ]; then
        echo "Deleting clusters..."
        run ionosctl db mariadb cluster delete --all -f
        sleep 30
    fi

    echo "Waiting for clusters to be deleted..."
    retry_command bash -c "[ -z \"\$(ionosctl db mariadb cluster list)\" ]"
}

@test "Create MariaDB Cluster" {
    datacenter_id=$(find_or_create_resource \
        "ionosctl datacenter list -M 1 -F location=${location},state=available -o json 2> /dev/null | jq -r '.items[] | .id'" \
        "ionosctl datacenter create --name \"CLI-Test-$(randStr 8)\" --location ${location} -o json 2> /dev/null | jq -r '.id'")
    [ -n "$datacenter_id" ] || fail "datacenter_id is empty"
    assert_regex "$datacenter_id" "$uuid_v4_regex"

    lan_id=$(find_or_create_resource \
        "ionosctl lan list -M 1 --datacenter-id ${datacenter_id} -F public=false-o json 2> /dev/null | jq -r '.items[] | .id'" \
        "sleep 30 && ionosctl lan create --datacenter-id ${datacenter_id} --public=false -o json 2> /dev/null | jq -r '.id'")
    [ -n "$lan_id" ] || fail "lan_id is empty"

    sleep 60

    echo "Trying to create mariadb cluster in datacenter $datacenter_id"
    run ionosctl db mariadb cluster create --name "CLI-Test-$(randStr 6)" --version 10.6 --user testuser1234 \
       --password "$(randStr 12)" --datacenter-id ${datacenter_id} --lan-id 1 --cidr 192.168.1.127/24 -o json 2> /dev/null
    assert_success

    cluster_id=$(echo "$output" | jq -r '.id')
    assert_regex "$cluster_id" "$uuid_v4_regex"
    echo "created mariadb cluster $cluster_id"

    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
    echo "$cluster_id" > /tmp/bats_test/cluster_id
}

@test "Backup Listing is successful" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    # Sadly cannot really assert backups output as this is a new cluster and no backups are available.
    # This is a stateful operation and its output cannot really be tested...
    # TODO: Improve me if possible

    run ionosctl db mariadb backup list 2> /dev/null
    assert_success

    run ionosctl db mariadb backup list --cluster-id "${cluster_id}" 2> /dev/null
    assert_success
}

@test "Find MariaDB Cluster" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)
    echo "Finding mariadb cluster $cluster_id"

    run ionosctl db mariadb cluster get --cluster-id "$cluster_id" -o json 2> /dev/null
    assert_success
    cluster_name=$(echo "$output" | jq -r '.properties.displayName')
    echo "Got name $cluster_name"

    run ionosctl db mariadb cluster list -n "${cluster_name}" -M 1 --cols ClusterId --no-headers 2> /dev/null
    assert_success
    assert_output "$cluster_id"
}

@test "Assert DNS resolves to CIDR" {
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    # Extract the DNS and CIDR from the JSON output
    clusters_json=$(ionosctl db mariadb cluster get -i "${cluster_id}" -o json)
    dns_name=$(echo "$clusters_json" | jq -r '.properties.dnsName')
    cidr=$(echo "$clusters_json" | jq -r '.properties.connections[0].cidr')

    echo "Looking up DNS name: $dns_name expecting an IP in CIDR: $cidr"
    ip=$(dig +short "$dns_name")
    echo "dig result: $ip"

    run ip_in_cidr "$ip" "$cidr"
    assert_success
}

teardown_file() {
    datacenter_id=$(cat /tmp/bats_test/datacenter_id)
    cluster_id=$(cat /tmp/bats_test/cluster_id)

    echo "cleaning up datacenter $datacenter_id and mariadb cluster $cluster_id"
    retry_command run ionosctl dbaas mariadb cluster delete --cluster-id "$cluster_id" -f
    sleep 120
    retry_command run ionosctl datacenter delete --datacenter_id "$datacenter_id" -f -w -t 1200

    run ionosctl token delete --token "$IONOS_TOKEN"
    unset IONOS_TOKEN

    rm -rf /tmp/bats_test
}
