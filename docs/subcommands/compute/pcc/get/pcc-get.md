---
description: Get a Private Cross-Connect
---

# PccGet

## Usage

```text
ionosctl pcc get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -D, --depth int32     Controls the detail depth of the response objects. Max depth is 10.
      --no-headers      When using text output, don't print headers
  -i, --pcc-id string   The unique Private Cross-Connect Id (required)
```

## Examples

```text
ionosctl pcc get --pcc-id PCC_ID
```

