---
description: "List IpConsumers"
---

# IpconsumerList

## Usage

```text
ionosctl ipconsumer list [flags]
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

Use this command to get a list of Resources where each IP address from an IpBlock is being used.

Required values to run command:

* IpBlock Id

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Ip Mac NicId ServerId ServerName DatacenterId DatacenterName K8sNodePoolId K8sClusterId] (default [Ip,NicId,ServerId,DatacenterId,K8sNodePoolId,K8sClusterId])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --ipblock-id string   The unique IpBlock Id (required)
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl ipconsumer list --ipblock-id IPBLOCK_ID
```

