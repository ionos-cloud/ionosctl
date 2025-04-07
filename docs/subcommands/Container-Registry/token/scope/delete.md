---
description: "Delete a token scope"
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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ScopeId DisplayName Type Actions]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -n, --scope-id int         Scope id (default -1)
  -t, --token-id string      Token ID
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry token scope delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME]
```

