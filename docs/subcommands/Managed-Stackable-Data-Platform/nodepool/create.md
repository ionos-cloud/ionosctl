---
description: "Create Dataplatform Nodepools"
---

# DataplatformNodepoolCreate

## Usage

```text
ionosctl dataplatform nodepool create [flags]
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

For `create` command:

```text
[c]
```

## Description

Node pools are the resources that powers the DataPlatformCluster.

The following requests allows to alter the existing resources, add or remove new resources to the cluster.

## Options

```text
  -A, --annotations stringToString   Annotations to set on a NodePool. It will overwrite the existing annotations, if there are any. Use the following format: --annotations KEY=VALUE,KEY=VALUE (default [])
  -u, --api-url string               Override default host URL. Preferred over the config file override 'dataplatform' and env var 'IONOS_API_URL' (default "https://api.ionos.com/dataplatform")
      --availability-zone string     The availability zone of the virtual datacenter region where the node pool resources should be provisioned
  -i, --cluster-id string            The UUID of the cluster the nodepool belongs to
      --cols strings                 Set of columns to be printed on output 
                                     Available columns: [Id Name Nodes Cores CpuFamily Ram Storage MaintenanceWindow State AvailabilityZone Labels Annotations ClusterId]
  -c, --config string                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int32                  The number of CPU cores per node
      --cpu-family string            A valid CPU family name or AUTO if the platform shall choose the best fitting option. Available CPU architectures can be retrieved from the datacenter resource
  -f, --force                        Force command to execute without user input
  -h, --help                         Print usage
  -L, --labels stringToString        Labels to set on a NodePool. It will overwrite the existing labels, if there are any. Use the following format: --labels KEY=VALUE,KEY=VALUE (default [])
      --maintenance-day string       Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur
      --maintenance-time string      Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59
  -n, --name string                  The name of your nodepool (required)
      --no-headers                   Don't print table headers when table output is used
      --node-count int32             The number of nodes that make up the node pool (required)
  -o, --output string                Desired output format [text|json|api-json] (default "text")
  -q, --quiet                        Quiet output
      --ram int32                    The RAM size for one node in MB. Must be set in multiples of 1024 MB, with a minimum size is of 2048 MB
      --storage-size int32           The size of the volume in GB. The size must be greater than 10GB
      --storage-type string          The type of hardware for the volume
  -t, --timeout int                  Timeout option for Request [seconds] (default 60)
  -v, --verbose                      Print step-by-step process when running command
  -w, --wait-for-request             Wait for the Request to be executed
```

