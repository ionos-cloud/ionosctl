#!/usr/bin/env bats

# tags: snapshot, volume, server, datacenter, cube, template

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    ssh-keygen -t rsa -b 4096 -N "" -f /tmp/bats_test/id_rsa

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

setup() {
    if [[ -f /tmp/bats_test/email ]] && [[ -f /tmp/bats_test/password ]]; then
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"
    fi
}

@test "Create temporary sub-user with relevant permissions" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "test-user-$(randStr 4)" --last-name "test-last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl group create --name "test-group-$(randStr 4)" \
        --create-dc --create-nic --create-backup --create-snapshot --reserve-ip \
        -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    sleep 10

    run ionosctl group user add --user-id "$(cat /tmp/bats_test/user_id)" \
        --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success
}

@test "Is temp user" {
    run ionosctl whoami
    assert_success
    assert_output "$(cat /tmp/bats_test/email)"
}


@test "Get and verify XS template" {
    run ionosctl template list -F name=XS -o json 2> /dev/null
    assert_success
    xs_output="$output"
    echo "$xs_output" | jq -r '.items[0].id' > /tmp/bats_test/template_id

    run ionosctl template get --template-id "$(cat /tmp/bats_test/template_id)" --cols Ram --no-headers
    assert_success
    assert_output "$(echo "$xs_output" | jq -r '.items[0].properties.ram')"

    run ionosctl template get --template-id "$(cat /tmp/bats_test/template_id)" --cols Cores --no-headers
    assert_success
    assert_output "$(echo "$xs_output" | jq -r '.items[0].properties.cores')"
}

@test "Create Datacenter" {
    run ionosctl datacenter create --name "snapshot-test-dc-$(randStr 8)" --location "de/txl" -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5
}

@test "Create Volume" {
    run ionosctl volume create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --size 10 --type "HDD" -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/volume_id
}

@test "Create Snapshot from Volume" {
    run ionosctl snapshot create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --volume-id "$(cat /tmp/bats_test/volume_id)" --name "snapshot-test-$(randStr 8)" -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/snapshot_id
}

@test "Create Volume from Snapshot" {
    run ionosctl volume create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --image-id "$(cat /tmp/bats_test/snapshot_id)" --size 10 --type "HDD" -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/volume_from_snapshot_id
}

@test "Create Server" {
    run ionosctl server create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --name "snapshot-test-server-$(randStr 8)" --cores 2 --ram 4096 -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/server_id
}

@test "Attach Volume to Server" {
    run ionosctl server volume attach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/volume_from_snapshot_id)" \
        -w -t 300 -o json 2> /dev/null
    assert_success
}

@test "Detach Volume from Server" {
    run ionosctl server volume detach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/volume_from_snapshot_id)" -w -t 300 -f
    assert_success
}

@test "Create CUBE Server using Snapshot" {
    run ionosctl server create --type CUBE --template-id "$(cat /tmp/bats_test/template_id)" \
        --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --image-id "$(cat /tmp/bats_test/snapshot_id)" \
        -w -t 400 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/cube_server_id
    assert_equal "$(echo "$output" | jq -r '.properties.type')" "CUBE"

    run ionosctl server get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --server-id "$(cat /tmp/bats_test/cube_server_id)" --no-headers --cols Type
    assert_success
    assert_output -p "CUBE"
}

@test "Detach and Cleanup" {
    run ionosctl server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --server-id "$(cat /tmp/bats_test/cube_server_id)" -w -t 300 -f
    assert_success

    run ionosctl volume delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
        --volume-id "$(cat /tmp/bats_test/volume_id)" -w -t 300 -f
    assert_success

    run ionosctl snapshot delete --snapshot-id "$(cat /tmp/bats_test/snapshot_id)" -f
    assert_success

    run ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -w -t 300 -f
    assert_success
}

teardown_file() {
    (
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        ionosctl datacenter delete -af
        ionosctl snapshot delete -af
    )

    ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
