---
description: Get a Request
---

# RequestGet

## Usage

```text
ionosctl request get [flags]
```

## Aliases

For `request` command:
```text
[req]
```

## Description

Use this command to get information about a specified Request.

Required values to run command:

* Request Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [RequestId Status Message] (default [RequestId,Status,Message])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                help for get
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -i, --request-id string   The unique Request Id (required)
```

## Examples

```text
ionosctl request get --request-id 20333e60-d65c-4a95-846b-08c48b871186 
RequestId                              Status   Message
20333e60-d65c-4a95-846b-08c48b871186   DONE     Request has been successfully executed
```

