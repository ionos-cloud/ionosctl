#!/usr/bin/env bats

# tags: user, group, token, whoami, login, logout, config, root

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    # Backup current config
    (mv "$(ionosctl config location)" "$(ionosctl config location).bak") || echo "No config file found."

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'

    mkdir -p /tmp/bats_test
}

# Create a temporary user. Image upload blocks the specific image name for that account, until
# support is contacted to delete that image from the FTP server. This ensures we block
# the image name for a user which will be deleted anyway, and not the testing account.
@test "Create Temporary User" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "first-$(randStr 4)" --last-name "last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success

    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id
}

@test "Upload image to FTP server" {
    # Export credentials of generated user
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN
    export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
    export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

    dd if=/dev/zero of=/tmp/bats_test/10KB.iso bs=1024 count=10
    dd if=/dev/zero of=/tmp/bats_test/10KB.vhd bs=1024 count=10
    dd if=/dev/zero of=/tmp/bats_test/10MB.iso bs=1M count=10
    dd if=/dev/zero of=/tmp/bats_test/10MB.vhd bs=1M count=10
    dd if=/dev/zero of=/tmp/bats_test/500MB.iso bs=1M count=500
    dd if=/dev/zero of=/tmp/bats_test/500MB.vhd bs=1M count=500

    # Upload all files at once to a single location
    first=$(randStr 8)
    run ionosctl image upload --image /tmp/bats_test/10KB.iso,/tmp/bats_test/10MB.iso,/tmp/bats_test/500MB.iso,/tmp/bats_test/10KB.vhd,/tmp/bats_test/10MB.vhd,/tmp/bats_test/500MB.vhd \
        --rename "$first-10KB,$first-10MB,$first-500MB,$first-10KB,$first-10MB,$first-500MB" --location vit --timeout 800 --verbose
    assert_success

    run ionosctl image list -F public=false --cols name --no-headers
    assert_success
    assert_output -p "$first"
    assert_equal "$(echo "$output" | wc -l)" 6
}
