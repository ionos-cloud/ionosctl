---
description: "Update a token's expiry or status (PATCH)"
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

Update an existing token in place (HTTP PATCH). Only the fields you pass are changed; the token's scopes and password are preserved (unlike 'replace', which regenerates the password and clears scopes).

Use --status disabled to revoke a token without deleting it (and enabled to re-activate it), or change its expiry with --expiry-date (absolute RFC3339) / --expiry-time (relative duration). To change what the token may access, use 'container-registry token scope'.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --expiry-date string   New absolute expiry date as an RFC3339 timestamp, e.g. 2025-01-02T15:04:05Z
      --expiry-time string   New expiry as a duration from now, combining y/m/d/h, e.g. 1y2d
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
      --status string        Token status: 'enabled' (usable) or 'disabled' (revoked)
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
      --token-id string      The unique ID of the token to update
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Revoke a token without deleting it
ionosctl container-registry token update --registry-id REGISTRY_ID --token-id TOKEN_ID --status disabled

# Extend a token's expiry to an absolute date
ionosctl container-registry token update --registry-id REGISTRY_ID --token-id TOKEN_ID --expiry-date 2026-01-01T00:00:00Z
```

