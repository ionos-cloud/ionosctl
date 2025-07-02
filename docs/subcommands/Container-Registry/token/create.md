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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --expiry-date string   Expiry date of the Token
      --expiry-time string   Time until the Token expires (ex: 1y2d)
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --name string          Name of the Token (required)
      --no-headers           Use --no-headers=false to show column headers (default true)
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID (required)
      --status string        Status of the Token
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry token create --registry-id [REGISTRY-ID] --name [TOKEN-NAME]
In order to save the token to a environment variable: export [ENV-VAL-NAME]=$(ionosctl cr token create --name [TOKEN-NAME] --registry-id [REGISTRY-ID]
```

