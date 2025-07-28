#!/usr/bin/env bats

# tags: label

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}


@test "Create temporary sub-user" {
    echo "$(randStr 16)@$(randStr 8).ionosctl.test" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/email
    echo "$(randStr 12)" > /tmp/bats_test/password

    run ionosctl user create --first-name "test-user-$(randStr 4)" --last-name "test-last-$(randStr 4)" \
        --email "$(cat /tmp/bats_test/email)" --password "$(cat /tmp/bats_test/password)" -o json 2> /dev/null
    assert_success
    echo "$output" | jq -r '.id' > /tmp/bats_test/user_id

    run ionosctl group create --name "test-group-$(randStr 4)" \
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

@test "Create Datacenter Label" {
    datacenter_name="cli-test-datacenter-$(randStr 8)"
    run ionosctl datacenter create --name "$datacenter_name" -o json 2> /dev/null
    assert_success

    datacenter_id=$(echo "$output" | jq -r '.id')
    echo "$datacenter_id"

    assert_output -p "\"name\": \"$datacenter_name\""

    echo "$datacenter_id" > /tmp/bats_test/datacenter_id
    echo "$datacenter_name" > /tmp/bats_test/datacenter_name

    label_key=$(echo "testlabelkey$(randStr 8)" | tr '[:upper:]' '[:lower:]')
    label_value=$(echo "testlabelvalue$(randStr 8)" | tr '[:upper:]' '[:lower:]')

    sleep 60

    run ionosctl label add --resource-type datacenter --datacenter-id "$datacenter_id" --label-key "$label_key" --label-value "$label_value" -o json 2> /dev/null
    assert_success
    assert_output -p "\"type\": \"label\""
    assert_output -p "\"key\": \"$label_key\""
    assert_output -p "\"value\": \"$label_value\""

}



@test "Create Ipblock Label" {
    ipblock_name="cli-test-ipblock-$(randStr 8)"

    run ionosctl ipblock create --name "$ipblock_name" -o json 2> /dev/null
    assert_success

    ipblock_id=$(echo "$output" | jq -r '.id')
    echo "$ipblock_id" > /tmp/bats_test/ipblock_id

    assert_output -p "\"type\": \"ipblock\""

    sleep 60

    gateway_id=$(echo "$output" | jq -r '.id')
    label_key=$(echo "testlabelkey$(randStr 8)" | tr '[:upper:]' '[:lower:]')
    label_value=$(echo "testlabelvalue$(randStr 8)" | tr '[:upper:]' '[:lower:]')

    run ionosctl label add --resource-type ipblock --ipblock-id "$ipblock_id" --label-key "$label_key" --label-value "$label_value" -o json 2> /dev/null
    assert_success
    assert_output -p "\"type\": \"label\""
    assert_output -p "\"key\": \"$label_key\""
    assert_output -p "\"value\": \"$label_value\""
}

@test "Delete  at least the two previously created labels" {

      initial_labels=$(run ionosctl label list --no-headers | wc -l)
      run ionosctl label remove --all --force
      remaining_labels=$(run ionosctl label list --no-headers | wc -l)

      diff=$(( initial_labels - remaining_labels ))
      assert_success
      run [ "$diff" -ge 2 ]

      run ionosctl ipblock delete --ipblock-id "$(cat /tmp/bats_test/ipblock_id)" -f
      assert_success

      run ionosctl datacenter delete --datacenter-id "$(cat /tmp/bats_test/datacenter_id)" -f
      assert_success
}

teardown_file() {
    (
        export IONOS_USERNAME="$(cat /tmp/bats_test/email)"
        export IONOS_PASSWORD="$(cat /tmp/bats_test/password)"

        ionosctl apigateway gateway delete -af
    )

    ionosctl user delete --user-id "$(cat /tmp/bats_test/user_id)" -f
    ionosctl group delete --group-id "$(cat /tmp/bats_test/group_id)" -f

    rm -rf /tmp/bats_test
}