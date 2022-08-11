---
description: Create a Data Platform Cluster
---

# DataplatformClusterCreate

## Usage

```text
ionosctl dataplatform cluster create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a new Data Platform Cluster. You must set the unique ID of the Datacenter, and the Name of the Cluster. If the other options are not set, the default values will be used. 

Required values to run command:

* Datacenter Id
* Name

## Options

```text
  -u, --api-url string            Override default host url (default "https://api.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name DataPlatformVersion MaintenanceWindow DatacenterId State] (default [ClusterId,Name,DataPlatformVersion,MaintenanceWindow,DatacenterId,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --datacenter-id string      The UUID of the virtual data center (VDC) the cluster is provisioned. (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
  -d, --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  -T, --maintenance-time string   Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format
  -n, --name string               The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. (required)
  -o, --output string             Desired output format [text|json] (default "text")
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout option for Cluster to be in AVAILABLE state[seconds] (default 1200)
  -v, --verbose                   Print step-by-step process when running command
  -V, --version string            The Data Platform version of your Cluster (default "1.1.0")
  -W, --wait-for-state            Wait for Cluster to be in AVAILABLE state
```

## Examples

```text
ionosctl dataplatform cluster create --datacenter-id DATACENTER_ID --name NAME --version DATA_PLATFORM_VERSION
```

