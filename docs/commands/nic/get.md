---
description: Get a NIC
---

# Get

## Usage

```text
ionosctl nic get [flags]
```

## Description

Use this command to get information about a specified NIC from specified Data Center and Server.

Required values to run command:
- Data Center Id
- Server Id
- NIC Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for get
      --ignore-stdin           Force command to execute without user input
      --nic-id string          The unique NIC Id [Required flag]
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
  -v, --verbose                Enable verbose output
```

## Examples

```text
ionosctl nic get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-id 2978400e-da90-405f-905e-8200d4f48158 
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   2       []
```

