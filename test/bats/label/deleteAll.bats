#!/usr/bin/env bats

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


@test "Create temporary sub-user with ApiGateway permissions" {
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

    echo "testlabelkey$(randStr 8)" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/label_key
    echo "testlabelvalue$(randStr 8)" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/label_value

    sleep 30

    run ionosctl label add --resource-type datacenter --datacenter-id "$datacenter_id" --label-key "$(cat /tmp/bats_test/label_key)" --label-value "$(cat /tmp/bats_test/label_value)" -o json 2> /dev/null
    assert_success
    assert_output -p "\"type\": \"label\""

}



@test "Create Ipblock Label" {
    ipblock_key=$(cat /tmp/bats_test/label_key)
    ipblock_value=$(cat /tmp/bats_test/label_value)
    ipblock_name="cli-test-ipblock-$(randStr 8)"

    run ionosctl ipblock create --name "$ipblock_name" -o json 2> /dev/null
    assert_success

    ipblock_id=$(echo "$output" | jq -r '.id')
    echo "$ipblock_id" > /tmp/bats_test/ipblock_id

    assert_output -p "\"type\": \"ipblock\""

    sleep 30

    echo "testlabelkey$(randStr 8)" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/label_key
    echo "testlabelvalue$(randStr 8)" | tr '[:upper:]' '[:lower:]' > /tmp/bats_test/label_value


    run ionosctl label add --resource-type ipblock --ipblock-id "$ipblock_id" --label-key "$(cat /tmp/bats_test/label_key)" --label-value "$(cat /tmp/bats_test/label_value)" -o json 2> /dev/null
    assert_success
    assert_output -p "\"type\": \"label\""
}

@test "Delete All Labels" {
      run ionosctl label list --no-headers | wc -l
      run ionosctl label remove --all --force
      # s-au sters
      assert_success
}