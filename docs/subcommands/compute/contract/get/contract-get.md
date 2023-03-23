---
description: Get information about the Contract Resources on your account
---

# ContractGet

## Usage

```text
ionosctl contract get [flags]
```

## Aliases

For `contract` command:

```text
[c]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about the Contract Resources on your account. Use `--resource-limits` flag to see specific Contract Resources Limits.

## Options

```text
  -D, --depth int32              Controls the detail depth of the response objects. Max depth is 10.
      --no-headers               When using text output, don't print headers
      --resource-limits string   Specify Resource Limits to see details about it
```

## Examples

```text
ionosctl contract get --resource-limits [ CORES|RAM|HDD|SSD|IPS|K8S ]
```

