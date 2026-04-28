#!/usr/bin/env bats

# tags: location, server, compat, datacenter, k8s, image, template, user

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

@test "location list: reachable via old and new form" {
    run ionosctl compute location list
    assert_success
    old_output="$output"

    run ionosctl location list
    assert_success
    assert_output "$old_output"
}

@test "datacenter list: reachable via old and new form" {
    run ionosctl compute datacenter list
    assert_success
    old_output="$output"

    run ionosctl datacenter list
    assert_success
    assert_output "$old_output"
}

@test "k8s version list: reachable via old and new form" {
    run ionosctl compute k8s version list
    assert_success
    old_output="$output"

    run ionosctl k8s version list
    assert_success
    assert_output "$old_output"
}

@test "image list: reachable via old and new form" {
    run ionosctl compute image list --cols ImageId --no-headers
    assert_success
    old_output="$output"

    run ionosctl image list --cols ImageId --no-headers
    assert_success
    assert_output "$old_output"
}

@test "template list: reachable via old and new form" {
    run ionosctl compute template list
    assert_success
    old_output="$output"

    run ionosctl template list
    assert_success
    assert_output "$old_output"
}

@test "user list: reachable via old and new form" {
    run ionosctl compute user list --cols UserId --no-headers
    assert_success
    old_output="$output"

    run ionosctl user list --cols UserId --no-headers
    assert_success
    assert_output "$old_output"
}
