---
description: Get a Firewall Rule
---

# FirewallruleGet

## Usage

```text
ionosctl firewallrule get [flags]
```

## Aliases

For `firewallrule` command:

```text
[f fr firewall]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information of a specified Firewall Rule.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* FirewallRule Id

## Options

```text
      --datacenter-id string     The unique Data Center Id (required)
      --depth int32              Controls the detail depth of the response objects. Max depth is 10.
  -i, --firewallrule-id string   The unique FirewallRule Id (required)
      --nic-id string            The unique NIC Id (required)
      --no-headers               When using text output, don't print headers
      --server-id string         The unique Server Id (required)
```

## Examples

```text
ionosctl firewallrule get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID
```

