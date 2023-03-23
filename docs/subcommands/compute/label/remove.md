---
description: Remove a Label from a Resource
---

# LabelRemove

## Usage

```text
ionosctl label remove [flags]
```

## Aliases

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a Label from a Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id
* Label Key

## Options

```text
  -a, --all                    Remove all Labels
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
      --resource-type string   Type of resource to remove labels from. Can be one of: datacenter, volume, server, snapshot, ipblock (required)
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label remove --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY
```

