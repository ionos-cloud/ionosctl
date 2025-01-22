#!/usr/bin/env bats

# tags: dns

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

@test "Create temporary sub-user with DNS permissions" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "test-user-$(randStr 4)" --last-name "test-last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl group create --name "test-group-$(randStr 4)" --access-dns \
        -w -t 300 -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/group_id

    run ionosctl group user add --user-id "$(cat /tmp/bats_test/user_id)" \
        --group-id "$(cat /tmp/bats_test/group_id)" -o json 2> /dev/null
    assert_success

    run ionosctl token generate --ttl 1h
    assert_success
    echo "$output" > /tmp/bats_test/token
}

@test "Create DNS Zone" {
    zone_name="cli-test-$(randStr 6).space"
    zone_name=$(echo "$zone_name" | tr '[:upper:]' '[:lower:]')
    run ionosctl dns zone create --name "$zone_name" --enabled false -o json 2> /dev/null
    assert_success

    zone_id=$(echo "$output" | jq -r '.id')
    assert_regex "$zone_id" "$uuid_v4_regex"

    assert_output -p "\"zoneName\": \"$zone_name\""
    assert_output -p "\"state\": \"AVAILABLE\""

    echo "created dns zone $zone_id ($zone_name)"
    echo "$zone_id" > /tmp/bats_test/zone_id
    echo "$zone_name" > /tmp/bats_test/zone_name
}

@test "List DNS Zones" {
    zone_name=$(cat /tmp/bats_test/zone_name)

    # List Zones (JSON output)
    run ionosctl dns zone list -o json 2> /dev/null
    assert_success
    assert_output -p "\"zoneName\": \"$zone_name\""

    # List Zones (Column output)
    run ionosctl dns zone list --cols name --no-headers
    assert_success
    assert_output -p "$zone_name"
}

@test "Create DNS Record" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    record_name="record-$(randStr 6)"
    record_name=$(echo "$record_name" | tr '[:upper:]' '[:lower:]')
    run ionosctl dns record create --zone "$zone_id" --name "$record_name" --type A --content 192.168.0.1 -o json 2> /dev/null
    assert_success

    record_id=$(echo "$output" | jq -r '.id')
    assert_regex "$record_id" "$uuid_v4_regex"

    # Verify specific fields
    assert_output -p "\"name\": \"$record_name\""
    assert_output -p "\"type\": \"A\""
    assert_output -p "\"content\": \"192.168.0.1\""
    assert_output -p "\"state\": \"AVAILABLE\""

    echo "created dns record $record_id ($record_name)"
    echo "$record_id" > /tmp/bats_test/record_id
    echo "$record_name" > /tmp/bats_test/record_name
}

@test "List DNS Records" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    record_name=$(cat /tmp/bats_test/record_name)

    # List Records (JSON output)
    run ionosctl dns record list --zone "$zone_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$record_name\""

    # List Records (Column output)
    run ionosctl dns record list --zone "$zone_id" --cols name --no-headers
    assert_success
    assert_output "$record_name"
}

@test "Get DNS Record" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    record_id=$(cat /tmp/bats_test/record_id)
    record_name=$(cat /tmp/bats_test/record_name)

    # Get Record by ID
    run ionosctl dns record get --zone "$zone_id" --record "$record_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$record_name\""
    assert_output -p "\"content\": \"192.168.0.1\""

    # Get Record by Name
    run ionosctl dns record get --zone "$zone_id" --record "$record_name" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$record_name\""
    assert_output -p "\"content\": \"192.168.0.1\""
}

@test "Update DNS Record" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    record_name=$(cat /tmp/bats_test/record_name)

    run ionosctl dns record update --zone "$zone_id" --record "$record_name" --ttl 120 -o json 2> /dev/null
    assert_success

    # Verify updated field
    assert_output -p "\"ttl\": 120"
}

@test "Zone File Operations" {
    zone_id=$(cat /tmp/bats_test/zone_id)

    # Get Zone File
    run ionosctl dns zone file get --zone "$zone_id" -o text 2> /dev/null
    assert_success

    echo "$output" > /tmp/bats_test/zone_file
    echo "test$(randStr 6) 60 IN A 1.2.3.4" >> /tmp/bats_test/zone_file

    # Update Zone File
    run ionosctl dns zone file update --zone "$zone_id" --zone-file /tmp/bats_test/zone_file -o json 2> /dev/null
    assert_success

    # Verify Zone File Update
    run ionosctl dns record list --zone "$zone_id" -o json 2> /dev/null
    assert_success

    record_count=$(echo "$output" | jq '.items' | jq length)
    assert [ "$record_count" -gt 0 ]
}

@test "Delete DNS Record" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    record_name=$(cat /tmp/bats_test/record_name)

    run ionosctl dns record delete --zone "$zone_id" --record "$record_name" -f
    assert_success
}

@test "Get and update DNS Zone File" {
    zone_id=$(cat /tmp/bats_test/zone_id)

    # Get the zone file
    run ionosctl dns zone file get --zone "$zone_id" -o text 2> /dev/null
    assert_success
    echo "$output" > /tmp/bats_test/zone_file
    echo "retrieved DNS zone file for zone $zone_id"

    # Update the zone file
    echo "test$(randStr 6) 60 IN A 192.168.0.2" >> /tmp/bats_test/zone_file
    run ionosctl dns zone file update --zone "$zone_id" --zone-file /tmp/bats_test/zone_file -o json 2> /dev/null
    assert_success
    echo "updated DNS zone file for zone $zone_id"

    # Verify new record added via zone file update
    run ionosctl dns record list --zone "$zone_id" -o json 2> /dev/null
    assert_success
    record_count=$(echo "$output" | jq '.items' | jq length)
    echo "$record_count" > /tmp/bats_test/record_count
    assert [ "$record_count" -gt 0 ]
}

@test "Create DNS Secondary Zone" {
    # Create a secondary DNS zone
    run ionosctl dns secondary-zone create --name "cli-test-$(randStr 6).space" --primary-ips 1.2.3.4,5.6.7.8 -o json 2> /dev/null
    assert_success
    zone_id=$(echo "$output" | jq -r '.id')
    assert_regex "$zone_id" "$uuid_v4_regex"

    echo "created secondary DNS zone $zone_id"
    echo "$zone_id" > /tmp/bats_test/secondary_zone_id
}

@test "List and retrieve DNS Secondary Zone by ID" {
    # List all secondary zones
    run ionosctl dns secondary-zone list -o json 2> /dev/null
    assert_success

    # Retrieve specific secondary zone by ID
    zone_id=$(cat /tmp/bats_test/secondary_zone_id)
    run ionosctl dns secondary-zone get --zone "$zone_id" -o json 2> /dev/null
    assert_success
}

@test "Update DNS Secondary Zone" {
    zone_id=$(cat /tmp/bats_test/secondary_zone_id)

    # Update secondary zone description
    run ionosctl dns secondary-zone update --zone "$zone_id" --description "updated secondary zone description" -o json 2> /dev/null
    assert_success
}

@test "Start and check Transfer for Secondary Zone" {
    zone_id=$(cat /tmp/bats_test/secondary_zone_id)

    # Start transfer
    run ionosctl dns secondary-zone transfer start --zone "$zone_id" -o json 2> /dev/null
    assert_success

    # Check transfer status
    run ionosctl dns secondary-zone transfer get --zone "$zone_id" -o json 2> /dev/null
    assert_success
}

@test "Delete DNS Zone" {
    zone_name=$(cat /tmp/bats_test/zone_name)

    run ionosctl dns zone delete --zone "$zone_name" -f
    assert_success
}

@test "Attempt to delete non-existent DNS Zone" {
    run ionosctl dns zone delete --zone "nonexistent-zone.com" -f
    assert_failure
    assert_output -p "could not find zone by name"
}

teardown_file() {
    (
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        ionosctl dns zone delete -af
    )

    ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}
