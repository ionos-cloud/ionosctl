---
description: Get a Label
---

# LabelGet

## Usage

```text
ionosctl label get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Label from a specified Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id
* Label Key

## Options

```text
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
      --no-headers             When using text output, don't print headers
      --resource-type string   Type of resource to get labels from. Can be one of: datacenter, volume, server, snapshot, ipblock (required)
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label get --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY
```

