---
description: "Get a Cross-Connect"
---

# PccGet

## Usage

```text
ionosctl pcc get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a specific Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32      Controls the detail depth of the response objects. Max depth is 10.
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -i, --pcc-id string    The unique Cross-Connect Id (required)
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl pcc get --pcc-id PCC_ID
```

