#!/usr/bin/env bats

load '../libs/bats-assert/load'
load '../libs/bats-support/load'

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

retry_command() {
    local n=1
    local max=3
    local delay=10
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
        assert_regex "$resource_id" "$uuid_v4_regex"
        echo "$resource_id"
        return 0
    fi

    run bash -c "$create_command"
    assert_success
    local new_resource_id=$(echo "$output")
    assert_regex "$new_resource_id" "$uuid_v4_regex"
    echo "$new_resource_id"
}
