---
description: "Delete a token"
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

Delete a token. Once deleted, its password stops working immediately for 'docker login'; any client using it must switch to another token.

Scope of deletion (choose one): a single token (--registry-id + --token-id), every token of one registry (--registry-id + --all-tokens), or every token in the entire contract (--all).

## Options

```text
  -a, --all                  Delete every token in every registry of the contract
      --all-tokens           Delete every token of the registry named by --registry-id
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
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
  -r, --registry-id string   The unique ID of the registry that owns the token(s)
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
      --token-id string      The unique ID of the token to delete
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Delete one token
ionosctl container-registry token delete --registry-id REGISTRY_ID --token-id TOKEN_ID

# Delete every token of a registry
ionosctl container-registry token delete --registry-id REGISTRY_ID --all-tokens
```

