---
description: "Block until a request reaches DONE (or FAILED)"
---

# RequestWait

## Usage

```text
ionosctl compute request wait [flags]
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

Use this command to block the terminal until a specific asynchronous request finishes. The CLI polls the request's status endpoint and returns once the request reaches a terminal state: it prints the request on success (DONE) or returns an error on FAILED.

This is the same polling that `--wait-for-request` performs inline on create/update/delete commands; use `request wait` when you already have a request ID and want to wait after the fact (e.g. in a script).

Use the global `--timeout` option to cap how long to wait, in seconds (default 600). If the request has not finished within the timeout, the command aborts with a timeout error.

Required values to run command:

* Request Id

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [RequestId CreatedDate CreatedBy Method Status Message Url Body Targets]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -i, --request-id string   The ID of the request to wait on (required)
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Wait up to the default 600s for a request to finish
ionosctl compute request wait --request-id REQUEST_ID

# Give up after 120 seconds
ionosctl compute request wait --request-id REQUEST_ID --timeout 120
```

