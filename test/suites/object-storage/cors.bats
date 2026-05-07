#!/usr/bin/env bats

# paths: commands/object-storage/*

load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-cors-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --location "$TEST_REGION"
    assert_success

    echo "created test bucket for CORS tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        # Clean up any leftover CORS config before deleting the bucket
        run ionosctl object-storage bucket cors delete --name "$TEST_BUCKET_NAME" -f
        run ionosctl object-storage object delete --name "$TEST_BUCKET_NAME" --all -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" -f
    fi
}

# --- validation ---

@test "object-storage bucket cors get: missing --name flag returns error" {
    run ionosctl object-storage bucket cors get
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket cors put: missing --name flag returns error" {
    run ionosctl object-storage bucket cors put
    assert_failure
}

@test "object-storage bucket cors put: missing --json-properties returns error" {
    run ionosctl object-storage bucket cors put --name some-bucket
    assert_failure
    assert_stderr -p "requires at least 2 options"
}

@test "object-storage bucket cors delete: missing --name flag returns error" {
    run ionosctl object-storage bucket cors delete
    assert_failure
    assert_stderr -p "requires at least 1 option"
}

@test "object-storage bucket cors get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage bucket cors get --name some-bucket 
    assert_failure
    assert_stderr -p "object storage credentials not found"
}

@test "object-storage bucket cors put: --json-properties-example prints example JSON" {
    run ionosctl object-storage bucket cors put --json-properties-example
    assert_success
    assert_output -p "CORSRules"
    assert_output -p "AllowedOrigins"
    assert_output -p "AllowedMethods"
}

@test "object-storage bucket cors put: nonexistent file returns error" {
    run ionosctl object-storage bucket cors put --name "$TEST_BUCKET_NAME" --json-properties "/tmp/nonexistent-cors-file.json"
    assert_failure
}

# --- CORS: put, get, delete lifecycle ---

@test "object-storage bucket cors put: apply CORS configuration" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"CORSRules":[{"AllowedOrigins":["http://www.example.com"],"AllowedMethods":["GET","PUT","POST"],"AllowedHeaders":["*"],"ExposeHeaders":["x-amz-request-id"],"MaxAgeSeconds":3600}]}
EOF
    run ionosctl object-storage bucket cors put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile"
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage bucket cors get: retrieve CORS configuration" {
    run ionosctl object-storage bucket cors get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "http://www.example.com"
    assert_output -p "GET"
}

@test "object-storage bucket cors delete: remove CORS configuration" {
    run ionosctl object-storage bucket cors delete --name "$TEST_BUCKET_NAME" -f
    assert_success
    assert_output -p "deleted successfully"
}

@test "object-storage bucket cors get: after delete returns error" {
    run ionosctl object-storage bucket cors get --name "$TEST_BUCKET_NAME"
    assert_failure
}

# --- CORS: multiple rules ---

@test "object-storage bucket cors put: multiple CORS rules" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"CORSRules":[{"AllowedOrigins":["http://example.com"],"AllowedMethods":["GET"],"AllowedHeaders":["Authorization"]},{"AllowedOrigins":["http://other.com"],"AllowedMethods":["POST","DELETE"],"AllowedHeaders":["*"],"MaxAgeSeconds":600}]}
EOF
    run ionosctl object-storage bucket cors put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile"
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage bucket cors get: multiple rules display" {
    run ionosctl object-storage bucket cors get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "http://example.com"
    assert_output -p "http://other.com"
}

@test "object-storage bucket cors delete: cleanup multiple rules" {
    run ionosctl object-storage bucket cors delete --name "$TEST_BUCKET_NAME" -f
    assert_success
}

# --- CORS: from pretty-printed JSON file ---

@test "object-storage bucket cors put: apply from pretty-printed JSON file" {
    local tmpfile
    tmpfile="$(mktemp)"
    cat > "$tmpfile" <<CORS
{
  "CORSRules": [
    {
      "AllowedOrigins": ["https://app.example.com"],
      "AllowedMethods": ["GET", "HEAD"],
      "AllowedHeaders": ["Content-Type"],
      "ExposeHeaders": ["ETag"],
      "MaxAgeSeconds": 1800
    }
  ]
}
CORS
    run ionosctl object-storage bucket cors put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile"
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage bucket cors get: file-based config applied correctly" {
    run ionosctl object-storage bucket cors get --name "$TEST_BUCKET_NAME"
    assert_success
    assert_output -p "https://app.example.com"
    assert_output -p "Content-Type"
}

@test "object-storage bucket cors delete: cleanup file-based config" {
    run ionosctl object-storage bucket cors delete --name "$TEST_BUCKET_NAME" -f
    assert_success
}
