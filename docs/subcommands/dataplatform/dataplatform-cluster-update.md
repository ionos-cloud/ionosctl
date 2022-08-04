---
description: Update a Data Platform Cluster
---

# DataplatformClusterUpdate

## Usage

```text
ionosctl dataplatform cluster update [flags]
```

## Aliases

For `dataplatform` command:

```text
[dp]
```

For `cluster` command:

```text
[c]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update attributes of a Data Platform Cluster.

Required values to run command:

* Cluster Id

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string         The unique ID of the Cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name DataPlatformVersion MaintenanceWindow DatacenterId State] (default [ClusterId,Name,DataPlatformVersion,MaintenanceWindow,DatacenterId,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -d, --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  -T, --maintenance-time string   Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  -n, --name string               The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose                   Print step-by-step process when running command
  -V, --version string            The Data Platform version of your cluster
  -W, --wait-for-state            Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dataplatform cluster update -i CLUSTER_ID -n CLUSTER_NAME
```

