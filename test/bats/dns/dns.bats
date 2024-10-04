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
    run ionosctl dns zone create --name "cli-test-$(randStr 6).com" --enabled false -o json 2> /dev/null
    assert_success

    zone_id=$(echo "$output" | jq -r '.id')
    assert_regex "$zone_id" "$uuid_v4_regex"

    echo "created dns zone $zone_id"
    echo "$zone_id" > /tmp/bats_test/zone_id
}

@test "List and retrieve DNS Zone by ID" {
    run ionosctl dns zone list -o json 2> /dev/null
    assert_success

    zone_id=$(cat /tmp/bats_test/zone_id)
    run ionosctl dns zone get --zone "${zone_id}" -o json 2> /dev/null
    assert_success
}

@test "Update DNS Zone" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    run ionosctl dns zone update --zone "${zone_id}" --description "updated description" --enabled false -o json 2> /dev/null
    assert_success
}

@test "Create DNS Record in Zone" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    run ionosctl dns record create --zone "${zone_id}" --name "test$(randStr 6)" --type A --content 1.2.3.4 -o json 2> /dev/null
    assert_success

    record_id=$(echo "$output" | jq -r '.id')
    assert_regex "$record_id" "$uuid_v4_regex"

    echo "created dns record $record_id"
    echo "$record_id" > /tmp/bats_test/record_id
}

@test "List and retrieve DNS Record by ID" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    run ionosctl dns record list --zone "${zone_id}" -o json 2> /dev/null
    assert_success

    record_count=$(echo "$output" | jq '.items' | jq length)
    assert [ "$record_count" -gt 0 ]
    echo "found $record_count records"
    echo "$record_count" > /tmp/bats_test/record_count

    record_id=$(cat /tmp/bats_test/record_id)
    run ionosctl dns record get --zone "${zone_id}" --record "${record_id}" -o json 2> /dev/null
    assert_success
}

@test "Update DNS Record" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    record_id=$(cat /tmp/bats_test/record_id)

    run ionosctl dns record update --zone "${zone_id}" --record "${record_id}" --ttl 120 -o json 2> /dev/null
    assert_success
}

@test "Get and update DNS Zone File" {
    zone_id=$(cat /tmp/bats_test/zone_id)
    run ionosctl dns zone file get --zone "${zone_id}" -o text 2> /dev/null
    assert_success

    echo "$output" > /tmp/bats_test/zone_file
    echo "test$(randStr 6) 60 IN A 1.2.3.4" >> /tmp/bats_test/zone_file

    run ionosctl dns zone file update --zone "${zone_id}" --zone-file /tmp/bats_test/zone_file -o json 2> /dev/null
    assert_success

    run ionosctl dns record list --zone "${zone_id}" -o json 2> /dev/null
    assert_success

    record_count=$(echo "$output" | jq '.items' | jq length)
    assert [ "$record_count" -gt "$(cat /tmp/bats_test/record_count)" ]
}

@test "Create DNS Secondary Zone" {
    run ionosctl dns secondary-zone create --name "cli-test-$(randStr 6).com" --primary-ips 1.2.3.4,5.6.7.8 -o json 2> /dev/null
    assert_success

    zone_id=$(echo "$output" | jq -r '.id')
    assert_regex "$zone_id" "$uuid_v4_regex"

    echo "created dns secondary zone $zone_id"
    echo "$zone_id" > /tmp/bats_test/secondary_zone_id
}

@test "List and retrieve DNS Secondary Zone by ID" {
    run ionosctl dns secondary-zone list -o json 2> /dev/null
    assert_success

    zone_id=$(cat /tmp/bats_test/secondary_zone_id)
    run ionosctl dns secondary-zone get --zone "${zone_id}" -o json 2> /dev/null
    assert_success
}

@test "Update DNS Secondary Zone" {
    zone_id=$(cat /tmp/bats_test/secondary_zone_id)
    run ionosctl dns secondary-zone update --zone "${zone_id}" --description "updated description" -o json 2> /dev/null
    assert_success
}

@test "Start and check Transfer for Secondary Zone" {
    zone_id=$(cat /tmp/bats_test/secondary_zone_id)
    run ionosctl dns secondary-zone transfer start --zone "${zone_id}" -o json 2> /dev/null
    assert_success

    run ionosctl dns secondary-zone transfer get --zone "${zone_id}" -o json 2> /dev/null
    assert_success
}

teardown_file() {
    zone_id=$(cat /tmp/bats_test/zone_id)
    secondary_zone_id=$(cat /tmp/bats_test/secondary_zone_id)

    echo "cleaning up dns zone $zone_id and secondary zone $secondary_zone_id with all related records"
    run ionosctl dns secondary-zone delete --zone "$secondary_zone_id" -f
    run ionosctl dns zone delete --zone "$zone_id" -f

    run ionosctl token delete --token "$IONOS_TOKEN" -f
    unset IONOS_TOKEN

    rm -rf /tmp/bats_test
}
