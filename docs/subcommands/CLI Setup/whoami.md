---
description: "Tells you who you are logged in as. Use `--provenance` to debug where your credentials are being used from"
---

# ConfigWhoami

## Usage

```text
ionosctl config whoami [flags]
```

## Aliases

For `config` command:

```text
[cfg]
```

## Description

This command will tell you the email of the user you are logged in as.
You can use '--provenance' flag to see which of these sources are being used. Note that If authentication fails, this flag is set by default.
If using a token, it will use the JWT's claims payload to find out your user UUID, then use the Users API on that UUID to find out your e-mail address.
If no token is present, the command will fall back to using the username and password for authentication.

AUTHENTICATION ORDER
ionosctl uses a layered approach for authentication, prioritizing sources in this order:
  1. Flags
  2. Environment variables
  3. Config file entries
Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.

## Options

```text
<<<<<<< HEAD
  -u, --api-url string   Override default host URL. Preferred over the config file override 'auth' and env var 'IONOS_API_URL' (default "https://api.ionos.com/auth/v1")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers       Don't print table headers when table output is used
      --offset int       Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -p, --provenance       If set, the command prints the layers of authentication sources, their order of priority, and which one was used. It also tells you if a token or username and password are being used for authentication.
      --query string     JMESPath query string to filter the output
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
=======
  -u, --api-url string    Override default host URL. Preferred over the config file override 'auth' and env var 'IONOS_API_URL' (default "https://api.ionos.com/auth/v1")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -p, --provenance        If set, the command prints the layers of authentication sources, their order of priority, and which one was used. It also tells you if a token or username and password are being used for authentication.
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
>>>>>>> 8e970fd7 (remove deprecated 'D' for 'datacenter-id' only on psql)
```

## Examples

```text
ionosctl cfg whoami
ionosctl cfg whoami --provenance
```

