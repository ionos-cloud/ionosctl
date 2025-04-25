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
  1. Global flags
  2. Environment variables
  3. Config file entries
Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.

## Options

```text
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -p, --provenance         If set, the command prints the layers of authentication sources, their order of priority, and which one was used. It also tells you if a token or username and password are being used for authentication.
  -q, --quiet              Quiet output
  -t, --timeout duration   Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait               Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl cfg whoami
ionosctl cfg whoami --provenance
```

