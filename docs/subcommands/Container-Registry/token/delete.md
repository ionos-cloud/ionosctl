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
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           When using text output, don't print headers
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

