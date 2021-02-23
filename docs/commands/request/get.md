---
description: Get a Request
---

# Get

## Usage

```text
ionosctl request get [flags]
```

## Description

Use this command to get information about a specified Request.

Required values to run command:
- Request Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [RequestId,Status,Message])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                help for get
      --ignore-stdin        Force command to execute without user input
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
      --request-id string   The unique Request Id. [Required flag]
  -v, --verbose             Enable verbose output
```

## Examples

```text
ionosctl request get --request-id 20333e60-d65c-4a95-846b-08c48b871186 
RequestId                              Status   Message
20333e60-d65c-4a95-846b-08c48b871186   DONE     Request has been successfully executed
```

