---
description: Add a Label to a Resource
---

# LabelAdd

## Usage

```text
ionosctl label add [flags]
```

## Aliases

For `add` command:

```text
[a]
```

## Description

Use this command to add a Label to a specific Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id
* Label Key
* Label Value

## Options

```text
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
      --label-value string     The unique Label Value (required)
      --resource-type string   Type of resource to add labels to. Can be one of: datacenter, volume, server, snapshot, ipblock (required)
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label add --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID  --label-key LABEL_KEY --label-value LABEL_VALUE

ionosctl label add --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY --label-value LABEL_VALUE
```

