#!/usr/bin/env bats

# tags: server, console, template, gpu

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
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

@test "Create temporary user with relevant permissions" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl compute user create --first-name "random-$(randStr 4)" --last-name "last-$(randStr 4)" \
     --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl compute group create --name "test-server-$(randStr 4)" \
     --create-dc --create-nic --reserve-ip \
     -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    sleep 10

    run ionosctl compute group user add --user-id "$(cat /tmp/bats_test/user_id)" \
     --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success

    (
        unset IONOS_TOKEN
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        run ionosctl token generate --ttl 1h
        assert_success
        echo "$output" > /tmp/bats_test/token
    )
}

@test "Is temp user" {
    run ionosctl whoami
    assert_success
    assert_output "$(cat /tmp/bats_test/email)"
}

@test "Create Datacenter" {
    run ionosctl compute datacenter create --name "server-test-$(randStr 8)" --location "es/vit" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5
}

@test "Create ENTERPRISE Server" {
    run ionosctl compute server create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --cores 1 --ram 1GB -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/server_id
}

@test "Server get --cols ServerId,Name,Cores,Ram" {
    run ionosctl compute server get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --cols ServerId,Name,Cores,Ram --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/server_id)"
}

@test "Server list --datacenter-id, verify in output" {
    run ionosctl compute server list --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --no-headers
    assert_success
    assert_output -p "$(cat /tmp/bats_test/server_id)"
}

@test "Get and verify XS template" {
    run ionosctl compute template list -F name=XS -o json 2> /dev/null
    assert_success
    xs_output="$output"
    echo "$xs_output" | jq -r '.items[0].id' > /tmp/bats_test/template_id

    run ionosctl compute template get --template-id "$(cat /tmp/bats_test/template_id)" --cols Ram --no-headers
    assert_success
    assert_output "$(echo "$xs_output" | jq -r '.items[0].properties.ram')"

    run ionosctl compute template get --template-id "$(cat /tmp/bats_test/template_id)" --cols Cores --no-headers
    assert_success
    assert_output "$(echo "$xs_output" | jq -r '.items[0].properties.cores')"
}

@test "Create CUBE Server with Direct Attached Storage" {
    # Find a suitable HDD image for the CUBE server
    run ionosctl compute image list -F imageAliases=ubuntu:latest -F location="es/vit" -F imageType=hdd --cols ImageId --no-headers
    assert_success
    echo "$output" | head -n 1 > /tmp/bats_test/hdd_image_id

    run ionosctl compute server create --name "bats-test-$(randStr 8)" --type "CUBE" \
     -k /tmp/bats_test/id_rsa.pub --template-id "$(cat /tmp/bats_test/template_id)" \
     --image-id "$(cat /tmp/bats_test/hdd_image_id)" --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     -w -t 400 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/cube_server_id
    assert_equal "$(echo "$output" | jq -r '.properties.type')" "CUBE"

    run ionosctl compute server get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/cube_server_id)" --no-headers --cols Type
    assert_success
    assert_output -p "CUBE"
}

@test "Create de/fra/2 Datacenter for GPU Server" {
    run ionosctl compute datacenter create --name "gpu-test-$(randStr 8)" --location "de/fra/2" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id_gpu
    sleep 5
}

@test "Create GPU Server" {
    run ionosctl compute server create --name "bats-gpu-test-$(randStr 8)" --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" \
     --type "GPU" --template-id "6913ed82-a143-4c15-89ac-08fb375a97c5" -w -t 600 -o json 2> /dev/null
    assert_success
    assert_output -p "GPU"
    echo "$output" | jq -r '.id' > /tmp/bats_test/gpu_server_id
}

@test "List GPUs for Server" {
    run ionosctl compute server gpu list --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" \
     --server-id "$(cat /tmp/bats_test/gpu_server_id)" -o json 2> /dev/null
    assert_success

    if [ "$(echo "$output" | jq -r '.items | length')" -gt 0 ]; then
        echo "$output" | jq -r '.items[0].id' > /tmp/bats_test/gpu_id
    else
        fail "No GPUs found for the server"
    fi
}

@test "Get GPU by ID" {
    run ionosctl compute server gpu get --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" \
     --server-id "$(cat /tmp/bats_test/gpu_server_id)" --gpu-id "$(cat /tmp/bats_test/gpu_id)" -o json 2> /dev/null
    assert_success
}

@test "Server Console is accessible. Token is valid." {
    run ionosctl compute server token get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --no-headers
    assert_success

    run ionosctl compute server console get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --no-headers
    assert_success

    run curl "$(echo "$output" | grep -o 'https://[^ ]*')"
    assert_success
    assert_output --partial "<title>Remote Console</title>"
}

@test "Delete GPU Server and Datacenter" {
    run ionosctl compute server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" \
     --server-id "$(cat /tmp/bats_test/gpu_server_id)" -f -w -t 600
    assert_success

    run ionosctl compute datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" -f -w -t 600
    assert_success
}

@test "Delete CUBE Server" {
    run ionosctl compute server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/cube_server_id)" -f -w -t 600
    assert_success
}

@test "Delete ENTERPRISE Server" {
    run ionosctl compute server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" -w -t 600 -f
    assert_success
}

@test "Delete Datacenter" {
    run ionosctl compute datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f -w -t 600
    assert_success
}

teardown_file() {
    (
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
        ionosctl compute datacenter delete -af
    )

    ionosctl compute user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl compute group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
