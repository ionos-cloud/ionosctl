---
description: Update a token's properties
---

# ContainerRegistryTokenUpdate

## Usage

```text
ionosctl container-registry token update [flags]
```

## Aliases

For `container-registry` command:

```text
[cr contreg cont-reg]
```

For `token` command:

```text
[t tokens]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a token's properties. You can update the token's expiry date and status.

## Options

```text
      --expiry-date string   Expiry date of the Token
      --expiry-time string   Time until the Token expires (ex: 1y2d)
  -r, --registry-id string   Registry ID
      --status string        Status of the Token
  -t, --token-id string      Token ID
```

## Examples

```text
ionosctl container-registry token update --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --expiry-date [EXPIRY-DATE] --status [STATUS]
```

