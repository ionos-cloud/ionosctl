---
description: "Update a token's properties"
---

# ContainerRegistryTokenUpdate

## Usage

```text
ionosctl container-registry token update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update a token's properties. You can update the token's expiry date and status.

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
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
      --status string        Status of the Token
  -t, --token-id string      Token ID
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry token update --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --expiry-date [EXPIRY-DATE] --status [STATUS]
```

