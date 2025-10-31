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

DEFAULT CONTROLS:
Ctrl + A\tGo to the beginning of the line (Home)
Ctrl + E\tGo to the end of the line (End)
Ctrl + P\tPrevious command (Up arrow)
Ctrl + N\tNext command (Down arrow)
Ctrl + F\tForward one character
Ctrl + B\tBackward one character
Ctrl + D\tDelete character under the cursor
Ctrl + H\tDelete character before the cursor (Backspace)
Ctrl + W\tCut the word before the cursor to the clipboard
Ctrl + K\tCut the line after the cursor to the clipboard
Ctrl + U\tCut the line before the cursor to the clipboard
Ctrl + L\tClear the screen

## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
      --limit int             pagination limit: Maximum number of items to return per request (default 50)
      --no-headers            Don't print table headers when table output is used
      --offset int            pagination offset: Number of items to skip before starting to collect the results
  -o, --output string         Desired output format [text|json|api-json] (default "text")
  -p, --persist-flag-values   Persist flag values between commands
  -q, --quiet                 Quiet output
  -v, --verbose count         Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl shell
```

