---
description: "Delete a Dataplatform Cluster by ID"
---

# DataplatformNodepoolDelete

## Usage

```text
ionosctl dataplatform nodepool delete [flags]
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

For `delete` command:

```text
[del d]
```

## Description

Delete a Dataplatform Cluster by ID

## Options

```text
  -a, --all                  Delete all clusters. If cluster ID is provided, delete all nodepools in given cluster
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cluster-id string    The unique ID of the cluster (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations ClusterId]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -i, --nodepool-id string   The unique ID of the nodepool (required)
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform cluster delete --cluster-id <cluster-id>
```

