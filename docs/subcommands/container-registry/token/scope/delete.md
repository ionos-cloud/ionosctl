---
description: Delete a token scope
---

# ContainerRegistryTokenScopeDelete

## Usage

```text
ionosctl container-registry token scope delete [flags]
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

For `delete` command:

```text
[d rm remove]
```

## Description

Use this command to delete a token scope of a container registry. If a name is provided, the first scope with that name will be deleted. It is possible to delete all scopes by providing the --all flag.

## Options

```text
  -a, --all                  List all scopes of all tokens of a registry.
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ScopeId TokenId DisplayName Type Actions]
  -r, --registry-id string   Registry ID
  -n, --scope-id int         Scope id (default -1)
  -t, --token-id string      Token ID
```

## Examples

```text
ionosctl container-registry token scope delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME]
```

