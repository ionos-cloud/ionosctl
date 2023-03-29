---
description: Create or replace a token
---

# ContainerRegistryTokenReplace

## Usage

```text
ionosctl container-registry token replace [flags]
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

For `replace` command:

```text
[r re]
```

## Description

Create or replace a token used to access a container registry

## Options

```text
      --expiry-date string   Expiry date of the Token
      --expiry-time string   Time until the Token expires (ex: 1y2d)
      --name string          Name of the Token (required)
  -r, --registry-id string   Registry ID
      --status string        Status of the Token
  -t, --token-id string      Token ID
```

## Examples

```text
ionosctl container-registry token replace --name [NAME] --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]
```

