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
  -A, --all                Delete the Tokens under your account (required)
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --contract int       Users with multiple contracts must provide the contract number, for which the tokens are deleted
  -C, --current            Delete the Token that is currently used. This requires a token to be set for authentication via environment variable IONOS_TOKEN or via config file (required)
  -E, --expired            Delete the Tokens that are currently expired (required)
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
      --timeout duration   Timeout for waiting for resource to reach desired state (default 1m0s)
  -t, --token string       The contents of a Token (required)
  -i, --token-id string    The unique Key ID of a Token (required)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait               Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl token delete --token-id TOKEN_ID

ionosctl token delete --token TOKEN

ionosctl token delete --expired

ionosctl token delete --current

ionosctl token delete --all
```

