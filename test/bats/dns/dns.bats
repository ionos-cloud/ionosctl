#!/usr/bin/env bats

# tags: dns

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    mkdir -p /tmp/bats_test
}

@test "Create DNS Zone" {
    zone_name="cli-test-$(randStr 6).com"
    zone_name=$(echo "$zone_name" | tr '[:upper:]' '[:lower:]')
    run ionosctl dns zone create --name "$zone_name" --enabled false -o json 2> /dev/null
    assert_success

    zone_id=$(echo "$output" | jq -r '.id')
    assert_regex "$zone_id" "$uuid_v4_regex"

    # Verify specific fields
    assert_output -p "\"zoneName\": \"$zone_name\""
    assert_output -p "\"state\": \"AVAILABLE\""

    echo "created dns zone $zone_id ($zone_name)"
    echo "$zone_id" > /tmp/bats_test/zone_id
    echo "$zone_name" > /tmp/bats_test/zone_name
}

@test "List DNS Zones" {
    zone_id=$(cat /tmp/bats_test/zone_id)
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
    record_id=$(cat /tmp/bats_test/record_id)
    record_name=$(cat /tmp/bats_test/record_name)

    # List Records (JSON output)
    run ionosctl dns record list --zone "$zone_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"$record_name\""

    # List Records (Column output)
    run ionosctl dns record list --zone "$zone_id" --cols name --max-results 1 --no-headers
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

@test "Delete DNS Record by ID" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    record_id=$(cat /tmp/bats_test/record_id)

    run ionosctl dns record delete --zone "$zone_id" --record "$record_id" -f
    assert_success
}

@test "Delete DNS Zone by Name" {
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
    if [[ -f /tmp/bats_test/zone_id ]]; then
        zone_id=$(cat /tmp/bats_test/zone_id)
        run ionosctl dns zone delete --zone "$zone_id" -f
    fi

    if [[ -f /tmp/bats_test/record_id ]]; then
        zone_id=$(cat /tmp/bats_test/zone_id)
        record_id=$(cat /tmp/bats_test/record_id)
        run ionosctl dns record delete --zone "$zone_id" --record "$record_id" -f
    fi

    run ionosctl token delete --token "$IONOS_TOKEN" -f
    unset IONOS_TOKEN

    rm -rf /tmp/bats_test
}
