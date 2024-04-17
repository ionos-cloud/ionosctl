#!/usr/bin/env bats

# tags: user, group, resource, whoami, login

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

location="es/vit"

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'

    mkdir -p /tmp/bats_test
}

@test "Create User" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "first-$(randStr 4)" --last-name "last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success

    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id
}

@test "Verify User via API" {
    skip "disabled, --cols not implemented for user list"

    user_id=$(cat /tmp/bats_test/user_id)
    run ionosctl user get --user-id "$user_id" -o json 2> /dev/null
    assert_success
    assert_equal "$(echo "$output" | jq -r '.id')" "$user_id"

    run ionosctl user get --user-id "$user_id" --cols UserId --no-headers
    assert_success
    assert_output "$user_id"

    run ionosctl user list -F email="$(cat /tmp/bats_test/email)" --cols UserId --no-headers
    assert_success
    assert_output "$user_id"
}

@test "'ionosctl login' as newly created user" {
    email="$(cat /tmp/bats_test/email)"
    password="$(cat /tmp/bats_test/password)"

    # login using force
    run ionosctl login --user "$email" --password "$password" --force
    assert_success
    assert_output -p "Config file updated successfully"

    # Generated token should correspond to the user. Parse it to identify User ID. Delete it
    run ionosctl config location
    assert_success
    jwt=$(jq -r '.["userdata.token"]' < "$output")
    run ionosctl token parse --token "$jwt" --cols UserId --no-headers
    assert_output "$(cat /tmp/bats_test/user_id)"

    # Simulated enter 'y' for the interactive prompt
    run bash -c "echo y | ionosctl login --user '$email' --password '$password'"
    assert_success
    assert_output -p "Do you want to replace it? [y/n]:"
}

teardown_file() {
    echo "cleaning up user $user_id"
    run ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f

    rm -rf /tmp/bats_test
}
