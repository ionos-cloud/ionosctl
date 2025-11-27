#!/usr/bin/env bats

# tags: cloudapi-v6, client

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
    mkdir -p /tmp/bats_test
}

@test "Query parameters 'limit', 'offset', 'filter', 'order-by' are sent correctly" {
    run ionosctl datacenter list --limit 2 --offset 3 --filters a=b --filters c=d --order-by name --api-url 'test'
    assert_output --partial 'limit=2'
    assert_output --partial 'offset=3'
    assert_output --partial 'filter.a=b'
    assert_output --partial 'filter.c=d'
    assert_output --partial 'order-by=name'
}

@test "no query parameters sent when none provided" {
    run ionosctl datacenter list --api-url 'test'
    refute_output --partial 'filter.'
    refute_output --partial 'order-by='
}

@test "single filter with --filters sent correctly" {
    run ionosctl datacenter list --filters status=active --api-url 'test'
    assert_output --partial 'filter.status=active'
}

@test "single filter with -F sent correctly" {
    run ionosctl datacenter list -F status=active --api-url 'test'
    assert_output --partial 'filter.status=active'
}

@test "Request list default depth is 2" {
    run ionosctl request list --api-url 'test'
    assert_output --partial 'depth=2'
}

@test "Server list default depth is 1" {
    run ionosctl server list --datacenter-id "foo" --api-url 'test'
    assert_output --partial 'depth=1'
}

@test "Using deprecated '--max-results' simply sets limit" {
    run ionosctl datacenter list --max-results 5 --api-url 'test'
    assert_output --partial 'limit=5'
}

@test "Using '--limit' on 'image list' sets maxResults as here 'limit' unsupported" {
    run ionosctl image list --limit 5 --api-url 'test'
    assert_output --partial 'maxResults=5'

    # img is an alias for image
    run ionosctl img list --limit 5 --api-url 'test'
    assert_output --partial 'maxResults=5'
}

@test "For Image API, not setting --limit results in no maxResults query parameter" {
    # This is due a to a bug in the Image API where setting a Filter and MaxResults together causes a weird behavior

    run ionosctl image list --filters name=Ubuntu --api-url 'test'
    refute_output --partial 'maxResults='

    run ionosctl img list --filters name=Ubuntu --api-url 'test'
    refute_output --partial 'maxResults='

    run ionosctl img list -F name=Ubuntu --api-url 'test'
    refute_output --partial 'maxResults='

    # test os.Args check doesnt break (i.e. changing arg doesnt break check)
    run ionosctl -F name=Ubuntu --api-url 'test' img list
    refute_output --partial 'maxResults='
}
