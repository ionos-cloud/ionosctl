---
description: "Delete one or multiple Tokens"
---

# TokenDelete

## Usage

```text
ionosctl token delete [flags]
```

## Aliases

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Token by token Id or multiple Tokens (based on a criteria: CURRENT, EXPIRED, ALL) from your account. With parameter values ALL and EXPIRED, 'Basic Authentication' or 'Token Authentication' tokens with valid credentials must be encapsulated in the header. With value CURRENT, only the 'Token Authentication' with valid credentials is required.

Required values to run command:

* Token Id/Token/CURRENT/EXPIRED/ALL

## Options

```text
  -a, --all               Delete the Tokens under your account (required)
  -u, --api-url string    Override default host URL. Preferred over the config file override 'auth' and env var 'IONOS_API_URL' (default "https://api.ionos.com/auth/v1")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --contract int      Users with multiple contracts must provide the contract number, for which the tokens are deleted
      --current           Delete the Token that is currently used. This requires a token to be set for authentication via environment variable IONOS_TOKEN or via config file (required)
      --expired           Delete the Tokens that are currently expired (required)
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -t, --token string      The contents of a Token (required)
  -i, --token-id string   The unique Key ID of a Token (required)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl token delete --token-id TOKEN_ID

ionosctl token delete --token TOKEN

ionosctl token delete --expired

ionosctl token delete --current

ionosctl token delete --all
```

