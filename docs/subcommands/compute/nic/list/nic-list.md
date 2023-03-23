---
description: List NICs
---

# NicList

## Usage

```text
ionosctl nic list [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of NICs on your account.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [name mac ips dhcp lan firewallActive firewallType deviceNumber pciSlot vnet]
* filter by metadata: [etag createdDate createdBy createdByUserId lastModifiedDate lastModifiedBy lastModifiedByUserId state]

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -F, --filters strings        Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             When using text output, don't print headers
      --order-by string        Limits results to those containing a matching value for a specific property
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl nic list --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

