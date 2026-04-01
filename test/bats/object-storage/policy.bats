#!/usr/bin/env bats

# tags: object-storage

BATS_LIBS_PATH="${LIBS_PATH:-../libs}" # fallback to relative path if not set
load "${BATS_LIBS_PATH}/bats-assert/load"
load "${BATS_LIBS_PATH}/bats-support/load"
load '../setup.bats'


setup_file() {
    if [[ -z "$IONOS_S3_ACCESS_KEY" || -z "$IONOS_S3_SECRET_KEY" ]]; then
        echo "IONOS_S3_ACCESS_KEY and IONOS_S3_SECRET_KEY must be set for object storage e2e tests" >&2
        return 1
    fi

    export TEST_REGION="${IONOS_S3_TEST_REGION:-eu-central-3}"
    export TEST_BUCKET_NAME="ionosctl-ci-pol-$(randStr 8 | tr '[:upper:]' '[:lower:]')"

    run ionosctl object-storage bucket create --name "$TEST_BUCKET_NAME" --region "$TEST_REGION"
    assert_success

    echo "created test bucket for policy tests: $TEST_BUCKET_NAME"
}

teardown_file() {
    if [[ -n "$TEST_BUCKET_NAME" ]]; then
        # Clean up any leftover policy before deleting the bucket
        run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f
        run ionosctl object-storage bucket delete --name "$TEST_BUCKET_NAME" --recursive -f
    fi
}

# --- validation ---

@test "object-storage policy get: missing --name flag returns error" {
    run ionosctl object-storage policy get 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage policy put: missing --name flag returns error" {
    run ionosctl object-storage policy put 2>&1
    assert_failure
}

@test "object-storage policy put: missing --json-properties returns error" {
    run ionosctl object-storage policy put --name some-bucket 2>&1
    assert_failure
    assert_output -p "requires at least 2 options"
}

@test "object-storage policy delete: missing --name flag returns error" {
    run ionosctl object-storage policy delete 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage policy status: missing --name flag returns error" {
    run ionosctl object-storage policy status 2>&1
    assert_failure
    assert_output -p "requires at least 1 option"
}

@test "object-storage policy get: missing S3 credentials returns error" {
    run env -u IONOS_S3_ACCESS_KEY -u IONOS_S3_SECRET_KEY \
        ionosctl object-storage policy get --name some-bucket \
        --config /dev/null 2>&1
    assert_failure
    assert_output -p "object storage credentials not found"
}

@test "object-storage policy put: --json-properties-example prints example JSON" {
    run ionosctl object-storage policy put --json-properties-example 2>/dev/null
    assert_success
    assert_output -p "Statement"
    assert_output -p "Effect"
    assert_output -p "Principal"
}

@test "object-storage policy put: nonexistent file returns error" {
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "/tmp/nonexistent-policy-file.json" 2>&1
    assert_failure
}

# --- policy format: basic Allow with Principal object ---

@test "object-storage policy put: basic Allow with Principal object" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Version":"2012-10-17","Statement":[{"Sid":"AllowGetObject","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: retrieve basic Allow policy" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "AllowGetObject"
    assert_output -p "s3:GetObject"
    assert_output -p "Allow"
}

@test "object-storage policy status: public policy reports IsPublic true" {
    run ionosctl object-storage policy status --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "true"
}

@test "object-storage policy delete: cleanup basic Allow policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
    assert_output -p "deleted successfully"
}

# --- policy format: wildcard action s3:* ---

@test "object-storage policy put: wildcard action s3:*" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Version":"2012-10-17","Statement":[{"Sid":"AllowAll","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:*"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}","arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: wildcard action displays correctly" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "s3:*"
    assert_output -p "AllowAll"
}

@test "object-storage policy delete: cleanup wildcard action policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy format: multiple actions ---

@test "object-storage policy put: multiple actions" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Version":"2012-10-17","Statement":[{"Sid":"MultiAction","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject","s3:PutObject","s3:DeleteObject"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: multiple actions display" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "s3:GetObject"
    assert_output -p "s3:PutObject"
    assert_output -p "s3:DeleteObject"
}

@test "object-storage policy delete: cleanup multiple actions policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy format: bucket-level resource ---

@test "object-storage policy put: bucket-level resource (ListBucket)" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Version":"2012-10-17","Statement":[{"Sid":"AllowListBucket","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:ListBucket"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: bucket-level resource displays correctly" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "s3:ListBucket"
    assert_output -p "arn:aws:s3:::${TEST_BUCKET_NAME}"
}

@test "object-storage policy delete: cleanup bucket-level policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy format: Principal as wildcard string ---

@test "object-storage policy put: Principal as wildcard string" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Version":"2012-10-17","Statement":[{"Sid":"WildcardPrincipal","Effect":"Allow","Principal":"*","Action":["s3:GetObject"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: wildcard Principal resolves to *" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "WildcardPrincipal"
    assert_output -p "*"
}

@test "object-storage policy delete: cleanup wildcard Principal policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy format: no optional fields (no Id, no Version, no Sid) ---

@test "object-storage policy put: minimal policy without optional fields" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: minimal policy displays correctly" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "Allow"
    assert_output -p "s3:GetObject"
}

@test "object-storage policy delete: cleanup minimal policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy format: policy from file ---

@test "object-storage policy put: apply policy from JSON file" {
    local tmpfile
    tmpfile="$(mktemp)"
    cat > "$tmpfile" <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "FromFile",
      "Effect": "Allow",
      "Principal": {"AWS": ["*"]},
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::${TEST_BUCKET_NAME}/*"]
    }
  ]
}
POLICY
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: file-based policy applied correctly" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "FromFile"
}

@test "object-storage policy delete: cleanup file-based policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy format: combined bucket + object resources ---

@test "object-storage policy put: combined bucket and object resources" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Version":"2012-10-17","Statement":[{"Sid":"AllowListAndGet","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:ListBucket","s3:GetObject"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}","arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: combined resources display" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "AllowListAndGet"
    assert_output -p "arn:aws:s3:::${TEST_BUCKET_NAME}"
}

@test "object-storage policy delete: cleanup combined resources policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy format: policy with Id field ---

@test "object-storage policy put: policy with Id field" {
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Id":"MyPolicyId","Version":"2012-10-17","Statement":[{"Sid":"WithId","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success
    assert_output -p "applied successfully"
}

@test "object-storage policy get: policy with Id retrieves correctly" {
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>/dev/null
    assert_success
    assert_output -p "WithId"
}

@test "object-storage policy delete: cleanup Id policy" {
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
}

# --- policy: delete and confirm absence ---

@test "object-storage policy delete: put then force-delete confirms absence" {
    # Put a policy so we have something to delete
    local tmpfile="$(mktemp)"
    cat > "$tmpfile" <<EOF
{"Version":"2012-10-17","Statement":[{"Sid":"ToDelete","Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::${TEST_BUCKET_NAME}/*"]}]}
EOF
    run ionosctl object-storage policy put --name "$TEST_BUCKET_NAME" --json-properties "$tmpfile" 2>/dev/null
    rm -f "$tmpfile"
    assert_success

    # Delete with --force
    run ionosctl object-storage policy delete --name "$TEST_BUCKET_NAME" -f 2>/dev/null
    assert_success
    assert_output -p "deleted successfully"

    # Confirm policy is gone
    run ionosctl object-storage policy get --name "$TEST_BUCKET_NAME" 2>&1
    assert_failure
}
