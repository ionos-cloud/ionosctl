#!/bin/bash

json_output=$(ionosctl ipblock create -o json -w)
new_json=$(echo "$json_output" | jq -r '{ipblock_id: .id, ips: .properties.ips}')

echo "$new_json"
