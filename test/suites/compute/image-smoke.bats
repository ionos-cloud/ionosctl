#!/usr/bin/env bats

# paths: commands/compute/image/*

# Smoke tests for image upload via a local FTPS mock server.
# These tests verify CLI-side logic only (FTP path routing, flag validation,
# rename, multi-image uploads, TLS certificate handling) without hitting real
# IONOS infrastructure. All tests use --skip-update.
#
# Requirements: python3, pyftpdlib<2 + pyOpenSSL, openssl
# Tests are skipped gracefully if dependencies are missing.
#
# Usage:
#   go build -o ./ionosctl .
#   IONOSCTL_BIN=./ionosctl LIBS_PATH=test/libs bats test/suites/compute/image-smoke.bats

load '../setup.bats'

FTPS_PORT="${FTPS_PORT:-2121}"
FTPS_ROOT="/tmp/bats_smoke_ftp"
FTPS_PID_FILE="/tmp/bats_smoke_ftp/ftps.pid"
FTPS_CERT="/tmp/bats_smoke_ftp/server.crt"
FTPS_KEY="/tmp/bats_smoke_ftp/server.key"
FTPS_USER="testuser"
FTPS_PASS="testpass"
TEST_DIR="/tmp/bats_smoke_test"
HELPERS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/../../helpers"
# Use locally built binary if IONOSCTL_BIN is set, otherwise fall back to PATH
IONOSCTL="${IONOSCTL_BIN:-ionosctl}"

# Common flags for all upload commands targeting the local FTPS server
FTP_FLAGS="--ftp-url localhost --ftp-port $FTPS_PORT --skip-verify --skip-update --timeout 30"
# Same but using --crt-path instead of --skip-verify
FTP_FLAGS_CRT="--ftp-url localhost --ftp-port $FTPS_PORT --crt-path $FTPS_CERT --skip-update --timeout 30"

setup_file() {
    # Check dependencies
    command -v python3 >/dev/null 2>&1 || skip "python3 not found"
    python3 -c "from pyftpdlib.handlers import TLS_FTPHandler" || skip "pyftpdlib with TLS not available (pip install 'pyftpdlib<2' pyOpenSSL)"
    command -v openssl >/dev/null 2>&1 || skip "openssl not found"

    # Kill stale server from a previous interrupted run
    if [ -f "$FTPS_PID_FILE" ]; then
        kill "$(cat "$FTPS_PID_FILE")" || true
        rm -f "$FTPS_PID_FILE"
    fi

    rm -rf "$FTPS_ROOT" "$TEST_DIR"
    mkdir -p "$FTPS_ROOT" "$TEST_DIR"

    # Generate self-signed certificate for FTPS (with SAN for Go TLS verification)
    openssl req -x509 -newkey rsa:2048 -keyout "$FTPS_KEY" -out "$FTPS_CERT" \
        -days 1 -nodes -subj "/CN=localhost" \
        -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

    # Start local FTPS server in background
    python3 "$HELPERS_DIR/ftps_server.py" \
        "$FTPS_PORT" "$FTPS_ROOT" "$FTPS_CERT" "$FTPS_KEY" "$FTPS_PID_FILE" \
        "$FTPS_USER" "$FTPS_PASS" &
    local server_pid=$!
    disown || true

    # Wait for server to be ready
    for i in $(seq 1 20); do
        if [ -f "$FTPS_PID_FILE" ]; then
            break
        fi
        sleep 0.25
    done

    if [ ! -f "$FTPS_PID_FILE" ]; then
        # Clean up the background process since teardown_file won't run
        kill "$server_pid" || true
        rm -rf "$FTPS_ROOT" "$TEST_DIR"
        fail "FTPS server failed to start"
    fi

    # Create test image files of different formats
    for ext in qcow2 vhd iso vmdk img raw vhdx cow qcow vpc vdi; do
        dd if=/dev/zero of="$TEST_DIR/test.$ext" bs=1024 count=10
    done
}

teardown_file() {
    if [ -f "$FTPS_PID_FILE" ]; then
        local pid
        pid="$(cat "$FTPS_PID_FILE")"
        kill "$pid" || true
        # Wait for the process to actually exit before removing its files
        for i in $(seq 1 20); do
            kill -0 "$pid" || break
            sleep 0.1
        done
    fi
    rm -rf "$FTPS_ROOT" "$TEST_DIR"
}

setup() {
    # Clear any real credentials so tests never hit IONOS infrastructure
    unset IONOS_USERNAME IONOS_PASSWORD IONOS_TOKEN
    export IONOS_USERNAME="$FTPS_USER"
    export IONOS_PASSWORD="$FTPS_PASS"

    # Clean upload directories between tests
    rm -f "$FTPS_ROOT/iso-images/"* || true
    rm -f "$FTPS_ROOT/hdd-images/"* || true
}

# =============================================================================
# Normal uploads  - path routing
# =============================================================================

@test "Upload .iso routes to iso-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.iso" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/iso-images/test.iso" ]
    [ ! -f "$FTPS_ROOT/hdd-images/test.iso" ]
}

@test "Upload .img routes to iso-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.img" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/iso-images/test.img" ]
    [ ! -f "$FTPS_ROOT/hdd-images/test.img" ]
}

@test "Upload .vhd routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.vhd" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.vhd" ]
    [ ! -f "$FTPS_ROOT/iso-images/test.vhd" ]
}

@test "Upload .vmdk routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.vmdk" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.vmdk" ]
}

@test "Upload .qcow2 routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.qcow2" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.qcow2" ]
}

@test "Upload .raw routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.raw" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.raw" ]
}

@test "Upload .vhdx routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.vhdx" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.vhdx" ]
}

@test "Upload .cow routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.cow" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.cow" ]
}

@test "Upload .qcow routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.qcow" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.qcow" ]
}

@test "Upload .vpc routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.vpc" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.vpc" ]
}

@test "Upload .vdi routes to hdd-images/" {
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.vdi" $FTP_FLAGS
    assert_success

    [ -f "$FTPS_ROOT/hdd-images/test.vdi" ]
}

# =============================================================================
# Multiple images
# =============================================================================

@test "Upload multiple images routes each to the correct directory" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso,$TEST_DIR/test.vhd" \
        $FTP_FLAGS

    assert_success

    [ -f "$FTPS_ROOT/iso-images/test.iso" ]
    [ -f "$FTPS_ROOT/hdd-images/test.vhd" ]
}

@test "Upload multiple HDD images to hdd-images/" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.vhd,$TEST_DIR/test.vmdk" \
        --rename "$random-a,$random-b" \
        $FTP_FLAGS

    assert_success

    [ -f "$FTPS_ROOT/hdd-images/$random-a.vhd" ]
    [ -f "$FTPS_ROOT/hdd-images/$random-b.vmdk" ]
}

@test "Upload multiple ISO-type images to iso-images/" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso,$TEST_DIR/test.img" \
        --rename "$random-a,$random-b" \
        $FTP_FLAGS

    assert_success

    [ -f "$FTPS_ROOT/iso-images/$random-a.iso" ]
    [ -f "$FTPS_ROOT/iso-images/$random-b.img" ]
}

@test "Upload mix of ISO and HDD images in one command" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso,$TEST_DIR/test.vhd,$TEST_DIR/test.vmdk,$TEST_DIR/test.img" \
        --rename "$random-1,$random-2,$random-3,$random-4" \
        $FTP_FLAGS

    assert_success

    [ -f "$FTPS_ROOT/iso-images/$random-1.iso" ]
    [ -f "$FTPS_ROOT/hdd-images/$random-2.vhd" ]
    [ -f "$FTPS_ROOT/hdd-images/$random-3.vmdk" ]
    [ -f "$FTPS_ROOT/iso-images/$random-4.img" ]
}

@test "Upload all supported HDD formats in one command" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.vhd,$TEST_DIR/test.vmdk,$TEST_DIR/test.qcow2,$TEST_DIR/test.raw,$TEST_DIR/test.vhdx,$TEST_DIR/test.cow,$TEST_DIR/test.qcow,$TEST_DIR/test.vpc,$TEST_DIR/test.vdi" \
        --rename "$random-1,$random-2,$random-3,$random-4,$random-5,$random-6,$random-7,$random-8,$random-9" \
        $FTP_FLAGS

    assert_success

    for i in 1 2 3 4 5 6 7 8 9; do
        ls "$FTPS_ROOT/hdd-images/$random-$i."* >/dev/null 2>&1
    done
}

# =============================================================================
# Rename
# =============================================================================

@test "Upload with --rename changes the filename on server" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --rename "$random-renamed" \
        $FTP_FLAGS

    assert_success

    [ -f "$FTPS_ROOT/iso-images/$random-renamed.iso" ]
    [ ! -f "$FTPS_ROOT/iso-images/test.iso" ]
}

@test "Rename preserves original file extension" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.vmdk" \
        --rename "$random-newname" \
        $FTP_FLAGS

    assert_success

    [ -f "$FTPS_ROOT/hdd-images/$random-newname.vmdk" ]
    # Must not exist without extension or with wrong extension
    [ ! -f "$FTPS_ROOT/hdd-images/$random-newname" ]
}

@test "Rename each image independently in multi-image upload" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso,$TEST_DIR/test.vhd,$TEST_DIR/test.qcow2" \
        --rename "$random-alpha,$random-beta,$random-gamma" \
        $FTP_FLAGS

    assert_success

    [ -f "$FTPS_ROOT/iso-images/$random-alpha.iso" ]
    [ -f "$FTPS_ROOT/hdd-images/$random-beta.vhd" ]
    [ -f "$FTPS_ROOT/hdd-images/$random-gamma.qcow2" ]
    # Originals should not exist
    [ ! -f "$FTPS_ROOT/iso-images/test.iso" ]
    [ ! -f "$FTPS_ROOT/hdd-images/test.vhd" ]
    [ ! -f "$FTPS_ROOT/hdd-images/test.qcow2" ]
}

# =============================================================================
# TLS certificate handling (--crt-path vs --skip-verify)
# =============================================================================

@test "Upload succeeds with --crt-path pointing to server certificate" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --rename "$random-crt" \
        $FTP_FLAGS_CRT

    assert_success

    [ -f "$FTPS_ROOT/iso-images/$random-crt.iso" ]
}

@test "Upload with --crt-path works for HDD images" {
    random=$(randStr 12)
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.vhd" \
        --rename "$random-crt-hdd" \
        $FTP_FLAGS_CRT

    assert_success

    [ -f "$FTPS_ROOT/hdd-images/$random-crt-hdd.vhd" ]
}

@test "Upload fails with --crt-path pointing to non-existent file" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --ftp-url localhost --ftp-port "$FTPS_PORT" \
        --crt-path "/tmp/bats_smoke_test/nonexistent.pem" \
        --skip-update --timeout 30

    assert_failure
}

@test "Upload fails with --crt-path pointing to wrong certificate" {
    # Generate a different self-signed cert that doesn't match the server
    openssl req -x509 -newkey rsa:2048 \
        -keyout "$TEST_DIR/wrong.key" -out "$TEST_DIR/wrong.crt" \
        -days 1 -nodes -subj "/CN=wrong.example.com"

    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --ftp-url localhost --ftp-port "$FTPS_PORT" \
        --crt-path "$TEST_DIR/wrong.crt" \
        --skip-update --timeout 10

    assert_failure
    assert_output -p "dialing FTP server failed"
}

@test "Upload fails without --skip-verify or --crt-path (untrusted cert)" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --ftp-url localhost --ftp-port "$FTPS_PORT" \
        --skip-update --timeout 10

    assert_failure
    assert_output -p "dialing FTP server failed"
}

@test "Multi-image upload with --crt-path" {
    random=$(randStr 12)

    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso,$TEST_DIR/test.vhd,$TEST_DIR/test.qcow2" \
        --rename "$random-a,$random-b,$random-c" \
        $FTP_FLAGS_CRT

    assert_success

    [ -f "$FTPS_ROOT/iso-images/$random-a.iso" ]
    [ -f "$FTPS_ROOT/hdd-images/$random-b.vhd" ]
    [ -f "$FTPS_ROOT/hdd-images/$random-c.qcow2" ]
}

# =============================================================================
# Validation  - invalid extensions
# =============================================================================

@test "Upload rejects .txt extension" {
    touch "$TEST_DIR/test.txt"
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.txt" $FTP_FLAGS

    assert_failure
    assert_output -p "invalid image extension"
}

@test "Upload rejects .zip extension" {
    touch "$TEST_DIR/test.zip"
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.zip" $FTP_FLAGS

    assert_failure
    assert_output -p "invalid image extension"
}

@test "Upload rejects .tar.gz extension" {
    touch "$TEST_DIR/test.tar.gz"
    run $IONOSCTL compute image upload --image "$TEST_DIR/test.tar.gz" $FTP_FLAGS

    assert_failure
    assert_output -p "invalid image extension"
}

@test "Upload rejects file with no extension" {
    touch "$TEST_DIR/noextension"
    run $IONOSCTL compute image upload --image "$TEST_DIR/noextension" $FTP_FLAGS

    assert_failure
    assert_output -p "invalid image extension"
}

@test "Upload rejects mix of valid and invalid extensions" {
    touch "$TEST_DIR/test.pdf"
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso,$TEST_DIR/test.pdf" \
        $FTP_FLAGS

    assert_failure
    assert_output -p "invalid image extension"
    # Nothing should be uploaded when validation fails pre-run
    [ ! -f "$FTPS_ROOT/iso-images/test.iso" ]
}

# =============================================================================
# Validation  - rename count mismatch
# =============================================================================

@test "Upload rejects when rename count < image count" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso,$TEST_DIR/test.vhd" \
        --rename "only-one" \
        $FTP_FLAGS

    assert_failure
    assert_output -p "different lengths"
}

@test "Upload rejects when rename count > image count" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --rename "name-one,name-two" \
        $FTP_FLAGS

    assert_failure
    assert_output -p "different lengths"
}

# =============================================================================
# Validation  - flag dependencies
# =============================================================================

@test "Upload fails without --image flag" {
    run $IONOSCTL compute image upload $FTP_FLAGS

    assert_failure
}

@test "--ftp-port without --ftp-url is rejected" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --ftp-port 2121 --skip-verify \
        --skip-update --location fra

    assert_failure
    assert_output -p "--ftp-port requires --ftp-url"
}

@test "--ftp-port 0 is rejected" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --ftp-url localhost --ftp-port 0 --skip-verify \
        --skip-update

    assert_failure
    assert_output -p "must be between 1 and 65535"
}

@test "--ftp-port 70000 is rejected" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --ftp-url localhost --ftp-port 70000 --skip-verify \
        --skip-update

    assert_failure
    assert_output -p "must be between 1 and 65535"
}

# =============================================================================
# Validation  - non-existent image file
# =============================================================================

@test "Upload fails when image file does not exist" {
    run $IONOSCTL compute image upload \
        --image "/tmp/bats_smoke_test/does-not-exist.qcow2" \
        $FTP_FLAGS

    assert_failure
}

# =============================================================================
# Validation  - bad credentials
# =============================================================================

@test "Upload fails with wrong FTP credentials" {
    export IONOS_USERNAME="wronguser"
    export IONOS_PASSWORD="wrongpass"

    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        $FTP_FLAGS --timeout 10

    assert_failure
    assert_output -p "dialing FTP server failed"
}

# =============================================================================
# Validation  - unreachable server
# =============================================================================

@test "Upload fails when FTP server is unreachable" {
    run $IONOSCTL compute image upload \
        --image "$TEST_DIR/test.iso" \
        --ftp-url localhost --ftp-port 19999 --skip-verify \
        --skip-update --timeout 5

    assert_failure
    assert_output -p "dialing FTP server failed"
}
