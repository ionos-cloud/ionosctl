---
description: "Delete a FlowLog from a NIC"
---

# FlowlogDelete

## Usage

```text
ionosctl compute flowlog delete [flags]
```

## Aliases

For `flowlog` command:

```text
[fl]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a Flow Log from a NIC, stopping traffic-metadata capture. Already-delivered log files remain in the Object Storage bucket and are not removed. Delete the Flow Log before deleting the bucket it writes to.

Use `--wait` (`-w`) to wait for the request to complete. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* FlowLog Id

## Options

```text
  -a, --all                    Delete all Flow Logs on the specified NIC. Mutually exclusive with --flowlog-id
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique ID of the Virtual Data Center that holds the server and NIC (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -i, --flowlog-id string      The unique ID of the Flow Log to delete (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --nic-id string          The unique ID of the NIC the Flow Log belongs to (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique ID of the server that owns the NIC (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Delete a single Flow Log, skipping the confirmation prompt and waiting for completion
ionosctl compute flowlog delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --flowlog-id FLOWLOG_ID -f -w

# Delete every Flow Log on a NIC
ionosctl compute flowlog delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --all
```

