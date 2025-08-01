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

Delete a token from a registry

## Options

```text
  -a, --all                  Delete all tokens from all registries
      --all-tokens           Delete all tokens from a registry
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -t, --token-id string      Token ID
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry token delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]
```

