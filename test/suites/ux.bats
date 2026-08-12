#!/usr/bin/env bats

# paths: commands/*

# UX behaviours that are independent of any single product:
#   - trailing "help" is treated as "--help" (#688)
# These invocations short-circuit before any API call, so they need no
# credentials and mutate nothing.

load './setup.bats'

@test "trailing help: 'create help' shows the command help" {
    run ionosctl dns zone create help
    assert_success
    assert_output --partial "USAGE:"
    assert_output --partial "ionosctl dns zone create"
}

@test "trailing help: works with a flag before help (#688 regression)" {
    # The shortcut used to only fire when 'help' was the sole argument, so a
    # preceding flag made it run 'create' and fail.
    run ionosctl dns zone create --name foo help
    assert_success
    assert_output --partial "USAGE:"
    assert_output --partial "ionosctl dns zone create"
}

@test "trailing help: works with a global value flag before help" {
    run ionosctl dns zone create -o json help
    assert_success
    assert_output --partial "USAGE:"
}

@test "trailing help: works with an inline flag value before help" {
    run ionosctl dns zone create --name=foo help
    assert_success
    assert_output --partial "USAGE:"
}

@test "trailing help: a trailing 'help' wins even after a value flag" {
    # A trailing "help" always shows help, even directly after a value flag.
    # To pass the literal "help" as a value, don't put it last (--zone=help).
    run ionosctl dns zone get --zone help
    assert_success
    assert_output --partial "USAGE:"
    assert_output --partial "ionosctl dns zone get"
}
