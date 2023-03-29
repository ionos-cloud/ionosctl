---
description: Get a token scopes
---

# ContainerRegistryTokenScopeList

## Usage

```text
ionosctl container-registry token scope list [flags]
```

## Aliases

For `token` command:

```text
[t tokens]
```

For `scope` command:

```text
[s scopes]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list all scopes of a token of a container registry.

## Options

```text
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ScopeId TokenId DisplayName Type Actions]
  -r, --registry-id string   Registry ID
  -t, --token-id string      Token ID
```

## Examples

```text
ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]
```

