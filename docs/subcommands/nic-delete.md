---
description: Delete a NIC
---

# NicDelete

## Usage

```text
ionosctl nic delete [flags]
```

## Description

This command deletes a specified NIC.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* NIC Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
  -t, --timeout int            Timeout option for Request for NIC deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NIC deletion to be executed
```

## Examples

```text
ionosctl nic delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-id 2978400e-da90-405f-905e-8200d4f48158 --force 
RequestId: 14a4bf17-48aa-4f87-b0dc-9c769a4cbdcb
Status: Command nic delete has been successfully executed
```

