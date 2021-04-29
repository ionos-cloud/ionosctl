---
description: Wait a Request
---

# Wait

## Usage

```text
ionosctl request wait [flags]
```

## Description

Use this command to wait for a specified Request to execute.

You can specify a timeout for the action to be executed using `--timeout` option.

Required values to run command:

* Request Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [RequestId,Status,Message])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force               Force command to execute without user input
  -h, --help                help for wait
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
      --request-id string   The unique Request Id (required)
      --timeout int         Timeout option waiting for request [seconds] (default 60)
```

## Examples

```text
ionosctl request wait --request-id 20333e60-d65c-4a95-846b-08c48b871186 
RequestId                              Status   Message
20333e60-d65c-4a95-846b-08c48b871186   DONE     Request has been successfully executed
```

