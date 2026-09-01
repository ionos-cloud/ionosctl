---
description: "List recent provisioning requests and their status"
---

# RequestList

## Usage

```text
ionosctl compute request list [flags]
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

Use this command to list the asynchronous requests on your account, most useful for reviewing recent activity or finding a failed request. The Status column shows QUEUED / RUNNING / DONE / FAILED, and Targets shows the resources each request acted on.

You can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.
Available Filters:
* filter by property: [method headers body url]
* filter by metadata: [createdDate createdBy etag requestStatus status message etag]

## Options

```text
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [RequestId CreatedDate CreatedBy Method Status Message Url Body Targets]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --latest int        Show only the N most recent requests (sorted by creation date, newest first). If unset, all requests are printed (DEPRECATED: Use --filters --order-by --max-results options instead!)
      --limit int         Maximum number of items to return per request (default 50)
      --method string     Show only requests with this HTTP method. Accepts POST/PUT/PATCH/DELETE, or the aliases CREATE (=POST) and UPDATE (=PUT+PATCH) (DEPRECATED: Use --filters --order-by --max-results options instead!)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# List all requests
ionosctl compute request list

# Show only DONE/FAILED status and target columns
ionosctl compute request list --cols RequestId,Method,Status,Message,Targets
```

