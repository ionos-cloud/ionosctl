---
description: "Delete (release) an IpBlock"
---

# IpblockDelete

## Usage

```text
ionosctl compute ipblock delete [flags]
```

## Aliases

For `ipblock` command:

```text
[ip ipb]
```

For `delete` command:

```text
[d]
```

## Description

Release a reserved IpBlock, returning all of its IPs to the pool and stopping billing for them.

An IP that is still assigned to a consumer (a NIC, NAT gateway, load balancer or IP-failover group) cannot be released - detach the IP from that resource first, then delete the block. Use `ionosctl compute ipconsumer list --ipblock-id <id>` to see what is still holding an IP.

Use `--wait` (`-w`) to block until the deletion completes. Use `--force` to skip the confirmation prompt.

Required values to run command:

* IpBlock Id

## Options

```text
  -a, --all                 Release every IpBlock on the contract (only those whose IPs are not in use). Use instead of --ipblock-id
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [IpBlockId Name Location Size Ips State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -i, --ipblock-id string   The unique IpBlock Id (required)
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute ipblock delete --ipblock-id IPBLOCK_ID --wait
```

