---
description: "Interactive shell - BETA"
---

# Shell

## Usage

```text
ionosctl shell [flags]
```

## Description

The ionosctl shell command launches an interactive shell environment, enabling a more dynamic and intuitive way to interact with the ionosctl CLI.
This shell is designed to enhance your command-line experience with advanced features and customizations, powered by the comptplus library.

CUSTOM CONTROLS: (your usual shell controls might not work)
- SHIFT + LEFT/RIGHT: Quickly navigate words left/right
- SHIFT + UP/DOWN: Quickly navigate to the beginning/end of the line
- SHIFT + DELETE: Delete previous word

## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --no-headers            Don't print table headers when table output is used
  -o, --output string         Desired output format [text|json|api-json] (default "text")
  -p, --persist-flag-values   Persist flag values between commands
  -q, --quiet                 Quiet output
  -v, --verbose               Print step-by-step process when running command
```

## Examples

```text
ionosctl shell
```

