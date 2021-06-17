---
description: List Requests
---

# RequestList

## Usage

```text
ionosctl request list [flags]
```

## Aliases

For `request` command:
```text
[req]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to list all Requests on your account

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [RequestId Status Message] (default [RequestId,Status,Message])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl request list
```

