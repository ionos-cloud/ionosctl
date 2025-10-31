---
description: "Create a new Token"
---

# TokenGenerate

## Usage

```text
ionosctl token generate [flags]
```

## Aliases

For `generate` command:

```text
[create]
```

## Description

Use this command to generate a new Token. Only the JSON Web Token, associated with user credentials, will be displayed.

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'auth' and env var 'IONOS_API_URL' (default "https://api.ionos.com/auth/v1")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --contract int     Users with multiple contracts can provide the contract number, for which the token is generated
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers       Don't print table headers when table output is used
      --offset int       Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
      --ttl string       Token Time to Live (TTL). Accepted formats: Y, M, D, h, m, s. Hybrids are also allowed (e.g. 1m30s). Min: 60s (1m) Max: 31536000s (1Y)
                         NOTE: Any values that do not match the format will be ignored. (default "1Y")
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl token generate
```

