---
description: Request Operations
---

# Request

## Usage

```text
ionosctl request [command]
```

## Description

The sub-commands of `ionosctl request` allow you to see information about requests on your account.
With the `ionosctl request` command, you can list, get or wait for a Request.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [RequestId,Status,Message])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl-config.json")
  -h, --help             help for request
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Enable verbose output
```

## See also

* [ionosctl](../)
* [ionosctl request get](get.md)
* [ionosctl request list](list.md)
* [ionosctl request wait](wait.md)

