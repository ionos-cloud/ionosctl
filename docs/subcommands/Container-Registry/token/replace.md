---
description: "Replace a token, regenerating its password (PUT)"
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

Replace an existing token (HTTP PUT). This overwrites the whole token with the values you pass, so any omitted property (expiry, status, and its scopes) is reset to its default - in particular, existing scopes are cleared. To keep scopes, prefer 'container-registry token update' or re-add them with 'container-registry token scope add'.

Replacing regenerates the token password, which is printed only once in the response (capture it - see the second example). The old password stops working. Set expiry with --expiry-date (absolute RFC3339) or --expiry-time (relative duration); the two are mutually exclusive.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --expiry-date string   Absolute expiry date as an RFC3339 timestamp, e.g. 2025-01-02T15:04:05Z. Mutually exclusive with --expiry-time; omit both to never expire
      --expiry-time string   Relative expiry as a duration from now, combining y/m/d/h, e.g. 1y2d. Mutually exclusive with --expiry-date
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -n, --name string          The name of the token (required)
      --no-headers           Use --no-headers=false to show column headers (default true)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   The unique ID of the registry that owns the token
      --status string        Token status: 'enabled' (usable) or 'disabled' (revoked)
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
      --token-id string      The unique ID of the token to replace
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Replace a token (prints a new, one-time password)
ionosctl container-registry token replace --registry-id REGISTRY_ID --token-id TOKEN_ID --name push-token

# Replace and capture the new password into an env var
export CR_TOKEN=$(ionosctl cr token replace --registry-id REGISTRY_ID --token-id TOKEN_ID --name ci-token --expiry-time 1y)
```

