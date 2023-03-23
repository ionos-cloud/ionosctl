---
description: Get a list of Peers from a Private Cross-Connect
---

# PccPeersList

## Usage

```text
ionosctl pcc peers list [flags]
```

## Aliases

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of Peers from a Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
      --no-headers      When using text output, don't print headers
      --pcc-id string   The unique Private Cross-Connect Id (required)
```

## Examples

```text
ionosctl pcc peers list --pcc-id PCC_ID
```

