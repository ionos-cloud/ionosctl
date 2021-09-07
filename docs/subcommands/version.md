---
description: Show the current version of CLI and SDK
---

# Version

## Usage

```text
ionosctl version [flags]
```

## Description

The `ionosctl version` command displays the current version of the ionosctl software, the SDK GO version and the latest Github release.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
      --updates          Check for latest updates for CLI
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl version

ionosctl version --updates
```

