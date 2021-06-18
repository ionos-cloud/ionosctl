---
description: Wait a Request
---

# RequestWait

## Usage

```text
ionosctl request wait [flags]
```

## Aliases

For `request` command:

```text
[req]
```

For `wait` command:

```text
[w]
```

## Description

Use this command to wait for a specified Request to execute.

You can specify a timeout for the Request to be executed using `--timeout` option.

Required values to run command:

* Request Id

## Options

```text
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [RequestId CreatedDate CreatedBy Method Status Message Url Body] (default [RequestId,CreatedDate,Method,Status,Message])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                help for wait
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -i, --request-id string   The unique Request Id (required)
  -t, --timeout int         Timeout option waiting for Request [seconds] (default 60)
```

## Examples

```text
ionosctl request wait --request-id REQUEST_ID
```

