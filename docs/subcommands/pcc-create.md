---
description: Create a Private Cross-Connect
---

# PccCreate

## Usage

```text
ionosctl pcc create [flags]
```

## Aliases

For `create` command:

```text
[c]
```

## Description

Use this command to create a Private Cross-Connect. You can specify the name and the description for the Private Cross-Connect.

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [PccId Name Description State] (default [PccId,Name,Description,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -d, --description string   The description for the Private Cross-Connect
  -f, --force                Force command to execute without user input
  -h, --help                 help for create
  -n, --name string          The name for the Private Cross-Connect (default "Cross Connect")
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for Request for Private Cross-Connect creation [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Private Cross-Connect creation to be executed
```

## Examples

```text
ionosctl pcc create --name NAME --description DESCRIPTION --wait-for-request
```

