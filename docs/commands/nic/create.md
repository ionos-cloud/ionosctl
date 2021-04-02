---
description: Create a NIC
---

# Create

## Usage

```text
ionosctl nic create [flags]
```

## Description

Use this command to create a new NIC on your account. You can specify the name, ips, dhcp and Lan Id the NIC will sit on. If the Lan Id does not exist it will be created.

You can wait for the action to be executed using `--wait` option.

Required values to run a command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
  -h, --help                   help for create
      --ignore-stdin           Force command to execute without user input
      --lan-id int             The LAN ID the NIC will sit on. If the LAN ID does not exist it will be created (default 1)
      --nic-dhcp               Set to false if you wish to disable DHCP on the NIC (default true)
      --nic-ips strings        IPs assigned to the NIC. This can be a collection
      --nic-name string        The name of the NIC
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
      --timeout int            Timeout option for NIC to be created [seconds] (default 60)
      --wait                   Wait for NIC to be created
```

## Examples

```text
ionosctl nic create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-name demoNic
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   1       []
RequestId: 67bdb2fb-b1ee-419a-9bcf-f8ea4b800653
Status: Command nic create has been successfully executed
```

