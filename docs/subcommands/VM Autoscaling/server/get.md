---
description: "Get a VM Autoscaling Server"
---

# VmAutoscalingServerGet

## Usage

```text
ionosctl vm-autoscaling server get [flags]
```

## Aliases

For `vm-autoscaling` command:

```text
[vmas vm-as vmasc vm-asc vmautoscaling]
```

For `server` command:

```text
[s sv vm vms servers]
```

For `get` command:

```text
[g]
```

## Description

Get a VM Autoscaling Server

## Options

```text
  -u, --api-url string     Override default host URL. Preferred over the config file override 'autoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [GroupServerId ServerId DatacenterId Name AvailabilityZone Cores RAM CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId NicMultiQueue]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int          Level of detail for response objects (default 1)
  -F, --filters strings    Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force              Force command to execute without user input
      --group-id string    ID of the autoscaling group that the server is a part of
  -h, --help               Print usage
      --limit int          Maximum number of items to return per request (default 50)
      --no-headers         Don't print table headers when table output is used
      --offset int         Number of items to skip before starting to collect the results
      --order-by string    Property to order the results by
  -o, --output string      Desired output format [text|json|api-json] (default "text")
      --query string       JMESPath query string to filter the output
  -q, --quiet              Quiet output
  -i, --server-id string   ID of the autoscaling server
  -v, --verbose count      Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl vm-autoscaling server get --group-id GROUP_ID --server-id SERVER_ID 
```

