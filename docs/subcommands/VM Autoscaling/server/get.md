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
  -u, --api-url string     Override default host URL. Preferred over the config file override 'vmautoscaling' and env var 'IONOS_API_URL' (default "https://api.ionos.com/autoscaling")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [GroupServerId ServerId DatacenterId Name AvailabilityZone Cores RAM CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32        Controls the detail depth of the response objects (default 1)
  -f, --force              Force command to execute without user input
      --group-id string    ID of the autoscaling group that the server is a part of
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -i, --server-id string   ID of the autoscaling server
  -v, --verbose            Print step-by-step process when running command
```

## Examples

```text
ionosctl vm-autoscaling server get --group-id GROUP_ID --server-id SERVER_ID 
```

