#!/usr/bin/env bats

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
