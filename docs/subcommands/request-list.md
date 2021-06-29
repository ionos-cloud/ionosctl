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

Use this command to list all Requests on your account.

Use flags to retrieve a list of Requests:

* sorting by the time the Request was created, starting from now in descending order, take the first N Requests: `ionosctl request list --latest N`
* sorting by method: `ionosctl request list --method REQUEST_METHOD`, where request method can be CREATE or POST, UPDATE or PATCH, PUT and DELETE
* sorting by both of the above options: `ionosctl request list --method REQUEST_METHOD --latest N`

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [RequestId CreatedDate CreatedBy Method Status Message Url Body] (default [RequestId,CreatedDate,Method,Status,Message])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
      --latest int       Show latest N Requests. If it is not set, all Requests will be printed
      --method string    Show only the Requests with this method. E.g CREATE, UPDATE, DELETE
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl request list --latest N
```

