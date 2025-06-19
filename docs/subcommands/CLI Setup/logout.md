---
description: "Remove credentials from your YAML config (and purge old JSON)"
---

# ConfigLogout

## Usage

```text
ionosctl config logout [flags]
```

## Aliases

For `config` command:

```text
[cfg]
```

## Description

This 'Quality of Life' command will:

  1. Clear out any sensitive fields in your YAML config.
  2. Afterwards, detect and optionally delete any legacy config.json alongside it.

You can skip the YAML logout and **only** purge the old JSON with:

    ionosctl logout --only-purge-old

AUTHENTICATION ORDER
ionosctl uses a layered approach for authentication, prioritizing sources in this order:
  1. Flags
  2. Environment variables
  3. Config file entries
Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.

## Options

```text
  -c, --config string   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force           Force command to execute without user input
  -h, --help            Print usage
      --no-headers      Don't print table headers when table output is used
  -o, --output string   Desired output format [text|json|api-json] (default "text")
  -q, --quiet           Quiet output
  -v, --verbose         Print step-by-step process when running command
```

## Examples

```text
ionosctl logout
ionosctl logout --only-purge-old
```

