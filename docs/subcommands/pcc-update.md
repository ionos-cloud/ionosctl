---
description: Update a Private Cross-Connect
---

# PccUpdate

## Usage

```text
ionosctl pcc update [flags]
```

## Aliases

For `update` command:

```text
[u up]
```

## Description

Use this command to update details about a specific Private Cross-Connect. Name and description can be updated.

Required values to run command:

* Pcc Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -d, --description string   The description for the Private Cross-Connect
  -f, --force                Force command to execute without user input
  -h, --help                 help for update
  -n, --name string          The name for the Private Cross-Connect
  -o, --output string        Desired output format [text|json] (default "text")
  -i, --pcc-id string        The unique Private Cross-Connect Id (required)
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for Request for Private Cross-Connect update [seconds] (default 60)
  -v, --verbose              see step by step process when running a command
  -w, --wait-for-request     Wait for the Request for Private Cross-Connect update to be executed
```

## Examples

```text
ionosctl pcc update --pcc-id PCC_ID --description DESCRIPTION
```

