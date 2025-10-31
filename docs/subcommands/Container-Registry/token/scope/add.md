---
description: "Add scopes to a token"
---

# ContainerRegistryTokenScopeAdd

## Usage

```text
ionosctl container-registry token scope add [flags]
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

For `add` command:

```text
[a ad]
```

## Description

Use this command to add scopes to a token of a container registry.

## Options

```text
  -a, --actions strings      Scope actions (required)
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ScopeId DisplayName Type Actions]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            pagination limit: Maximum number of items to return per request (default 50)
  -n, --name string          Scope name (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID (required)
  -t, --token-id string      Token ID
  -y, --type string          Scope type (required)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME] --actions [SCOPE-ACTIONS], --type [SCOPE-TYPE]
```

