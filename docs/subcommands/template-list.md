---
description: List Templates
---

# TemplateList

## Usage

```text
ionosctl template list [flags]
```

## Aliases

For `template` command:

```text
[tpl]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of available public Templates.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [TemplateId Name Cores Ram StorageSize] (default [TemplateId,Name,Cores,Ram,StorageSize])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          see step by step process when running a command
```

## Examples

```text
ionosctl template list
```

