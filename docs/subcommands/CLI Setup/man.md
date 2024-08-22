---
description: "Generate manpages for ionosctl"
---

# Man

## Usage

```text
ionosctl man [flags]
```

## Aliases

For `man` command:

```text
[manpages]
```

## Description

WARNING: This command is only supported on Linux.

The 'man' command allows you to generate manpages for ionosctl in a given directory. By default, the manpages will be compressed using gzip, but you can skip this step by using the '--skip-compression' flag.

## Options

```text
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --skip-compression    Skip compressing manpages with gzip, just generate them
      --target-dir string   Target directory where manpages will be generated. Must be an absolute path (default "/tmp/ionosctl-man")
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl man --target-dir /tmp/ionosctl-man
```

