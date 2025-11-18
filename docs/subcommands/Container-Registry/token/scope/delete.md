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
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [ScopeId DisplayName Type Actions]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -n, --scope-id int         Scope id (default -1)
  -t, --token-id string      Token ID
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry token scope delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME]
```

