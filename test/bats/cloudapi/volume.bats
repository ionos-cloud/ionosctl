#!/usr/bin/env bats

# tags: volume, cdrom, image

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

    run ionosctl compute group create --name "test-vol-$(randStr 4)" \
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
    run ionosctl compute datacenter create --name "vol-test-$(randStr 8)" --location "es/vit" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/datacenter_id
    sleep 5
}

@test "Create Server" {
    run ionosctl compute server create --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" --name "bats-test-$(randStr 8)" \
     --cores 1 --ram 1GB -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/server_id
}

@test "Image list with filters" {
    run ionosctl compute image list -F imageType=hdd -F location=es/vit --no-headers --cols ImageId
    assert_success
    # Should return at least one image
    [ -n "$output" ]
}

@test "Create volume with HDD image and SSH key" {
    run ionosctl compute image list -F imageAliases=ubuntu:latest -F location="es/vit" -F imageType=hdd --cols ImageId --no-headers
    assert_success
    echo "$output" | head -n 1 > /tmp/bats_test/hdd_image_id

    run ionosctl compute volume create --type "SSD Premium" --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --name "bats-test-$(randStr 8)" --size 50 --image-id "$(cat /tmp/bats_test/hdd_image_id)" \
     --ssh-key-paths /tmp/bats_test/id_rsa.pub -t 600 -w -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/volume_id
}

@test "Attach volume to server" {
    run ionosctl compute server volume attach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/volume_id)" -t 600 -w
    assert_success
}

@test "Attach a CD-ROM with an ISO image" {
    run ionosctl compute image list -F imageAliases=ubuntu:latest -F location="es/vit" -F imageType=CDROM --cols ImageId --no-headers
    assert_success
    echo "$output" | head -n 1 > /tmp/bats_test/iso_image_id

    run ionosctl compute server cdrom attach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --cdrom-id "$(cat /tmp/bats_test/iso_image_id)" --server-id "$(cat /tmp/bats_test/server_id)" -w -t 600 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/cdrom_id
}

@test "Detach volume" {
    run ionosctl compute server volume detach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --volume-id "$(cat /tmp/bats_test/volume_id)" -w -t 600 -f
    assert_success
}

@test "Detach CD-ROM" {
    run ionosctl compute server cdrom detach --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --server-id "$(cat /tmp/bats_test/server_id)" --cdrom-id "$(cat /tmp/bats_test/cdrom_id)" -w -t 600 -f
    assert_success
}

@test "Delete volume" {
    run ionosctl compute volume delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" \
     --volume-id "$(cat /tmp/bats_test/volume_id)" -f -w -t 600
    assert_success
}

@test "Delete Server" {
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
