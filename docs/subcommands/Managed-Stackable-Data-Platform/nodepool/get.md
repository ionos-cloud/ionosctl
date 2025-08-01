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
  -u, --api-url string       Override default host URL. Preferred over the config file override 'dataplatform' and env var 'IONOS_API_URL' (default "https://api.ionos.com/dataplatform")
      --cluster-id string    The unique ID of the cluster. Must conform to the UUID format
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations ClusterId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -i, --nodepool-id string   The unique ID of the nodepool. Must conform to the UUID format
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform nodepool get
```

