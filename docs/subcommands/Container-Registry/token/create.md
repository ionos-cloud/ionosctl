---
description: "Create a token to authenticate against a registry"
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

Create a new token used to authenticate 'docker login'/'docker push'/'docker pull' against a registry.

The generated password is printed only once, in this response, and cannot be retrieved afterwards - capture it now (see the second example). A freshly created token has no scopes yet, so it cannot pull or push until you grant scopes with 'container-registry token scope add'.

Set an expiry with either --expiry-date (an absolute RFC3339 timestamp) or --expiry-time (a relative duration such as 1y2d); the two are mutually exclusive. Omit both for a token that never expires. Use --status disabled to create the token pre-revoked.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'containerregistry' and env var 'IONOS_API_URL' (default "https://api.ionos.com/containerregistries")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [TokenId DisplayName ExpiryDate CredentialsUsername CredentialsPassword Status RegistryId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --expiry-date string   Absolute expiry date as an RFC3339 timestamp, e.g. 2025-01-02T15:04:05Z. Mutually exclusive with --expiry-time; omit both to never expire
      --expiry-time string   Relative expiry as a duration from now, combining y (years), m (months), d (days), h (hours), e.g. 1y2d or 6m. Mutually exclusive with --expiry-date
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -n, --name string          The name of the token, unique within the registry (required)
      --no-headers           Use --no-headers=false to show column headers (default true)
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -r, --registry-id string   The unique ID of the registry that will own this token (required)
      --status string        Initial status of the token: 'enabled' (usable) or 'disabled' (revoked). Defaults to enabled
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a token (note the printed password - it is shown only once)
ionosctl container-registry token create --registry-id REGISTRY_ID --name push-token

# Create a token expiring in 1 year and 2 days, capturing the password into an env var
export CR_TOKEN=$(ionosctl cr token create --registry-id REGISTRY_ID --name ci-token --expiry-time 1y2d)
```

