---
description: Get a Private Cross-Connect
---

# PccGet

## Usage

```text
ionosctl pcc get [flags]
```

## Description

Use this command to retrieve details about a specific Private Cross-Connect.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [PccId,Name,Description,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for get
  -o, --output string    Desired output format [text|json] (default "text")
      --pcc-id string    The unique Private Cross-Connect Id (required)
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl pcc get --pcc-id e2337b40-52d9-48d2-bcbc-41c5abc29d11 
PccId                                  Name   Description
e2337b40-52d9-48d2-bcbc-41c5abc29d11   test   test test
```

