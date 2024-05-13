#!/usr/bin/env bats

# tags: server, volume, image, console, nic, lan, ipblock

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    mkdir -p /tmp/bats_test

    ssh-keygen -t rsa -b 4096 -N "" -f /tmp/bats_test/id_rsa

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

# Creating a temp user allows us to simply delete all datacenters if something fails
# Prevents deletion of unrelated resources
@test "Create temporary user with relevant permissions" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "random-$(randStr 4)" --last-name "last-$(randStr 4)" \
     --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl group create --name "test-volumes-$(randStr 4)" \
     --create-dc --create-nic --reserve-ip -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    run ionosctl group user add --user-id "$(cat /tmp/bats_test/user_id)" \
     --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success
}

@test "Create Datacenter" {
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    run ionosctl datacenter create --name "volumes-test-$(randStr 8)" --location "es/vit" -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5
}

@test "Create Server" {
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    # CPU-Family should be selected correctly by default
    run ionosctl server create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --cores 1 --ram 4GB -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/server_id
}

@test "Make servers accessible via NIC and attach IP. Verify LAN is created" {
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    run ionosctl lan create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --public -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/lan_id

    run ionosctl ipblock create --location "es/vit" --size 1 --name "bats-test-$(randStr 8)" -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.properties.ips[0]' > /tmp/bats_test/ip

    run ionosctl nic create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" --name "bats-test-$(randStr 8)" --ips "$(cat /tmp/bats_test/ip)" -w -o json 2> /dev/null
    assert_success
    sleep 5
}

@test "Creating a nic with a non-existent LAN ID will create a LAN" {
    skip "todo"

    # A LAN is created by default
    run ionosctl lan list --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --no-headers --cols LanId
    assert_success
    assert_output "1"
}

@test "Attach a volume with an HDD image" {
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    # Find a suitable image
    run ionosctl image list -F imageAliases=ubuntu:20 -F location="es/vit" -F imageType=hdd --cols ImageId --no-headers
    assert_success
    echo "$output" | head -n 1 > /tmp/bats_test/image_id

    run ionosctl volume create --type "SSD Premium" --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --name "bats-test-$(randStr 8)" --size 50 --image-id "$(cat /tmp/bats_test/image_id)" \
     --ssh-key-paths /tmp/bats_test/id_rsa.pub -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/volume_id

    run ionosctl server volume attach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/volume_id)" -w
    assert_success
}

@test "Server Console is accessible" {
    run ionosctl server console get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --no-headers
    assert_success
    # Use curl to fetch the HTML content of the URL
    run curl "$output"
    assert_success
    assert_output --partial "<title>Remote Console</title>"
}

@test "ssh into the server" {
    run ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -i /tmp/bats_test/id_rsa \
     root@"$(cat /tmp/bats_test/ip)" "echo 'SSH into the server successful'"
    assert_success
}

@test "Delete Server" {
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    run ionosctl server delete \
     --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" -f
    assert_success
}

@test "Delete Datacenter" {
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    run ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f
    assert_success
}

teardown_file() {
    # use a temporary subshell to switch to the temp user
    (
        # Overwrite IONOS_USERNAME and IONOS_PASSWORD with values from temporary files
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        # Execute commands using the temporary user
        ionosctl ipblock delete -af
        ionosctl datacenter delete -af
    )

    # original IONOS_USERNAME IONOS_PASSWORD are restored
    ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
