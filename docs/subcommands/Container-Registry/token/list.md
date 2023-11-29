---
description: "List all tokens"
---

# ContainerRegistryTokenList

## Usage

```text
ionosctl container-registry token list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List all tokens for your container registry

## Options

```text
  -a, --all                  List all tokens, including expired ones
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -r, --registry-id string   Registry ID
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl container-registry token list --registry-id [REGISTRY-ID]
```

