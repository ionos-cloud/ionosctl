---
description: Show the current version
---

# Version

## Usage

```text
ionosctl version [flags]
```

## Description

The `ionosctl version` command displays the current version of the ionosctl software and the latest Github release.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for version
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
      --updates          Check for latest updates for CLI
  -v, --verbose          see step by step process when running a command
```

## Examples

```text
ionosctl version
```

