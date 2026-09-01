---
description: "List the resources consuming each IP in an IP block"
---

# IpconsumerList

## Usage

```text
ionosctl compute ipconsumer list [flags]
```

## Aliases

For `ipconsumer` command:

```text
[ipc]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to list, for every IP address in a reserved IP block, the resource that is currently using it: the NIC (and its MAC), the owning server and datacenter, and any Kubernetes cluster / node pool.

An empty result means none of the block's addresses are in use, so the block can be safely released.

Required values to run command:

* IpBlock Id (get it from `ionosctl compute ipblock list`)

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Ip Mac NicId ServerId ServerName DatacenterId DatacenterName K8sNodePoolId K8sClusterId]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --ipblock-id string   The ID of the reserved IP block whose addresses you want to inspect (required)
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
# List consumers of every IP in a block
ionosctl compute ipconsumer list --ipblock-id IPBLOCK_ID

# Show only the IP, server and datacenter columns
ionosctl compute ipconsumer list --ipblock-id IPBLOCK_ID --cols Ip,ServerName,DatacenterName
```

