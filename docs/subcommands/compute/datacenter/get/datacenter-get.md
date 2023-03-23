---
description: Get a Data Center
---

# DatacenterGet

## Usage

```text
ionosctl datacenter get [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a Virtual Data Center by using its ID.

Required values to run command:

* Data Center Id

## Options

```text
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --no-headers             When using text output, don't print headers
```

## Examples

```text
ionosctl datacenter get --datacenter-id DATACENTER_ID
```

