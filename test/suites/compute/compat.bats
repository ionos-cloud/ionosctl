#!/usr/bin/env bats

# paths: commands/compute/location/*, commands/compute/server/*, commands/compute/datacenter/*, commands/compute/k8s/*, commands/compute/image/*, commands/compute/template/*, commands/compute/user/*, commands/compute/helpers/*, commands/compute/waiter/*

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
