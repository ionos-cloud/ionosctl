---
description: Create a NIC
---

# NicCreate

## Usage

```text
ionosctl nic create [flags]
```

## Aliases

For `nic` command:
```text
[n]
```

## Description

Use this command to create/add a new NIC to the target Server. You can specify the name, ips, dhcp and Lan Id the NIC will sit on. If the Lan Id does not exist it will be created.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run a command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --dhcp                   Set to false if you wish to disable DHCP on the NIC (default true)
  -f, --force                  Force command to execute without user input
  -F, --format strings         Collection of fields to be printed on output (default [NicId,Name,Dhcp,LanId,Ips,State])
  -h, --help                   help for create
      --ips strings            IPs assigned to the NIC. This can be a collection
      --lan-id int             The LAN ID the NIC will sit on. If the LAN ID does not exist it will be created (default 1)
  -n, --name string            The name of the NIC
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
  -t, --timeout int            Timeout option for Request for NIC creation [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NIC creation to be executed
```

## Examples

```text
ionosctl nic create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --name demoNic
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   1       []
RequestId: 67bdb2fb-b1ee-419a-9bcf-f8ea4b800653
Status: Command nic create has been successfully executed
```

