#!/usr/bin/env bats

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"

setup_file() {
    uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    ip_regex='^([0-9]{1,3}\.){3}[0-9]{1,3}(\/[0-9]{1,2})?$'
}

randStr() {
    size=$1
    cat /dev/urandom | tr -dc 'a-zA-Z0-9' | head -c "$size"
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

retry_command() {
    local max=6
    local delay=30

    local n=1
    while true; do
        "$@" && break || {
            if [[ $n -lt $max ]]; then
                ((n++))
                echo "Command failed. Attempt $n/$max:"
                sleep $delay;
            else
                fail "$n attempts exhausted for command: $*"
                break
            fi
        }
    done
}

find_or_create_resource() {
    local find_command="$1"
    local create_command="$2"

    run bash -c "$find_command"
    if [ "$status" -eq 0 ] && [ -n "$output" ]; then
        local resource_id=$(echo "$output")
        echo "$resource_id"
        return 0
    fi

    run bash -c "$create_command"
    assert_success
    local new_resource_id=$(echo "$output")
    echo "$new_resource_id"
}
