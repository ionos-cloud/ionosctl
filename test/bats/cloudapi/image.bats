#!/usr/bin/env bats

# tags: user, group, token, whoami, login, logout, config, root

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    mkdir -p /tmp/bats_test

    dd if=/dev/zero of=/tmp/bats_test/10KB.iso bs=1024 count=10
    dd if=/dev/zero of=/tmp/bats_test/10KB.vhd bs=1024 count=10
}

# Create a temporary user. Image upload blocks the specific image name for that account, until
# support is contacted to delete that image from the FTP server. This ensures we block
# the image name for a user which will be deleted anyway, and not the testing account.
@test "Create Temporary User" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "random-$(randStr 4)" --last-name "last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success

    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id
}

@test "Upload multiple images to a single FTP server" {
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    # Upload all files at once to a single location
    random=$(randStr 8)
    run ionosctl image upload --image /tmp/bats_test/10KB.iso,/tmp/bats_test/10KB.vhd \
        --rename "$random-10KB,$random-10KB" --location vit --timeout 800 --verbose
    assert_success
    echo "$output" > /tmp/bats_test/upload_output

    # Check if both images are uploaded
    run ionosctl image list -F public=false --cols name --no-headers
    assert_success
    assert_output -p "$random"
    assert_equal "$(echo "$output" | wc -l)" 2

    # Check if the images are uploaded with the correct names and their properties were set
    run ionosctl image list -F "name=$random-10KB.vhd" --cols cloudInit --no-headers
    assert_success
    assert_output "V1"
    run ionosctl image list -F "name=$random-10KB.iso" -o json
    assert_success
    rhp=$(echo "$output" | jq -r '.items[0].properties.ramHotPlug')
    assert_equal "$rhp" "true"
}

@test "Upload multiple images to a multiple FTP servers" {
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    random=$(randStr 8)
    run ionosctl image upload --image /tmp/bats_test/10KB.iso,/tmp/bats_test/10KB.vhd \
        --rename "$random-10KB,$random-10KB" --location vit,lhr --timeout 800 --cols ImageId --no-headers
    assert_success
    imageIds=$output
    assert_equal "$(echo "$imageIds" | wc -l)" 4

    # Change licence-type to LINUX via update for one of the images
    imageId=$(echo "$imageIds" | head -n 1)
    run ionosctl image update --image-id "$imageId" --licence-type LINUX --cols name --no-headers
    assert_success
    assert_output -p "$random-10KB"

    run ionosctl image list -F public=false --cols location --no-headers
    assert_success
    vit_exists=$(echo "$output" | grep -c "vit")
    lhr_exists=$(echo "$output" | grep -c "lhr")
    assert_equal "$vit_exists" 1
    assert_equal "$lhr_exists" 1
}

@test "Upload image to a mock FTP server" {
    # Make a mock ftp server
    mkdir -p /tmp/bats_test/ftp
    docker run -d -p 21:21 -p 30000-30009:30000-30009 -e "FTP_USER=ftpuser" -e "FTP_PASS=ftppass" -v /tmp/bats_test/ftp:/home/ftpuser/ftp atmoz/sftp:alpine

    # Upload a file
    random=$(randStr 8)
    run ionosctl image upload --image /tmp/bats_test/10KB.iso \
        --ftp-url "ftp://ftpuser:ftppass@localhost:21" --verbose --skip-verify
    assert_success
    echo "$output" > /tmp/bats_test/upload_output

    # Check if file made it to docker ftp server
    run "docker exec -it $(docker ps -q) ls /home/ftpuser/ftp"
    assert_success
    assert_output -p "10KB.iso"
}

@test "Can attach private images as CD-ROM" {
    skip "todo"
}
@test "Can attach private images as HDD" {
    skip "todo"
}
@test "Can attach private images as DAS" {
    skip "todo"
}

@test "Creator of sub-user can delete sub-user private image" {
    run ionosctl image list -F public=false --cols ImageId --no-headers
    assert_success
    num=$(echo "$output" | wc -l)
    image_to_delete=$(echo "$output" | head -n 1)

    run ionosctl image delete --image-id "$image_to_delete" --force --wait-for-request
    assert_success

    # num has decreased
    run ionosctl image list -F public=false --cols ImageId --no-headers
    assert_success
    assert_equal "$(echo "$output" | wc -l)" $((num-1))
}

@test "Can delete all private images" {
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    run ionosctl image delete -af --wait-for-request
    assert_success

#    run ionosctl image list -F public=false --cols ImageId --no-headers
#    assert_success
#    # Disabled check because inconsistent status reports by API resulting in flaky test
#    assert_equal "$(echo "$output" | wc -l)" 0
}

@test "Can delete temp user" {
    run ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" --force
    assert_success
}

teardown_file() {
    rm -rf /tmp/bats_test
}
