#!/usr/bin/env bats

# tags: server, template, volume, cdrom, image, console, nic, lan, ipblock, backupunit, snapshot, gpu

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    rm -rf /tmp/bats_test # Setup could hang if already exists
    mkdir -p /tmp/bats_test

    ssh-keygen -t rsa -b 4096 -N "" -f /tmp/bats_test/id_rsa

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
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
     --create-dc --create-nic --create-backup --reserve-ip \
     -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    sleep 10

    run ionosctl group user add --user-id "$(cat /tmp/bats_test/user_id)" \
     --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success

    (
        # Create temporary token using the temporary user's credentials
        unset IONOS_TOKEN
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        run ionosctl token generate --ttl 1h
        assert_success
        echo "$output" > /tmp/bats_test/token
    )
}

@test "Is temp user" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl whoami
    assert_success
    assert_output "$(cat /tmp/bats_test/email)"
}

@test "Create Datacenter" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    # NOTE: In this test suite we also create a CUBE Server. Cubes can only work with INTEL_SKYLAKE family
    # If you want to change the location, make sure it supports INTEL_SKYLAKE!
    run ionosctl datacenter create --name "volumes-test-$(randStr 8)" --location "es/vit" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5
}

@test "Create Server" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    # CPU-Family should be selected correctly by default
    run ionosctl server create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --cores 1 --ram 4GB -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/server_id
}

@test "Reserve IP. Create NIC" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl lan create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --public -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/lan_id

    run ionosctl ipblock create --location "es/vit" --size 1 --name "bats-test-$(randStr 8)" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.properties.ips[0]' > /tmp/bats_test/ip
    echo "$output" | jq -r '.id' > /tmp/bats_test/ipblock_id
    run ionosctl nic create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" --name "bats-test-$(randStr 8)" --ips "$(cat /tmp/bats_test/ip)" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/nic_id
    sleep 5
}

@test "Creating a nic with a non-existent LAN ID will create a LAN" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl nic create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" \
     --lan-id 123 -w -t 600 -o json 2> /dev/null
    assert_success
    sleep 5

    # A LAN is created by default
    run ionosctl lan get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --lan-id 123 --no-headers --cols Public
    assert_success
    assert_output "false"
}

@test "Attach a volume with an HDD image" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    # Find a suitable image
    run ionosctl image list -F imageAliases=ubuntu:latest -F location="es/vit" -F imageType=hdd --cols ImageId --no-headers
    assert_success
    echo "$output" | head -n 1 > /tmp/bats_test/hdd_image_id

    # Create a volume with a custom b64-encoded userdata cloud config script
    echo -e "#cloud-config\nruncmd:\n - [ mkdir, -p, \"/root/test\" ]\n" | base64 -w 0 > /tmp/bats_test/userdata
    run ionosctl volume create --type "SSD Premium" --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --name "bats-test-$(randStr 8)" --size 50 --image-id "$(cat /tmp/bats_test/hdd_image_id)" \
     --ssh-key-paths /tmp/bats_test/id_rsa.pub --user-data "$(cat /tmp/bats_test/userdata)" -t 600 -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/volume_id

    run ionosctl server volume attach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/volume_id)" -t 600 -w
    assert_success
}

@test "Attach a CD-ROM with an ISO image" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    # Find a suitable image
    run ionosctl image list -F imageAliases=ubuntu:latest -F location="es/vit" -F imageType=CDROM --cols ImageId --no-headers
    assert_success
    echo "$output" | head -n 1 > /tmp/bats_test/iso_image_id

    run ionosctl server cdrom attach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --cdrom-id "$(cat /tmp/bats_test/iso_image_id)" --server-id "$(cat /tmp/bats_test/server_id)" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/cdrom_id
}

@test "Attach a volume with a backupunit public image" {
    skip "Disabled because flaky - backupunit creation is not stable"

    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl backupunit create --name "bats$(randStr 6)" --email "$(cat /tmp/bats_test/email)" \
     --password "$(cat /tmp/bats_test/password)" -w -t 600 -o json 2> /dev/null
    assert_success
    assert_regex "$output" "$uuid_v4_regex"
    echo "$output" | jq -r '.id' > /tmp/bats_test/backupunit_id

    # get-sso-url
    run ionosctl backupunit get-sso-url --backupunit-id "$(cat /tmp/bats_test/backupunit_id)" -o json 2> /dev/null
    assert_success

    run ionosctl image list -F location="es/vit" -F cloudInit=V1 -F imageType=hdd -F imageAliases=ubuntu:20 --cols ImageId --no-headers
    assert_success
    image_id="$output"
    assert_regex "$image_id" "$uuid_v4_regex"
    echo "$image_id" | head -n 1 > /tmp/bats_test/ubuntu_image_id

    run ionosctl volume create --type "HDD" --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --name "bats-test-$(randStr 8)" --size 50 --image-id "$(cat /tmp/bats_test/ubuntu_image_id)" \
     --backupunit-id "$(cat /tmp/bats_test/backupunit_id)" --ssh-key-paths /tmp/bats_test/id_rsa.pub \
     -t 600 -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/backup_volume_id

    # attach
    run ionosctl server volume attach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/backup_volume_id)" \
     -t 450 -w -o json 2> /dev/null
    assert_success
}

@test "Server Console is accessible. Token is valid." {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    # Get the token from ionosctl server token get command
    run ionosctl server token get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --no-headers
    assert_success

    run ionosctl server console get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --no-headers
    assert_success

    # expect the remote console page to be accessible
    run curl "$(echo "$output" | grep -o 'https://[^ ]*')"
    assert_success
    assert_output --partial "<title>Remote Console</title>"
}

@test "SSH into the server. Userdata created a directory" {
    skip "Skipping temporarily, as now, sometimes, SSH is not immediately available after server creation"

    local retries=5
    local delay=10
    while ((retries-- > 0)); do
        run ssh -o StrictHostKeyChecking=no -i /tmp/bats_test/id_rsa root@"$(cat /tmp/bats_test/ip)" 'ls /root/test'
        assert_success
        ((retries > 0)) && sleep $delay
    done
    assert_success  # fail if all retries exhausted
}
@test "Detach Volume, CD-ROM" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

#    run ionosctl server volume detach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
#     --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/backup_volume_id)" -w -t 600 -f
#    assert_success

    run ionosctl server volume detach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/volume_id)" -w -t 600 -f
    assert_success

    run ionosctl server cdrom detach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --cdrom-id "$(cat /tmp/bats_test/cdrom_id)" -w -t 600 -f
    assert_success
}

@test "Delete NIC, LAN" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl nic delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --nic-id "$(cat /tmp/bats_test/nic_id)" -w -f -t 600
    assert_success

    run ionosctl lan delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --lan-id "$(cat /tmp/bats_test/lan_id)" -w -f -t 600
    assert_success
}

@test "Delete Server" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl server delete \
     --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/server_id)" -w -t 600 -f
    assert_success
}

@test "Get and verify XS template" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

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

@test "Create Cube Server with Direct Attached Storage" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl server create --name "bats-test-$(randStr 8)" --type "CUBE" \
     -k /tmp/bats_test/id_rsa.pub --template-id "$(cat /tmp/bats_test/template_id)" \
     --image-id "$(cat /tmp/bats_test/hdd_image_id)" --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     -w -t 400 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/cube_server_id
    assert_equal "$(echo "$output" | jq -r '.properties.type')" "CUBE"

    run ionosctl server get --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/cube_server_id)" --no-headers --cols Type
    assert_success
    assert_output -p "CUBE"
}

@test "Create de/fra/2 Datacenter for GPU Server" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl datacenter create --name "gpu-test-$(randStr 8)" --location "de/fra/2" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id_gpu
    sleep 5
}

@test "Create GPU Server" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    # Create a GPU server using the test-only GPU server custom template
    run ionosctl server create --name "bats-gpu-test-$(randStr 8)" --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" \
     --type "GPU" --template-id "6913ed82-a143-4c15-89ac-08fb375a97c5" -w -t 600 -o json 2> /dev/null

    assert_success
    assert_output -p "GPU"
    echo "$output" | jq -r '.id' > /tmp/bats_test/gpu_server_id
}

@test "List GPUs for Server" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl server gpu list --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" \
     --server-id "$(cat /tmp/bats_test/gpu_server_id)" -o json 2> /dev/null
    assert_success

    if [ "$(echo "$output" | jq -r '.items | length')" -gt 0 ]; then
        echo "$output" | jq -r '.items[0].id' > /tmp/bats_test/gpu_id
    else
        fail "No GPUs found for the server"
    fi
}

@test "Get GPU by ID" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl server gpu get --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" \
     --server-id "$(cat /tmp/bats_test/gpu_server_id)" --gpu-id "$(cat /tmp/bats_test/gpu_id)" -o json 2> /dev/null
    assert_success
}

@test "Delete GPU Server" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" --server-id "$(cat /tmp/bats_test/gpu_server_id)" -f -w -t 600
    assert_success
}

@test "Delete GPU Datacenter" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id_gpu)" -f -w -t 600
    assert_success
}

@test "Delete CUBE" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl server delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --server-id "$(cat /tmp/bats_test/cube_server_id)" -f -w -t 600
    assert_success
}

@test "Delete Volumes" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl volume delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --volume-id "$(cat /tmp/bats_test/volume_id)" -f -w -t 600
    assert_success

}

@test "Delete Datacenter" {
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f -w -t 600
    assert_success
}

@test "Delete IPBlock" {
    skip "Disabled because flaky - teardown should handle deletion fine anyway"
    sleep 60

    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl ipblock delete -i "$(cat /tmp/bats_test/ipblock_id)" -f -w -t 600
    assert_success
}

@test "Delete Backupunit" {
    skip "Disabled because flaky - teardown should handle deletion fine anyway"

    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

    run ionosctl backupunit delete --backupunit-id "$(cat /tmp/bats_test/backupunit_id)" -f
    assert_success
}

teardown_file() {
    # use a temporary subshell to switch to the temp user
    (
        # Overwrite IONOS_TOKEN with the value from the temporary file
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"

        # Execute commands using the temporary user
        ionosctl ipblock delete -af
        ionosctl datacenter delete -af
        ionosctl backupunit delete -af
    )

    # original IONOS_USERNAME IONOS_PASSWORD are restored
    ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
