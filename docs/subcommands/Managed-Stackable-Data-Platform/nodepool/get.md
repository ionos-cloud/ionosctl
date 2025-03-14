---
description: "Get Dataplatform Nodepool by cluster and nodepool id"
---

# DataplatformNodepoolGet

## Usage

```text
ionosctl dataplatform nodepool get [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
```

For `nodepool` command:

```text
[np]
```

For `get` command:

```text
[g]
```

## Description

Get Dataplatform Nodepool by cluster and nodepool id

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cluster-id string    The unique ID of the cluster. Must conform to the UUID format
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations ClusterId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -i, --nodepool-id string   The unique ID of the nodepool. Must conform to the UUID format
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for polling the request (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait                 Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl dataplatform nodepool get
```

