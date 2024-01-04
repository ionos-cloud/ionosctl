#!/usr/bin/env bats

setup() {
  uuid_v4_regex='^[0-9a-f]{8}-[0-9a-f]{4}-[4][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
}

@test "Create a server under a NAT" {
  run ionosctl ipblock create --location de/txl -o json -w
  echo $output
  [ "$status" -eq 0 ]
  ipblock_id=$(echo "$output" | jq -r '.id')
  [[ $ipblock_id =~ $uuid_v4_regex ]]

  run ionosctl datacenter create --location de/txl -o json -w
  echo $output
  [ "$status" -eq 0 ]
  datacenter_id=$(echo "$output" | jq -r '.id')
  [[ $datacenter_id =~ $uuid_v4_regex ]]

  run ionosctl server create --datacenter-id $datacenter_id --cpu-family INTEL_SKYLAKE --cores 1 --ram 1GB -o json -w
  echo $output
  [ "$status" -eq 0 ]
  server_id=$(echo "$output" | jq -r '.id')
  [[ $server_id =~ $uuid_v4_regex ]]

  run ionosctl ipblock list -F location=de/txl -M 1 -o json
  echo $output
  [ "$status" -eq 0 ]
  found_ip=$(echo "$output" | jq -r '.items[] | .properties.ips[0]')

  run ionosctl natgateway create --datacenter-id $datacenter_id --ips $found_ip -o json
  echo $output
  [ "$status" -eq 0 ]
  nat_id=$(echo "$output" | jq -r '.id')
  [[ $nat_id =~ $uuid_v4_regex ]]
}

teardown() {
  echo "cleaning up..."
  run ionosctl server delete -f --datacenter-id $datacenter_id --server-id $server_id
  [ "$status" -eq 0 ]

  run ionosctl datacenter delete -f --datacenter-id $datacenter_id
  [ "$status" -eq 0 ]

  run ionosctl ipblock delete -f --ipblock-id $ipblock_id
  [ "$status" -eq 0 ]

  run ionosctl natgateway delete -f --natgateway-id $nat_id
  [ "$status" -eq 0 ]

  echo $(ionosctl resource list)
}
