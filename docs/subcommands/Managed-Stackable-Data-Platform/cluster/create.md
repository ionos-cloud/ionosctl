---
description: "Create Dataplatform Cluster"
---

# DataplatformClusterCreate

## Usage

```text
ionosctl dataplatform cluster create [flags]
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

For `create` command:

```text
[c]
```

## Description

The cluster will be provisioned in the datacenter matching the provided datacenterID. Therefore the datacenter must be created upfront and must be accessible by the user issuing the request

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id Name Version MaintenanceWindow DatacenterId State]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -i, --datacenter-id string      The ID of the connected datacenter (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur (required)
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59 (required)
  -n, --name string               The name of your cluster (required)
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
  -t, --timeout duration          Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose                   Print step-by-step process when running command
      --version string            The version of your dataplatform cluster (default "same as 'dataplatform version active'")
  -w, --wait                      Polls the request continuously until the operation is completed
```

