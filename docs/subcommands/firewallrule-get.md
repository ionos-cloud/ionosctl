---
description: Get a Firewall Rule
---

# FirewallruleGet

## Usage

```text
ionosctl firewallrule get [flags]
```

## Description

Use this command to retrieve information of a specified Firewall Rule.

Required values to run command:

* Data Center Id
* Server Id
*Nic Id
* FirewallRule Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
      --firewallrule-id string   The unique FirewallRule Id (required)
  -f, --force                    Force command to execute without user input
  -F, --format strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,State])
  -h, --help                     help for get
      --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id (required)
```

## Examples

```text
ionosctl firewallrule get --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --firewallrule-id f537ff0e-8b2c-4ce6-8a92-297a5ad08ca1 
FirewallRuleId                         Name        Protocol   PortRangeStart   PortRangeEnd   State
f537ff0e-8b2c-4ce6-8a92-297a5ad08ca1   test        TCP        80               80             AVAILABLE
```

