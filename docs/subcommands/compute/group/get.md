---
description: Get a Group
---

# GroupGet

## Usage

```text
ionosctl group get [flags]
```

## Aliases

For `group` command:

```text
[g]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Group.

Required values to run command:

* Group Id

## Options

```text
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
  -i, --group-id string   The unique Group Id (required)
      --no-headers        When using text output, don't print headers
```

## Examples

```text
ionosctl group get --group-id GROUP_ID
```

