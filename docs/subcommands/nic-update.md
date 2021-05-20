---
description: Update a NIC
---

# NicUpdate

## Usage

```text
ionosctl nic update [flags]
```

## Aliases

For `nic` command:
```text
[n]
```

## Description

Use this command to update the configuration of a specified NIC. Some restrictions are in place: The primary address of a NIC connected to a Load Balancer can only be changed by changing the IP of the Load Balancer. You can also add additional reserved, public IPs to the NIC.

The user can specify and assign private IPs manually. Valid IP addresses for private networks are 10.0.0.0/8, 172.16.0.0/12 or 192.168.0.0/16.

The value for firewallActive can be toggled between true and false to enable or disable the firewall. When the firewall is enabled, incoming traffic is filtered by all the firewall rules configured on the NIC. When the firewall is disabled, then all incoming traffic is routed directly to the NIC and any configured firewall rules are ignored.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* NIC Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --dhcp                   Boolean value that indicates if the NIC is using DHCP (true) or not (false) (default true)
  -f, --force                  Force command to execute without user input
  -F, --format strings         Collection of fields to be printed on output (default [NicId,Name,Dhcp,LanId,Ips,State])
  -h, --help                   help for update
      --ips strings            IPs assigned to the NIC
      --lan-id int             The LAN ID the NIC sits on (default 1)
  -n, --name string            The name of the NIC
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
  -t, --timeout int            Timeout option for Request for NIC update [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for NIC update to be executed
```

## Examples

```text
ionosctl nic update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-id 2978400e-da90-405f-905e-8200d4f48158 --lan-id 2 --wait-for-request
1.2s Waiting for request... DONE
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   2       []
RequestId: b0361cf3-06b2-4cca-ae13-4035ace9f265
Status: Command nic update & wait have been successfully executed
```

