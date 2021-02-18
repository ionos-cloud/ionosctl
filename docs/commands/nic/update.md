---
description: Update a NIC
---

# Update

## Usage

```text
ionosctl nic update [flags]
```

## Description

Use this command to update the configuration of a specified NIC. Some restrictions are in place:
The primary address of a NIC connected to a Load Balancer can only be changed by changing the IP of the Load Balancer. 
You can also add additional reserved, public IPs to the NIC.

The user can specify and assign private IPs manually. Valid IP addresses for private networks are 10.0.0.0/8, 172.16.0.0/12 or 192.168.0.0/16.

The value for firewallActive can be toggled between true and false to enable or disable the firewall. When the firewall is enabled, incoming traffic is filtered by all the firewall rules configured on the NIC. 
When the firewall is disabled, then all incoming traffic is routed directly to the NIC and any configured firewall rules are ignored.

You can wait for the action to be executed using `--wait` option.

Required values to run command:
- Data Center Id
- NIC Id

## Options

```text
      --add-ip string          Add IP
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [NicId,Name,Dhcp,LanId,Ips])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
      --datacenter-id string   The unique Data Center Id
  -h, --help                   help for update
      --ignore-stdin           Force command to execute without user input
      --lan-id int             The LAN ID the NIC sits on (default 1)
      --nic-dhcp               Boolean value that indicates if the NIC is using DHCP (true) or not (false) (default true)
      --nic-id string          The unique NIC Id [Required flag]
      --nic-name string        	The name of the NIC
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --remove-ip string       Remove IP
      --server-id string       The unique Server Id
      --timeout int            Timeout option [seconds] (default 60)
  -v, --verbose                Enable verbose output
      --wait                   Wait for NIC to be updated
```

## Examples

```text
ionosctl nic update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-id 2978400e-da90-405f-905e-8200d4f48158 --lan-id 2 --wait 
⧖ Waiting for request: b0361cf3-06b2-4cca-ae13-4035ace9f265
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   2       []
✔ RequestId: b0361cf3-06b2-4cca-ae13-4035ace9f265
✔ Status: Command nic update and request have been successfully executed
```

