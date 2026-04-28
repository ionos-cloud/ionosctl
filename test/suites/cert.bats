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

#Provider
@test "Create provider" {
    provider_name="provider-name-$(randStr 3)"
    email="sdk-go-v6@cloud.ionos.com"
    server="https://acme-v02.api.letsencrypt.org/directory"
    key_id="provider-key-$(randStr 3)"
    key_secret="provider-secret-$(randStr 3)"

    run ionosctl certmanager provider create \
      --name "$provider_name" \
      --email "$email" \
      --server "$server" \
      --key-id "$key_id" \
      --key-secret "$key_secret" \
      -o json 2> /dev/null
    assert_success

    provider_id=$(echo "$output" | jq -r '.id')
    echo "$provider_id" > /tmp/bats_test/provider_id

    assert_output -p "\"name\": \"$provider_name\""
}

@test "Get Provider" {
  provider_id=$(cat /tmp/bats_test/provider_id)
  run ionosctl certmanager provider get --provider-id "$provider_id" -o json 2> /dev/null
  assert_success
  assert_output -p "\"id\": \"$provider_id\""
}

@test "Update Provider" {
  new_provider_name="new-provider-$(randStr 3)"
  provider_id=$(cat /tmp/bats_test/provider_id)
  run ionosctl certmanager provider update --provider-id "$provider_id" --name "$new_provider_name" -o json 2> /dev/null
  assert_success
  assert_output -p "\"name\": \"$new_provider_name\""
}

#Autocertificate
@test "Create Zone for Autocertificate" {
  #Delete zone if it already exist
  zone_name="devsdkionos.net"
  run ionosctl dns zone delete --zone "$zone_name" -f 2> /dev/null

  run ionosctl dns zone create --name "$zone_name" -o json 2> /dev/null
  assert_success
  assert_output -p "\"zoneName\": \"$zone_name\""

}

@test "Create Autocertificate" {
    provider_id=$(cat /tmp/bats_test/provider_id)

    autocertificate_name="autocertificate-name-$(randStr 3)"
    common_name="devsdkionos.net"
    key_algorithm="rsa4096"
    alternative_names="devsdkionos.net"

    run ionosctl certmanager autocertificate create \
      --name "$autocertificate_name" \
      --provider-id "$provider_id" \
      --common-name "$common_name" \
      --key-algorithm "$key_algorithm" \
      --subject-alternative-names "$alternative_names" \
      -o json 2> /dev/null
    assert_success

    autocertificate_id=$(echo "$output" | jq -r '.id')
    echo "$autocertificate_id" > /tmp/bats_test/autocertificate_id

    assert_output -p "\"name\": \"$autocertificate_name\""
}

@test "Get Autocertificate" {
  autocertificate_id=$(cat /tmp/bats_test/autocertificate_id)
  run ionosctl certmanager autocertificate get --autocertificate-id "$autocertificate_id" -o json 2> /dev/null
  assert_success
  assert_output -p "\"id\": \"$autocertificate_id\""
}

@test "Update Autocertificate" {
  new_autocertificate_name="new-autocertificate-$(randStr 3)"
  autocertificate_id=$(cat /tmp/bats_test/autocertificate_id)
  run ionosctl certmanager autocertificate update --autocertificate-id "$autocertificate_id" --name "$new_autocertificate_name" -o json 2> /dev/null
  assert_success
  assert_output -p "\"name\": \"$new_autocertificate_name\""
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

@test "Delete Autocertificate" {
  autocertificate_id=$(cat /tmp/bats_test/autocertificate_id)
  run ionosctl certmanager autocertificate delete --autocertificate-id "$autocertificate_id" -f 2> /dev/null
  assert_success
}

@test "Delete Provider" {
  provider_id=$(cat /tmp/bats_test/provider_id)
  run ionosctl certmanager provider delete --provider-id "$provider_id" -f 2> /dev/null
  assert_success
}

teardown_file() {
  ionosctl certmanager cert delete --all --force
  rm -rf /tmp/bats_test
}
