---
description: Delete a token
---

# ContainerRegistryTokenDelete

## Usage

```text
ionosctl container-registry token delete [flags]
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

For `delete` command:

```text
[d del rm]
```

## Description

Delete a token from a registry

## Options

```text
  -a, --all                  Delete all tokens from all registries
      --all-tokens           Delete all tokens from a registry
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status]
  -r, --registry-id string   Registry ID
  -t, --token-id string      Token ID
```

## Examples

```text
ionosctl container-registry token delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]
```

