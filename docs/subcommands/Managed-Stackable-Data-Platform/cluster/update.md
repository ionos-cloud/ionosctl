---
description: "Update a Dataplatform Cluster by ID"
---

# DataplatformClusterUpdate

## Usage

```text
ionosctl dataplatform cluster update [flags]
```

## Aliases

For `dataplatform` command:

```text
[mdp dp stackable managed-dataplatform]
```

For `cluster` command:

```text
[c]
```

For `update` command:

```text
[u]
```

## Description

Modifies the specified DataPlatformCluster by its distinct cluster ID. The fields in the request body are applied to the cluster. Note that the application to the cluster itself is performed asynchronously. You can check the sync state by querying the cluster with the GET method

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --cidr strings              The list of IPs and subnet for your cluster. Note the following unavailable IP ranges: 10.233.114.0/24 (required)
  -i, --cluster-id string         The unique ID of the cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id Name Version MaintenanceWindow DatacenterId State]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --lan-id string             The numeric LAN ID with which you connect your cluster (required)
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur (required)
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59 (required)
  -n, --name string               The name of the cluster
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
  -v, --verbose                   Print step-by-step process when running command
      --version string            The version of the cluster
```

## Examples

```text
ionosctl dataplatform cluster update --cluster-id <cluster-id>
```

