---
description: Firewall Rule Operations
---

# FirewallRule

## Usage

```text
ionosctl firewallrule [command]
```

## Description

The sub-commands of `ionosctl firewallrule` allow you to create, list, get, update, delete Firewall Rules.

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
      --force                  Force command to execute without user input
  -h, --help                   help for firewallrule
      --nic-id string          The unique NIC Id [Required flag]
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id [Required flag]
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl firewallrule create](create.md) | Create a Firewall Rule |
| [ionosctl firewallrule delete](delete.md) | Delete a FirewallRule |
| [ionosctl firewallrule get](get.md) | Get a Firewall Rule |
| [ionosctl firewallrule list](list.md) | List Firewall Rules |
| [ionosctl firewallrule update](update.md) | Update a FirewallRule |

