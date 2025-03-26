---
description: "List Servers that are managed by VM-Autoscaling Groups"
---

# VmAutoscalingServerList

## Usage

```text
ionosctl vm-autoscaling server list [flags]
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

For `list` command:

```text
[l ls]
```

## Description

List Servers that are managed by VM-Autoscaling Groups

## Options

```text
  -a, --all               If set, list all servers of all groups
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [GroupServerId ServerId DatacenterId Name AvailabilityZone Cores Ram CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32       Controls the detail depth of the response objects (default 1)
  -f, --force             Force command to execute without user input
  -i, --group-id string   ID of the autoscaling group to list servers from
  -h, --help              Print usage
      --no-headers        Don't print table headers when table output is used
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for polling the request (default 60)
  -v, --verbose           Print step-by-step process when running command
  -w, --wait              Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl vm-autoscaling server list --group-id GROUP_ID
ionosctl vm-autoscaling server list --all ALL
```

