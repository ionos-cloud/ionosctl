#!/usr/bin/env bats

# tags: cdn

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'


setup_file() {
    export IONOS_TOKEN=$(ionosctl token generate)
    mkdir -p /tmp/bats_test
}

@test "Create cdn distribution" {
      RULES='[
        {
          "prefix": "/api",
          "scheme": "http/https",
          "upstream": {
            "caching": true,
            "geoRestrictions": {
              "allowList": ["CN", "RU"]
            },
            "host": "clitest.example.com",
            "rateLimitClass": "R500",
            "sniMode": "distribution",
            "waf": true
          }
        },
        {
          "prefix": "/api2",
          "scheme": "http/https",
          "upstream": {
            "caching": false,
            "geoRestrictions": {
              "blockList": ["CN", "RU"]
            },
            "host": "server2.example.com",
            "rateLimitClass": "R10",
            "sniMode": "origin",
            "waf": false
          }
        }
      ]'
    run ionosctl cdn ds create --domain site.$(randStr 6).com --certificate-id 24e82787-a60e-4f18-8764-856bafc378b4 --routing-rules "$RULES" -o json 2> /dev/null
    assert_success

    distribution_id=$(echo "$output" | jq -r '.id')
    assert_regex "$distribution_id" "$uuid_v4_regex"

    echo "created cdn distribution $distribution_id"
    echo "$distribution_id" > /tmp/bats_test/distribution_id
}

@test "List and retrieve cdn distribution by ID" {
    run ionosctl cdn ds list -o json 2> /dev/null
    assert_success

    distribution_id=$(cat /tmp/bats_test/distribution_id)
    run ionosctl cdn ds get --distribution-id "${distribution_id}" -o json 2> /dev/null
    assert_success
}

@test "Update CDN distribution" {
    RULES_UPDATED='
    [
        {
            "prefix": "/api",
            "scheme": "http/https",
            "upstream": {
                "caching": true,
                "geoRestrictions": {
                    "allowList": [
                        "CN",
                        "RU"
                    ]
                },
                "host": "clitest.example.com",
                "rateLimitClass": "R500",
                "sniMode": "distribution",
                "waf": true
            }
        },
        {
            "prefix": "/api2",
            "scheme": "http/https",
            "upstream": {
                "caching": false,
                "geoRestrictions": {
                    "blockList": [
                        "CN",
                        "RU"
                    ]
                },
                "host": "server2.example.com",
                "rateLimitClass": "R10",
                "sniMode": "origin",
                "waf": false
            }
        }
    ]'
     distribution_id=$(cat /tmp/bats_test/distribution_id)
     run ionosctl cdn ds update --distribution-id "${distribution_id}" --routing-rules "$RULES_UPDATED" -o json 2> /dev/null
     assert_success
}

@test "Get routing rules for CDN distribution" {
    distribution_id=$(cat /tmp/bats_test/distribution_id)
    run ionosctl cdn ds rr get --distribution-id "${distribution_id}" -o json 2> /dev/null
    assert_success

    record_count=$(echo "$output" | jq '.items' | jq length)
    assert [ "$record_count" -eq 2 ]
}

teardown_file() {
    distribution_id=$(cat /tmp/bats_test/distribution_id)

    echo "cleaning up cdn distribution $distribution_id"
    run ionosctl cdn ds delete --distribution-id "$distribution_id" -f

    run ionosctl token delete --token "$IONOS_TOKEN" -f
    unset IONOS_TOKEN

    rm -rf /tmp/bats_test
}
