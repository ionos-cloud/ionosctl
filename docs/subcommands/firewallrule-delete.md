---
description: Delete a FirewallRule
---

# FirewallruleDelete

## Usage

```text
ionosctl firewallrule delete [flags]
```

## Description

Use this command to delete a specified Firewall Rule from a Virtual Data Center.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Firewall Rule Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
      --firewallrule-id string   The unique FirewallRule Id (required)
  -f, --force                    Force command to execute without user input
  -F, --format strings           Columns to be printed in the standard output. Example: --cols "ResourceId,Name" (default [FirewallRuleId,Name,Protocol,PortRangeStart,PortRangeEnd,State])
  -h, --help                     help for delete
      --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id (required)
  -t, --timeout int              Timeout option for Request for Firewall Rule deletion [seconds] (default 60)
  -w, --wait-for-request         Wait for Request for Firewall Rule deletion to be executed
```

## Examples

```text
ionosctl firewallrule delete --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --firewallrule-id e7c4e91a-d3e3-42db-bfb1-2d5e9ebc952b 
Warning: Are you sure you want to delete firewall rule (y/N) ? 
y
RequestId: 481b6e7c-0c31-4395-81e4-36fad877b77b
Status: Command firewallrule delete has been successfully executed
```

