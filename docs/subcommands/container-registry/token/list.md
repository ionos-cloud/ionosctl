---
description: List all tokens
---

# ContainerRegistryTokenList

## Usage

```text
ionosctl container-registry token list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List all tokens for your container registry

## Options

```text
  -a, --all                  List all tokens, including expired ones
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status]
  -r, --registry-id string   Registry ID
```

## Examples

```text
ionosctl container-registry token list --registry-id [REGISTRY-ID]
```

