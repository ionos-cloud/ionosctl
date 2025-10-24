---
description: "Create or replace a token"
---

# ContainerRegistryTokenReplace

## Usage

```text
ionosctl container-registry token replace [flags]
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

For `replace` command:

```text
[r re]
```

## Description

Create or replace a token used to access a container registry

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --expiry-date string   Expiry date of the Token
      --expiry-time string   Time until the Token expires (ex: 1y2d)
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -n, --name string          Name of the Token (required)
      --no-headers           Use --no-headers=false to show column headers (default true)
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
      --status string        Status of the Token
  -t, --token-id string      Token ID
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl container-registry token replace --name [NAME] --registry-id [REGISTRY-ID] --token-id [TOKEN-ID]
In order to save the token to a environment variable: export [ENV-VAL-NAME]=$(ionosctl cr token replace --name [TOKEN-NAME] --registry-id [REGISTRY-ID] --token-id [TOKEN-ID]
```

