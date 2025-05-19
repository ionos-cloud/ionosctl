---
description: "Print your config file's path"
---

# ConfigLocation

## Usage

```text
ionosctl config location [flags]
```

## Aliases

For `config` command:

```text
[cfg]
```

For `location` command:

```text
[location loc]
```

## Description

Print your config file's path

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl cfg loc
```

