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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Ip Mac NicId ServerId ServerName DatacenterId DatacenterName K8sNodePoolId K8sClusterId] (default [Ip,NicId,ServerId,DatacenterId,K8sNodePoolId,K8sClusterId])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32         Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --ipblock-id string   The unique IpBlock Id (required)
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Print step-by-step process when running command
```

## Examples

```text
ionosctl ipconsumer list --ipblock-id IPBLOCK_ID
```

