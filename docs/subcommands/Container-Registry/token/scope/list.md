---
description: "Get a token scopes"
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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ScopeId TokenId DisplayName Type Actions]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -t, --token-id string      Token ID
  -v, --verbose count        Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]
```

