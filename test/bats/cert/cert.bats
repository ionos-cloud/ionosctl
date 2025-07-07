#!/usr/bin/env bats

# tags: cert, certificate, certificate-manager, certmanager

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'

setup_file() {
  mkdir -p /tmp/bats_test
  uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'

  run ionosctl token generate --ttl 4h
  if [[ "$status" -ne 0 ]]; then
    echo "Token generation failed: $output" >&2
    exit 1
  fi

  echo "$output" > /tmp/bats_test/token
}

setup() {
  if [[ -f /tmp/bats_test/token ]]; then
    export IONOS_TOKEN="$(cat /tmp/bats_test/token)"
  fi
}

@test "Alias: certs list works" {
  run ionosctl certs list -o json 2> /dev/null
  assert_success
}

@test "Alias: certificate-manager legacy create (hidden)" {
  cert_name="bats-hidden-alias-$(randStr 6)"
  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -subj "/CN=hidden.example.com" \
    -keyout /tmp/bats_test/key_hidden.pem \
    -out /tmp/bats_test/cert_hidden.pem

  run ionosctl certificate-manager create \
    --certificate-name "$cert_name" \
    --certificate "$(cat /tmp/bats_test/cert_hidden.pem)" \
    --certificate-chain "$(cat /tmp/bats_test/cert_hidden.pem)" \
    --private-key "$(cat /tmp/bats_test/key_hidden.pem)" \
    -o json 2> /dev/null
  assert_success

  cert_id=$(echo "$output" | jq -r '.id')
  echo "$cert_id" > /tmp/bats_test/cert_hidden_id
  assert_regex "$cert_id" "$uuid_v4_regex"
}

@test "Delete hidden alias cert" {
  cert_id=$(cat /tmp/bats_test/cert_hidden_id)
  run ionosctl cert delete --certificate-id "$cert_id" -f 2> /dev/null
  assert_success
}

@test "certmanager api-version fallback outputs v2 and warning" {
  run ionosctl certmanager api-version
  assert_success
  assert_output -p "v2.0"
  assert_output -p "deprecated"
}

@test "List certificates" {
  run ionosctl certmanager cert list -o json 2> /dev/null
  assert_success
}

@test "Create certificate from flags" {
  # generate key + cert
  cert_name="bats-cert-flag-$(randStr 6)"
  echo "$cert_name" > /tmp/bats_test/cert_name

  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -subj "/CN=bats.example.com" \
    -keyout /tmp/bats_test/key.pem \
    -out /tmp/bats_test/cert.pem

  run ionosctl certmanager cert create \
    --certificate-name "$cert_name" \
    --certificate "$(cat /tmp/bats_test/cert.pem)" \
    --certificate-chain "$(cat /tmp/bats_test/cert.pem)" \
    --private-key "$(cat /tmp/bats_test/key.pem)" \
    -o json 2> /dev/null
  assert_success

  cert_id=$(echo "$output" | jq -r '.id')
  echo "$cert_id" > /tmp/bats_test/cert_id
  assert_regex "$cert_id" "$uuid_v4_regex"
}

@test "Get certificate" {
  cert_id=$(cat /tmp/bats_test/cert_id)
  run ionosctl certmanager cert get --certificate-id "$cert_id" -o json 2> /dev/null
  assert_success
  assert_output -p "\"id\": \"$cert_id\""
}

@test "Create certificate from file paths" {
  cert_name="bats-cert-file-$(randStr 6)"
  echo "$cert_name" > /tmp/bats_test/cert_name_file

  openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
    -subj "/CN=batsfile.example.com" \
    -keyout /tmp/bats_test/key_file.pem \
    -out /tmp/bats_test/cert_file.pem

  run ionosctl certmanager cert create \
    --certificate-name "$cert_name" \
    --certificate-path /tmp/bats_test/cert_file.pem \
    --certificate-chain-path /tmp/bats_test/cert_file.pem \
    --private-key-path /tmp/bats_test/key_file.pem \
    -o json 2> /dev/null
  assert_success

  cert_id=$(echo "$output" | jq -r '.id')
  echo "$cert_id" > /tmp/bats_test/cert_id_file
  assert_regex "$cert_id" "$uuid_v4_regex"
}

@test "Update certificate name" {
  cert_id=$(cat /tmp/bats_test/cert_id_file)
  new_name="bats-cert-updated-$(randStr 4)"
  run ionosctl certmanager cert update --certificate-id "$cert_id" --certificate-name "$new_name" -o json 2> /dev/null
  assert_success

  run ionosctl certmanager cert get --certificate-id "$cert_id" -o json 2> /dev/null
  assert_output -p "\"name\": \"$new_name\""
}

@test "Delete certificate (from flags)" {
  cert_id=$(cat /tmp/bats_test/cert_id)
  run ionosctl certmanager cert delete --certificate-id "$cert_id" -f 2> /dev/null
  assert_success
}

@test "Delete certificate (from files)" {
  cert_id=$(cat /tmp/bats_test/cert_id_file)
  run ionosctl certmanager cert delete --certificate-id "$cert_id" -f 2> /dev/null
  assert_success
}

teardown_file() {
  ionosctl certmanager cert delete --all --force
  rm -rf /tmp/bats_test
}
