---
description: Create a new token
---

# ContainerRegistryTokenCreate

## Usage

```text
ionosctl container-registry token create [flags]
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

For `create` command:

```text
[c]
```

## Description

Create a new token used to access a container registry

## Options

```text
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status]
      --expiry-date string   Expiry date of the Token
      --expiry-time string   Time until the Token expires (ex: 1y2d)
      --name string          Name of the Token (required)
  -r, --registry-id string   Registry ID (required)
      --status string        Status of the Token
```

## Examples

```text
ionosctl container-registry token create --registry-id [REGISTRY-ID] --name [TOKEN-NAME]
```

