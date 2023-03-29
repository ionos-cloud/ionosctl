---
description: Get all Resources of a Type or a specific Resource Type
---

# ResourceGet

## Usage

```text
ionosctl resource get [flags]
```

## Aliases

For `resource` command:

```text
[res]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get all Resources of a Type or a specific Resource Type using its Type and ID.

Required values to run command:

* Type

## Options

```text
      --no-headers           When using text output, don't print headers
  -i, --resource-id string   The ID of the specific Resource to retrieve information about
      --type string          The specific Type of Resources to retrieve information about (required)
```

## Examples

```text
ionosctl resource get --resource-type ipblock
```

