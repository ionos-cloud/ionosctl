---
description: Delete a FirewallRule
---

# FirewallruleDelete

## Usage

```text
ionosctl firewallrule delete [flags]
```

## Aliases

For `firewallrule` command:

```text
[f fr firewall]
```

For `delete` command:

```text
[d]
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
  -a, --all                      Delete all the Firewalls.
      --datacenter-id string     The unique Data Center Id (required)
      --depth int32              Controls the detail depth of the response objects. Max depth is 10.
  -i, --firewallrule-id string   The unique FirewallRule Id (required)
      --nic-id string            The unique NIC Id (required)
      --server-id string         The unique Server Id (required)
  -t, --timeout int              Timeout option for Request for Firewall Rule deletion [seconds] (default 60)
  -w, --wait-for-request         Wait for Request for Firewall Rule deletion to be executed
```

## Examples

```text
ionosctl firewallrule delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID
```

