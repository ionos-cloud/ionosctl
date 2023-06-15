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

Tells you who you are logged in as. Use `--provenance` to debug where your credentials are being used from

## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
  -o, --output string         Desired output format [text|json] (default "text")
  -p, --provenance userdata   If set, prints a JSON object which explains the source of each configuration variable. All-caps rules represent env vars. Rules prefixed with userdata represent vars read from config. Otherwise, the rules represent flags.
  -q, --quiet                 Quiet output
  -v, --verbose               Print step-by-step process when running command
```

## Examples

```text
ionosctl whoami
```

