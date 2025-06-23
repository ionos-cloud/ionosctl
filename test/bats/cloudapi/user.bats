#!/usr/bin/env bats

# tags: user, group, token, whoami, login, logout, cfg, root

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

@test "60 second token will expire" {
    run ionosctl token generate --ttl 60s
    assert_success
    echo "$output" > /tmp/bats_test/token_60s

    # Wait for token to expire
    sleep 61
    export IONOS_TOKEN=$(cat /tmp/bats_test/token_60s)
    run ionosctl whoami
    assert_output -p "Authentication failed!"
    assert_output -p "Using token for authentication"
}

@test "Create User" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "first-$(randStr 4)" --last-name "last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success

    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id
}

@test "Verify User via API" {
    user_id=$(cat /tmp/bats_test/user_id)
    email=$(cat /tmp/bats_test/email)

    run ionosctl user get --user-id "$user_id" -o json 2> /dev/null
    assert_success
    assert_equal "$(echo "$output" | jq -r '.id')" "$user_id"

    run ionosctl user get --user-id "$user_id" --cols email --no-headers
    assert_success
    assert_output "$email"

    run ionosctl user list -F email="$(cat /tmp/bats_test/email)" --cols UserId --no-headers
    assert_success
    assert_output "$user_id"
}

@test "Add new User to a new Group" {
    user_id=$(cat /tmp/bats_test/user_id)

    group_name="group-$(randStr 8)"
    run ionosctl group create --s3privilege=true --name "$group_name" --cols GroupId --no-headers
    assert_success
    group_id=$output
    echo "$group_id" > /tmp/bats_test/group_id

    sleep 5

    run ionosctl group user add --group-id "$group_id" \
        --user-id "$user_id" --cols UserId --no-headers 2> /dev/null
    assert_success
    assert_output "$user_id"

    sleep 5

    run ionosctl group user list --group-id "$group_id" --cols UserId --no-headers
    assert_success
    assert_output "$user_id"
}

@test "Create and verify S3Key" {
    skip "Test disabled as S3Key creation is flaky with error: \"The user needs to be part of a group that has ACCESS_S3_OBJECT_STORAGE privilege\""

    user_id=$(cat /tmp/bats_test/user_id)
    run ionosctl user s3key create --user-id "$user_id" -o json 2> /dev/null
    assert_success
    access_key=$(echo "$output" | jq -r '.id')
    secret_key=$(echo "$output" | jq -r '.properties.secretKey')

    # TODO: Make a request to the S3 server to test the credentials

    run ionosctl user s3key list --user-id "$user_id" --cols S3KeyId --no-headers
    assert_output -p "$access_key"
    assert_success

    run ionosctl user s3key get --user-id "$user_id" --s3key-id "$access_key" -o json 2> /dev/null
    assert_success
    assert_equal "$access_key" "$(echo "$output" | jq -r '.id')"
    assert_equal "$secret_key" "$(echo "$output" | jq -r '.properties.secretKey')"

    run ionosctl user s3key delete --user-id "$user_id" --s3key-id "$access_key" -f
    assert_success
}

@test "Test 'ionosctl token' commands" {
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN

    email="$(cat /tmp/bats_test/email)"
    password="$(cat /tmp/bats_test/password)"
    user_id=$(cat /tmp/bats_test/user_id)

    run ionosctl login --user "$email" --password "$password" --force
    assert_success

    # Generate a token and ensure it belongs to this user
    run ionosctl token generate --ttl 1h
    assert_success
    jwt="$output"

    # Parse JWT to get the UserId
    run ionosctl token parse --token "$jwt" --cols UserId --no-headers
    assert_output "$user_id"

    run ionosctl token list --cols TokenId --no-headers
    assert_success
    assert_output -p "$uuid_v4_regex"

    # delete JWT
    run ionosctl token delete --token "$jwt" -f
    assert_success

    # ensure JWT no longer works
    run ionosctl login --token "$jwt" -f
    assert_failure
    assert_output -p "401 Unauthorized"

}

@test "Test 'ionosctl cfg' commands" {
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN

    email="$(cat /tmp/bats_test/email)"
    password="$(cat /tmp/bats_test/password)"
    user_id=$(cat /tmp/bats_test/user_id)

    function check_user_token() {
        local email=$1
        local user_id=$2

        # Fetch JWT from config file location and parse it
        run ionosctl config location
        assert_success
        jwt=$(jq -r '.["userdata.token"]' < "$output")

        # Parse JWT to get the UserId
        run ionosctl token parse --token "$jwt" --cols UserId --no-headers
        assert_output "$user_id"

        # Verify the user identity
        run ionosctl whoami
        assert_success
        assert_output "$email"

        # Verify config file is used in absence of environment variables
        run ionosctl whoami --provenance
        assert_success
        assert_output -p "* [3] Config file settings (userdata.token, userdata.name, userdata.password) (USED)"
        assert_output -p "- Using token for authentication."

        # Verify environment variables are used when present
        export IONOS_TOKEN="$jwt"
        run ionosctl whoami --provenance
        assert_success
        assert_output -p "* [2] Environment Variables (IONOS_TOKEN, IONOS_USERNAME, IONOS_PASSWORD) (USED)"
        assert_output -p "- Using token for authentication."
        unset IONOS_TOKEN

        run ionosctl logout
        assert_success
        assert_output -p "De-authentication successful"
    }

    # login using force
    run ionosctl login --user "$email" --password "$password" --force
    assert_success
    assert_output -p "Config file updated successfully"
    check_user_token "$email" "$user_id"

    # Simulated enter 'y' for the interactive prompt
    run bash -c "echo y | ionosctl login --user '$email' --password '$password'"
    assert_success
    assert_output -p "Do you want to replace it? [y/n]:"
    # Simulated enter username for the interactive prompt
    rm "$(ionosctl cfg location)"
    run bash -c "echo $email | ionosctl login --password '$password'"
    assert_success
    assert_output -p "Config file updated successfully."
    check_user_token "$email" "$user_id"
}

@test "Config file should only work for permissions 600" {
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN

    email="$(cat /tmp/bats_test/email)"
    password="$(cat /tmp/bats_test/password)"

    run ionosctl login --user "$email" --password "$password" --force
    assert_success

    run ionosctl config location
    assert_success

    # Check if the file has the correct permissions
    run stat -c "%a" "$(ionosctl cfg location)"
    assert_success
    assert_output "600"
    rm -f "$config_file"

    # Create a new file with incorrect permissions
    loc="$(ionosctl cfg location)"
    touch "$loc"
    chmod 644 "$loc"

    run ionosctl logout
    assert_failure
    assert_output -p "expected 600, got 644"

    # Should use the config (with incorrect permissions) in this case
    unset IONOS_USERNAME IONOS_PASSWORD
    run ionosctl datacenter list
    assert_failure
    assert_output -p "failed reading auth config file"

    rm "$(ionosctl cfg location)"
}

teardown_file() {
    user_id=$(cat /tmp/bats_test/user_id)
    group_id=$(cat /tmp/bats_test/group_id)

    echo "cleaning up user $user_id and group $group_id"
    run ionosctl user delete --user-id "$user_id" -f
    run ionosctl group delete --group-id "$group_id" -f
    run ionosctl token delete -af

    rm -rf /tmp/bats_test

    # Rollback config
    echo "rolling back config file"
    rm -f "$(ionosctl config location)"
    mv "$(ionosctl config location).bak" "$(ionosctl config location)" || echo "No config file found."
}
