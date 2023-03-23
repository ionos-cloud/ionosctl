---
description: Get a token
---

# ContainerRegistryTokenGet

## Usage

```text
ionosctl container-registry token get [flags]
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

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve information about a single token of a container registry.

## Options

```text
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status]
  -r, --registry-id string   Registry ID
  -t, --token-id string      Token ID
```

## Examples

```text
ionosctl container-registry token get --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]
```

