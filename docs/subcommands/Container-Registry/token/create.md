---
description: "Create a new token"
---

# ContainerRegistryTokenCreate

## Usage

```text
ionosctl container-registry token create [flags]
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

For `create` command:

```text
[c]
```

## Description

Create a new token used to access a container registry

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --expiry-date string   Expiry date of the Token
      --expiry-time string   Time until the Token expires (ex: 1y2d)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -n, --name string          Name of the Token (required)
      --no-headers           Use --no-headers=false to show column headers (default true)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID (required)
      --status string        Status of the Token
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry token create --registry-id [REGISTRY-ID] --name [TOKEN-NAME]
In order to save the token to a environment variable: export [ENV-VAL-NAME]=$(ionosctl cr token create --name [TOKEN-NAME] --registry-id [REGISTRY-ID]
```

