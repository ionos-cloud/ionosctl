---
description: List Private Cross-Connects
---

# PccList

## Usage

```text
ionosctl pcc list [flags]
```

## Aliases

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of existing Private Cross-Connects available on your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl pcc list 
PccId                                  Name   Description
e2337b40-52d9-48d2-bcbc-41c5abc29d11   test   test test
4b9c6a43-a338-11eb-b70c-7ade62b52cc0   test   test
```

