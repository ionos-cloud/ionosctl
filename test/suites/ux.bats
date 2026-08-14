#!/usr/bin/env bats

# paths: commands/*

# UX behaviours that are independent of any single product:
#   - trailing "help" is treated as "--help" (#688)
#   - mistyped commands suggest the closest match (did-you-mean, #689)
#   - command errors don't dump the whole flag list (#690)
#   - mistyped flags suggest the closest known flag (flag-suggest)
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

# --- did-you-mean: mistyped commands (#689) ---

@test "did-you-mean: mistyped subcommand suggests the closest match" {
    run ionosctl server craete
    assert_failure
    assert_stderr --partial 'unknown command "craete" for "ionosctl server"'
    assert_stderr --partial "Did you mean this?"
    assert_stderr --partial "create"
    assert_stderr --partial "Run 'ionosctl server --help' for usage."
}

@test "did-you-mean: mistyped top-level command suggests the closest match" {
    run ionosctl dnn
    assert_failure
    assert_stderr --partial 'unknown command "dnn" for "ionosctl"'
    assert_stderr --partial "Did you mean this?"
    assert_stderr --partial "dns"
}

@test "did-you-mean: mistyped command exits non-zero (no silent help)" {
    # Regression: a non-runnable parent used to print help and exit 0 on a typo.
    run ionosctl server craete
    assert_failure
    refute_stderr --partial "USAGE:"
}

# --- no flag dump on errors (#690) ---

@test "no-flag-dump: unknown flag error omits the flag list" {
    run ionosctl dns zone list --nope
    assert_failure
    assert_stderr --partial "unknown flag: --nope"
    refute_stderr --partial "Global Flags:"
    refute_stderr --partial "--verbose"
}

@test "no-flag-dump: missing required flag error omits the flag list" {
    run ionosctl dns zone create
    assert_failure
    refute_stderr --partial "Global Flags:"
    refute_stderr --partial "--verbose"
    # A pointer to --help is still offered.
    assert_stderr --partial "--help"
}

# --- flag-suggest: mistyped flags (unknown-flag did-you-mean) ---

@test "flag-suggest: mistyped local flag suggests the closest match" {
    run ionosctl server list --datacentr
    assert_failure
    assert_stderr --partial "unknown flag: --datacentr"
    assert_stderr --partial "Did you mean this?"
    assert_stderr --partial "--datacenter-id"
    assert_stderr --partial "Run with --help"
}

@test "flag-suggest: omitting the -id suffix still matches (--datacenter -> --datacenter-id)" {
    # Many ID flags end in -id/-ids, which users routinely drop; the suffix is
    # stripped before measuring distance so the bare name still suggests.
    run ionosctl server list --datacenter
    assert_failure
    assert_stderr --partial "unknown flag: --datacenter"
    assert_stderr --partial "--datacenter-id"
}

@test "flag-suggest: mistyped inherited global flag suggests the closest match" {
    run ionosctl dns zone list --verbos
    assert_failure
    assert_stderr --partial "unknown flag: --verbos"
    assert_stderr --partial "Did you mean this?"
    assert_stderr --partial "--verbose"
}

@test "flag-suggest: an unrelated unknown flag offers no suggestion and no flag dump" {
    # Nothing within edit distance 2, so we fall back to the plain unknown-flag
    # error (#690: still no wall of flags).
    run ionosctl dns zone list --zzzzzzzz
    assert_failure
    assert_stderr --partial "unknown flag: --zzzzzzzz"
    refute_stderr --partial "Did you mean this?"
    refute_stderr --partial "Global Flags:"
}
