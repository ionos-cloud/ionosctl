#!/usr/bin/env bats

# Auto-detect LIBS_PATH if not set (e.g. running bats directly instead of via test/run.sh)
export LIBS_PATH="${LIBS_PATH:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../libs" && pwd)}"

load "${LIBS_PATH}/bats-assert/load"
load "${LIBS_PATH}/bats-support/load"

bats_require_minimum_version 1.5.0

# Override bats' run to separate stderr from stdout (requires bats >= 1.5.0).
# - $output contains only stdout, so jq parsing works without 2>/dev/null
# - $stderr is captured separately and dumped (along with stdout) on failure
eval "$(declare -f run | sed '1s/run/__bats_original_run/')"

# Redact sensitive data from a string (IPs, JWTs, UUIDs).
# For CLI flag values, use redact_args instead.
redact() {
    sed -E \
        -e 's/eyJ[A-Za-z0-9_-]+\.eyJ[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+/***JWT***/g' \
        -e 's/[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}(\/[0-9]{1,2})?/***IP***/g' \
        -e 's/([0-9a-f]{4})[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{8}([0-9a-f]{4})/\1...\2/g'
}

# Redact CLI args: walks $@ and replaces entire values after sensitive flags.
# Handles multi-word values (e.g. PEM contents passed via $(cat file.pem)).
redact_args() {
    local redact_next=false
    local sensitive='--db-password|--password|--private-key|--psk-key|--key-secret|--certificate-chain|--certificate|--token|--secret|--email|--user'
    for arg in "$@"; do
        if $redact_next; then
            printf '***REDACTED*** '
            redact_next=false
        elif [[ "$arg" =~ ^($sensitive)$ ]]; then
            printf '%s ' "$arg"
            redact_next=true
        elif [[ "$arg" =~ ^($sensitive)= ]]; then
            printf '%s=***REDACTED*** ' "${arg%%=*}"
        else
            printf '%s ' "$arg"
        fi
    done
}

skip_if_suite_failed() {
    if [[ -f "$BATS_FILE_TMPDIR/suite_failed" ]]; then
        skip "skipped due to prior test failure ($(cat "$BATS_FILE_TMPDIR/suite_failed"))"
    fi
}

# === SKIP-AFTER-FAILURE ===
# When a test fails, all subsequent tests in the same suite are skipped.
# This relies on setup() and teardown() defined here.
#
# IMPORTANT: If you add setup() or teardown() to your test, 
# you MUST call "skip_if_suite_failed" (if setup()) 
# or mark_suite_failed_on_test_failure (if teardown()) as first line
# Without these calls, skip-after-failure silently stops working for that suite.
setup() {
    skip_if_suite_failed
    if [[ -f /tmp/bats_test/token ]]; then
        export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
    fi
}

run() {
    # In setup_file/teardown_file, BATS_TEST_TMPDIR is unset and --separate-stderr fails.
    # Fall back to original run behavior outside test functions.
    if [[ -z "${BATS_TEST_TMPDIR:-}" ]]; then
        __bats_original_run "$@"
        return
    fi
    skip_if_suite_failed
    __bats_original_run --separate-stderr "$@"
    if [[ "$status" -ne 0 ]]; then
        __failed_count=$(( ${__failed_count:-0} + 1 ))
        local n=$__failed_count
        local cmd
        cmd=$(redact_args "$@" | redact)
        local out
        out=$(echo "$output" | redact)
        local err
        err=$(echo "$stderr" | redact)

        eval "__failed_cmd_$n=\$cmd"
        eval "__failed_out_$n=\$out"
        eval "__failed_err_$n=\$err"
    elif [[ -n "${stderr:-}" && "${stderr}" == *"Warning:"* ]]; then
        # Capture warnings from successful commands for diagnostics.
        # Printed in teardown only if the test ultimately fails.
        __warn_count=$(( ${__warn_count:-0} + 1 ))
        local n=$__warn_count
        local cmd
        cmd=$(redact_args "$@" | redact)
        local err
        err=$(echo "$stderr" | redact)
        eval "__warn_cmd_$n=\$cmd"
        eval "__warn_err_$n=\$err"
    fi
}

# Print deferred diagnostics only when the test failed, then mark suite as failed.
teardown() {
    if [[ -z "${BATS_TEST_COMPLETED:-}" ]]; then
        local n=${__failed_count:-0}
        if [[ $n -gt 0 ]]; then
            if [[ $n -eq 1 ]]; then
                echo "=== Failed command ==="
            else
                echo "=== Failed commands ($n) ==="
            fi
            for i in $(seq 1 "$n"); do
                local cmd err out
                eval "cmd=\$__failed_cmd_$i"
                eval "out=\$__failed_out_$i"
                eval "err=\$__failed_err_$i"

                if [[ $n -gt 1 ]]; then
                    echo "[$i] $cmd"
                else
                    echo "\$ $cmd"
                fi
                if [[ -n "$out" ]]; then
                    echo "stdout: $out"
                fi
                if [[ -n "$err" ]]; then
                    echo "stderr: $err"
                fi
                if [[ $i -lt $n ]]; then
                    echo ""
                fi
            done
        fi
        local wn=${__warn_count:-0}
        if [[ $wn -gt 0 ]]; then
            echo "=== Warnings from successful commands ==="
            for i in $(seq 1 "$wn"); do
                local wcmd werr
                eval "wcmd=\$__warn_cmd_$i"
                eval "werr=\$__warn_err_$i"
                echo "\$ $wcmd"
                echo "stderr: $werr"
            done
        fi
        echo "$BATS_TEST_NAME" > "$BATS_FILE_TMPDIR/suite_failed"
    fi
}

mark_suite_failed_on_test_failure() {
    if [[ -z "${BATS_TEST_COMPLETED:-}" ]]; then
        echo "$BATS_TEST_NAME" > "$BATS_FILE_TMPDIR/suite_failed"
    fi
}

# assert_stderr: like assert_output but checks $stderr (captured via --separate-stderr).
# Supports: assert_stderr -p "substring"  |  assert_stderr "exact match"
assert_stderr() {
    local mode="equal"
    local expected

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -p|--partial) mode="partial"; shift ;;
            *) expected="$1"; shift ;;
        esac
    done

    if [[ "$mode" == "partial" ]]; then
        if [[ "$stderr" != *"$expected"* ]]; then
            echo "-- stderr does not contain substring --"
            echo "substring : $expected"
            echo "stderr    : $(echo "$stderr" | redact)"
            echo "--"
            return 1
        fi
    else
        if [[ "$stderr" != "$expected" ]]; then
            echo "-- stderr is not equal to expected --"
            echo "expected : $expected"
            echo "actual   : $(echo "$stderr" | redact)"
            echo "--"
            return 1
        fi
    fi
}

# refute_stderr: like refute_output but checks $stderr.
# Supports: refute_stderr -p "substring"  |  refute_stderr --partial "substring"
refute_stderr() {
    local mode="equal"
    local expected

    while [[ $# -gt 0 ]]; do
        case "$1" in
            -p|--partial) mode="partial"; shift ;;
            *) expected="$1"; shift ;;
        esac
    done

    if [[ "$mode" == "partial" ]]; then
        if [[ "$stderr" == *"$expected"* ]]; then
            echo "-- stderr should not contain substring --"
            echo "substring : $expected"
            echo "stderr    : $(echo "$stderr" | redact)"
            echo "--"
            return 1
        fi
    else
        if [[ "$stderr" == "$expected" ]]; then
            echo "-- stderr should not be equal to --"
            echo "value : $(echo "$stderr" | redact)"
            echo "--"
            return 1
        fi
    fi
}

# assert_output_not_empty: fail if $output is empty. Use after assert_success
# before parsing output with jq, to catch --wait silent output loss.
assert_output_not_empty() {
    if [[ -z "$output" ]]; then
        echo "-- output is empty (expected non-empty) --"
        if [[ -n "${stderr:-}" ]]; then
            echo "stderr: $(echo "$stderr" | redact)"
        fi
        echo "--"
        return 1
    fi
}

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

randStr() {
    size=$1
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c "$size"
}

# Function to add padding to base64url encoded strings
add_padding() {
    local input=$1
    local padded=$input
    while [ $(( ${#padded} % 4 )) -ne 0 ]; do
        padded="${padded}="
    done
    echo "$padded"
}

generate_ssh_key() {
    ssh_key_path="id_rsa_$(randStr 8)"
    ssh-keygen -t rsa -b 4096 -f "$ssh_key_path" -N "" >/dev/null 2>&1
    echo $ssh_key_path
}

# Bash function to check if an IP is in a CIDR
# First argument: IP
# Second argument: CIDR
# Returns: 0 if IP is in CIDR, 1 otherwise
ip_in_cidr() {
    local ip=$1
    local cidr=$2

    local ip_dec=$(ip_to_dec "$ip")
    local cidr_ip=$(echo $cidr | cut -d '/' -f 1)
    local cidr_ip_dec=$(ip_to_dec "$cidr_ip")

    local prefix=$(echo $cidr | cut -d '/' -f 2)
    local mask=$((0xffffffff << (32 - prefix) & 0xffffffff))

    # Check if IP is within the CIDR block
    if (( (ip_dec & mask) == (cidr_ip_dec & mask) )); then
        return 0
    else
        return 1
    fi
}

# Helper function to convert IP address to decimal
ip_to_dec() {
    local ip=$1
    local a b c d

    IFS='.' read -r a b c d <<< "$ip"
    echo $(( (a << 24) + (b << 16) + (c << 8) + d ))
}

# retry_until retries a command until a condition is met
# First argument: command to run
# Second argument: condition to check
# Third argument: max retries (default: 6)
# Fourth argument: delay between retries (default: 30)
# Returns: 0 if condition is met, 1 otherwise
retry_until() {
    local command="$1"
    local condition="$2"
    local max_retries=${3:-6}
    local delay=${4:-20}
    local n=1

    while true; do
        output=$(eval "$command")
        if eval "$condition"; then
            return 0
        else
            echo "$output"
            if [[ $n -lt $max_retries ]]; then
                ((n++))
                sleep $delay
            else
                fail "$n attempts exhausted for command: $command, condition: $condition"
                return 1
            fi
        fi
    done
}

# Old func which simply calls retry_until expecting status code 0
retry_command() {
    local cmd_to_run="$*"
    local command_condition='[ $? -eq 0 ]'

    retry_until "\"$cmd_to_run\"" "$command_condition"
}

find_or_create_resource() {
    local find_command="$1"
    local create_command="$2"

    bash -c "$find_command"
    if [ "$status" -eq 0 ] && [ -n "$output" ]; then
        echo "$output"
        return 0
    fi

    bash -c "$create_command"
    echo "$output"
}
