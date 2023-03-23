---
description: Get a Volume
---

# VolumeGet

## Usage

```text
ionosctl volume get [flags]
```

## Aliases

For `volume` command:

```text
[v vol]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about a Volume using its ID.

Required values to run command:

* Data Center Id
* Volume Id

## Options

```text
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --no-headers             When using text output, don't print headers
  -i, --volume-id string       The unique Volume Id (required)
```

## Examples

```text
ionosctl volume get --datacenter-id DATACENTER_ID --volume-id VOLUME_ID
```

