#!/usr/bin/env bats

# tags: apigateway, api-gateway, route, gateway, customDomains, upstreams
# nume de fisiere only sau directoare

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    rm -rf /tmp/bats_test
    mkdir -p /tmp/bats_test

    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    # verifica daca este ipv4
}

setup() {
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

#content
@test "Create ApiGateway Gateway" {
    gateway_name="cli-test-gateway"
    run ionosctl apigateway gateway create --name "$gateway_name" -o json 2> /dev/null
    assert_success

    gateway_id=$(echo "$output" | jq -r '.id')

    assert_output -p "\"name\": \"cli-test-gateway\""
    assert_output -p "\"status\": \"PROVISIONING\""

    echo "created apigateway gateway $gateway_id ($gateway_name)"
    echo "$gateway_id" > /tmp/bats_test/gateway_id
    echo "$gateway_name" > /tmp/bats_test/gateway_name
}

@test "List ApiGateway Gateways" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)

    # List Gateway (JSON output)
    run ionosctl apigateway gateway list -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"cli-test-gateway\""

    # List Gateway (Column output)
    run ionosctl apigateway gateway list --cols name --no-headers
    assert_success
    assert_output -p "cli-test-gateway"
}

@test "Get ApiGateway Gateways" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)

    # Get Gateway (JSON output)
    run ionosctl apigateway gateway get --gateway-id "$gateway_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"status\": \"AVAILABLE\""
}

@test "Create ApiGateway Route" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)
    route_name="route"
    run ionosctl apigateway route create --gateway-id "$gateway_id" --name "name_route" --paths "path_route" --methods "GET" --host "host.com" -o json 2> /dev/null
    assert_success

    route_id=$(echo "$output" | jq -r '.id')


    # Verify specific fields
    assert_output -p "\"name\": \"name_route\""
    assert_output -p "\"path_route\""
    assert_output -p "\"GET\""
    assert_output -p "\"host\": \"host.com\""
    assert_output -p "\"status\": \"PROVISIONING\""

    echo "created apigateway route $route_id ($route_name)"
    echo "$route_id" > /tmp/bats_test/route_id
    echo "$route_name" > /tmp/bats_test/route_name
}

@test "List ApiGateway Routes" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)

    # List Records (JSON output)
    run ionosctl apigateway route list --gateway-id "$gateway_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"name_route\""

    # List Records (Column output)
    run ionosctl apigateway route list --gateway-id "$gateway_id" --cols name --no-headers
    assert_success
    assert_output "name_route"
}


@test "Get ApiGateway Route" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)
    route_id=$(cat /tmp/bats_test/route_id)

    # Get Route by ID
    run ionosctl apigateway route get --gateway-id "$gateway_id" --route-id "$route_id" -o json 2> /dev/null
    assert_success
    assert_output -p "\"name\": \"name_route\""
}

@test "Gateway CustomDomains Operations" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)
    certificate_id="bc4738cc-0a51-4f0b-ac20-ae65d930454e"

    # Add CustomDomains
    run ionosctl apigateway gateway customdomains add --gateway-id "$gateway_id" --name name.com --certificate-id "$certificate_id" -o json 2> /dev/null
    assert_success
    assert_output -p "name.com"
    assert_output -p "\"certificateId\": \"$certificate_id\""

    # List CustomDomains
    run ionosctl apigateway gateway customdomains list --gateway-id "$gateway_id" -o json 2> /dev/null
    assert_success
    assert_output -p "name.com"

    #Delete CustomDomains
    ionosctl apigateway gateway customdomains remove --gateway-id "$gateway_id" --custom-domains-id 0
    assert_success
}

@test "Route Upstreams Operations" {
  gateway_id=$(cat /tmp/bats_test/gateway_id)
  route_id=$(cat /tmp/bats_test/route_id)

  # Add Upstream
  run ionosctl apigateway route upstreams add --gateway-id "$gateway_id" --route-id "$route_id" --host name.com -o json 2> /dev/null
  assert_success
  assert_output -p "\"host\": \"name.com\""

  #List Upstream
  run ionosctl apigateway route upstreams list --gateway-id "$gateway_id" --route-id "$route_id" -o json 2> /dev/null
  assert_success
  assert_output -p "\"host\": \"name.com\""

  #Delete Upstream
  run ionosctl apigateway route upstreams remove --gateway-id "$gateway_id" --route-id "$route_id" --upstream-id 1 -o json 2> /dev/null
  assert_success
}

@test "Update ApiGateway Gateway" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)

    # Get Gateway (JSON output)
    run ionosctl apigateway gateway update --gateway-id "$gateway_id" --name cli-test-gateway_update -o json 2> /dev/null
    assert_success
    assert_output -p "\"status\": \"PROVISIONING\""
    assert_output -p "\"name\": \"cli-test-gateway_update\""
}

@test "Update ApiGateway Route" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)
    route_id=$(cat /tmp/bats_test/route_id)

    run ionosctl apigateway route update --gateway-id "$gateway_id" --route-id "$route_id" --name name_route_updated -o json 2> /dev/null
    assert_success

    # Verify updated field
    assert_output -p "\"name\": \"name_route_updated\""
}

@test "Delete ApiGateway Gateway" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)
    run ionosctl apigateway gateway delete --gateway-id "$gateway_id" -f
    assert_success
}

@test "Delete ApiGateway Route" {
    gateway_id=$(cat /tmp/bats_test/gateway_id)
    route_id=$(cat /tmp/bats_test/route_id)
    run ionosctl apigateway route delete --gateway-id "$gateway_id" --route-id "$route_id" -f
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
