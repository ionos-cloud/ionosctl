---
description: "Get a specified Token"
---

# TokenGet

## Usage

```text
ionosctl token get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a Token by using its ID.

Required values to run command:

* Token Id

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'auth' and env var 'IONOS_API_URL' (default "https://api.ionos.com/auth/v1")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --contract int      Users with multiple contracts must provide the contract number, for which the token information is displayed
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -t, --token string      The contents of a Token (required)
  -i, --token-id string   The unique Key ID of a Token (required)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl token get --token-id TOKEN_ID

ionosctl token get --token TOKEN
```

