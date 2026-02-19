#!/usr/bin/env bats

# tags: location, server, compat

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

@test "Commands reachable via both old and new form" {
    # Read-only command: location list
    run ionosctl compute location list
    assert_success
    old_output="$output"

    run ionosctl location list
    assert_success
    assert_output "$old_output"
}
