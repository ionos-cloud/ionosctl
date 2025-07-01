#!/usr/bin/env bats

# tags: user, group, token, whoami, login, logout, cfg, root, config

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

# TODO: We need a test that if bad token in config file, an approriate error is shown especailly for whoami

setup_file() {
    # Backup current config
    (mv "$(ionosctl config location)" "$(ionosctl config location).bak") || echo "No config file found."
    # backup legacy config.json at same dir as YAML if exists
    legacy_cfg="$(dirname "$(ionosctl config location)")/config.json"
    if [ -f "$legacy_cfg" ]; then
        mv "$legacy_cfg" "${legacy_cfg}.bak"
    fi

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
    assert_output -p "failed getting username via token"
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
    export IONOS_TOKEN="$jwt"
    run ionosctl whoami
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
        location="$output"
        jwt=$(grep '^[[:space:]]*token:' "$location" \
              | sed -E 's/^[[:space:]]*token:[[:space:]]*//')
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
        assert_output -p "* [3] credentials from config file: token (USED)"

        # Verify environment variables are used when present
        export IONOS_TOKEN="$jwt"
        run ionosctl whoami --provenance
        assert_success
        assert_output -p "* [1] environment variable: IONOS_TOKEN (USED)"
        unset IONOS_TOKEN

        run ionosctl logout
        assert_success
        assert_output -p "Removed credentials from $location but kept URL overrides"
    }

    # login using force
    run ionosctl login --whitelist=dns --user "$email" --password "$password" --force
    assert_success
    assert_success
    assert_output -p "Config file generated at"
    echo "email: $email, user_id: $user_id"
    check_user_token "$email" "$user_id"

    # Simulated enter 'y' for the interactive prompt
    run bash -c "echo y | ionosctl login --whitelist=dns --user '$email' --password '$password'"
    assert_success
    assert_output -p "Do you want to replace it? [y/n]:"

    # Simulated enter username for the interactive prompt
    rm "$(ionosctl cfg location)"
    run bash -c "echo $email | ionosctl login --whitelist=dns --password '$password' --force"
    assert_success
    assert_output -p "Enter your username:"
    assert_output -p "Config file generated at"
    check_user_token "$email" "$user_id"
}

@test "cfg login --example prints sample YAML but does not write file" {
    rm $(ionosctl cfg location) # ensure clean state
    run ionosctl config login --example --whitelist=dns
    echo "done"
    assert_success
    # should include profiles:
    assert_output --partial "profiles:"

    echo "location..."
    run cat $(ionosctl config location)
    assert_failure
    echo "cat location failed!"
    assert_output --partial "No such file or directory"
}

@test "whoami without env uses config token and prints email" {
    if [ -f "$(ionosctl cfg location)" ]; then
        rm "$(ionosctl cfg location)"
    fi

    unset IONOS_TOKEN IONOS_USERNAME IONOS_PASSWORD
    run ionosctl config login --whitelist=dns --user "$(cat /tmp/bats_test/email)" \
      --password "$(cat /tmp/bats_test/password)" --force
    assert_success

    # whoami returns the email
    run ionosctl config whoami
    assert_success
    assert_output --partial "$email"
}

@test "whoami --provenance shows correct auth layer" {
    if [ -f "$(ionosctl cfg location)" ]; then
        rm "$(ionosctl cfg location)"
    fi

    run ionosctl config whoami --provenance
    assert_success
    assert_output --partial "[2] environment variables: IONOS_USERNAME, IONOS_PASSWORD (USED)"

    run ionosctl token generate --ttl 60s
    assert_success
    export IONOS_TOKEN="$output"

    run ionosctl config whoami --provenance
    assert_success
    assert_output --partial "[1] environment variable: IONOS_TOKEN (USED)"

    unset IONOS_TOKEN IONOS_USERNAME IONOS_PASSWORD
    run ionosctl config login --whitelist=dns --user "$(cat /tmp/bats_test/email)" \
      --password "$(cat /tmp/bats_test/password)" --force
    assert_success

    run ionosctl config whoami --provenance
    assert_success
    # config-file token is 3rd in priority
    assert_output --partial "[3] credentials from config file: token (USED)"
}

@test "logout clears credentials but preserves config file" {
  # get YAML path
  run ionosctl config location
  assert_success
  cfg_path="$output"

  # perform logout
  run ionosctl config logout
  assert_success
  assert_output --partial "Removed credentials"

  # YAML file still exists
  run bash -c "[ -f \"$cfg_path\" ]"
  assert_success

  unset IONOS_TOKEN IONOS_USERNAME IONOS_PASSWORD
  run ionosctl config whoami
  assert_output --partial "authentication failed: no credentials found"
}

@test "logout --only-purge-old deletes legacy config.json without touching YAML" {
    # locate YAML and create a fake legacy JSON beside it
    run ionosctl config location
    assert_success
    cfg_path="$output"
    legacy_json="$(dirname "$cfg_path")/config.json"

    cat > "$legacy_json" <<EOF
{
  "userdata.token":"LEGACY",
  "userdata.name":"foo",
  "userdata.password":"bar"
}
EOF

    run bash -c "[ -f \"$legacy_json\" ]"
    assert_success

    # invoke only-purge-old and answer "y"
    run bash -c "echo y | ionosctl config logout --only-purge-old"
    assert_success
    assert_output --partial "Detected legacy config.json"
    assert_output --partial "Delete legacy"

    # legacy JSON should be removed
    run bash -c "[ ! -f \"$legacy_json\" ]"
    assert_success

    # YAML remains
    run bash -c "[ -f \"$cfg_path\" ]"
    assert_success
}

@test "logout skips purge when user answers no" {
    # recreate fake legacy JSON
    run ionosctl config location
    assert_success
    cfg_path="$output"
    legacy_json="$(dirname "$cfg_path")/config.json"

     cat > "$legacy_json" <<EOF
{
  "userdata.token":"LEGACY",
  "userdata.name":"foo",
  "userdata.password":"bar"
}
EOF
    run bash -c "[ -f \"$legacy_json\" ]"
    assert_success

    # invoke only-purge-old and answer "n"
    run bash -c "echo n | ionosctl config logout --only-purge-old"
    assert_success
    assert_output --partial "Detected legacy config.json"
    assert_output --partial "Delete legacy"

    # since answered "n", JSON should still exist
    run bash -c "[ -f \"$legacy_json\" ]"
    assert_success
}

@test "login settings flags are applied" {
    unset IONOS_TOKEN IONOS_USERNAME IONOS_PASSWORD

    email="$(cat /tmp/bats_test/email)"
    password="$(cat /tmp/bats_test/password)"

    run ionosctl config login --custom-names dns=foobar --whitelist=dns --user "$email" --password "$password" \
      --force --version 2.2 --profile-name custom-name --environment made-with-bats
    assert_success
    assert_output -p "Config file generated at"

    # check the config file
    run ionosctl config location
    assert_success
    cfg_path="$output"

    run cat "$cfg_path"
    assert_success
    assert_output --partial "version: 2.2"
    assert_output --partial "name: custom-name"
    assert_output --partial "name: made-with-bats"
    assert_output --partial "environment: made-with-bats"
    assert_output --partial "name: foobar"
}

@test "login --custom-names : --whitelist allows both pre-apply and post-apply name" {
    unset IONOS_TOKEN IONOS_USERNAME IONOS_PASSWORD

    email="$(cat /tmp/bats_test/email)"
    password="$(cat /tmp/bats_test/password)"

    run ionosctl config login --custom-names dns=post-apply-works --whitelist=post-apply-works --user "$email" \
      --password "$password" --force
    assert_success
    assert_output -p "Config file generated at"

    run ionosctl config location
    assert_success
    cfg_path="$output"

    run cat "$cfg_path"
    assert_success
    assert_output --partial "name: post-apply-works"

    run ionosctl config login --custom-names dns=pre-apply-works --whitelist=dns --user "$email" \
      --password "$password" --force
    assert_success
    assert_output -p "Config file generated at"

    run ionosctl config location
    assert_success
    cfg_path="$output"

    run cat "$cfg_path"
    assert_success
    assert_output --partial "name: pre-apply-works"
}

@test "overriding auth (location-less) URL with a bad URL" {
    run ionosctl config location
    assert_success
    cfg_path="$output"

    if [ -f "$cfg_path" ]; then
        rm "$cfg_path"
    fi
    cat > "$cfg_path" <<EOF
version: 1.0
currentProfile: user
profiles:
    - name: user
      environment: prod
      credentials:
        token: <not-important-for-this-test>
environments:
    - name: prod
      products:
        - name: auth
          endpoints:
            - name: https://bad.url-example.com/auth/v1
              skipTlsVerify: false
EOF
    run bash -c "[ -f \"$cfg_path\" ]"
    assert_success

    run ionosctl token list
    assert_failure
    assert_output -p "Error: Get \"https://bad.url-example.com/auth/v1/tokens\": dial tcp: lookup bad.url-example.com"
}

@test "overriding dns (location-based) URL with a new location with bad URL" {
    run ionosctl config location
    assert_success
    cfg_path="$output"

    if [ -f "$cfg_path" ]; then
        rm "$cfg_path"
    fi
    cat > "$cfg_path" <<EOF
version: 1.0
currentProfile: user
profiles:
    - name: user
      environment: prod
      credentials:
        token: <not-important-for-this-test>
environments:
    - name: prod
      products:
        - name: dns
          endpoints:
            - location: de/fra
              name: https://dns.de-fra.ionos.com
              skipTlsVerify: false
            - location: new/loc
              name: https://dns.new-loc.ionos.com
              skipTlsVerify: false
EOF
    run bash -c "[ -f \"$cfg_path\" ]"
    assert_success

    run ionosctl dns zone list --location new/loc
    assert_failure
    assert_output -p "Error: Get \"https://dns.new-loc.ionos.com/zones\": dial tcp: lookup dns.new-loc.ionos.com"
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
    # Restore legacy config if it was backed up
    legacy_cfg="$(dirname "$(ionosctl config location)")/config.json"
    if [ -f "${legacy_cfg}.bak" ]; then
        mv "${legacy_cfg}.bak" "$legacy_cfg"
    fi
}
