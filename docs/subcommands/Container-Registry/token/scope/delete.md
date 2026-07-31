---
description: "Remove a scope from a token"
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

Remove a scope from a token, narrowing what the token may access (removing all scopes leaves the token unable to pull or push).

Identify the scope to remove by its zero-based --scope-id (the ScopeId column from 'scope list'), or pass --all to remove every scope from the token. Internally the token is deleted and re-created with the remaining scopes; its expiry, status and name are preserved.

## Options

```text
  -a, --all                  Remove every scope from the token
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
  -r, --registry-id string   The unique ID of the registry that owns the token
  -n, --scope-id int         Zero-based index of the scope to remove (the ScopeId shown by 'scope list') (default -1)
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
      --token-id string      The unique ID of the token to remove a scope from
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Remove the scope at index 0
ionosctl container-registry token scope delete --registry-id REGISTRY_ID --token-id TOKEN_ID --scope-id 0

# Remove all scopes from a token
ionosctl container-registry token scope delete --registry-id REGISTRY_ID --token-id TOKEN_ID --all
```

