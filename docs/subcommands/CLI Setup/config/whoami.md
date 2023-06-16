---
description: Tells you who you are logged in as. Use `--provenance` to debug where your credentials are being used from
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

'ionosctl whoami' will tell you the email of the user you are logged in as.
ionosctl prioritizes different sources of authentication, in the following order:
  1. Token provided as a global flag
  2. Token, Username and Password provided as environment variables
  3. Token, Username and Password read from the config file
Additionally, the token has priority over username & password in the layer it is originating from. (meaning that a token in your config file will be overridden by IONOS_USERNAME & IONOS_PASSWORD)
You can use '--provenance' flag to see which of these sources are being used.
If using a token, it will use the JWT's claims payload to find out your user UUID, then use the Users API on that UUID to find out your e-mail address.
If no token is present, the command will fall back to using the username and password for authentication.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -p, --provenance       If set, the command prints the layers of authentication sources, their order of priority, and which one was used. It also tells you if a token or username and password are being used for authentication.
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl cfg whoami
ionosctl cfg whoami --provenance
```

