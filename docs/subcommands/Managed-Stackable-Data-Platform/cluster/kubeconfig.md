---
description: "Get a Dataplatform Cluster's Kubeconfig by ID"
---

# DataplatformClusterKubeconfig

## Usage

```text
ionosctl dataplatform cluster kubeconfig [flags]
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

For `kubeconfig` command:

```text
[k]
```

## Description

Get a Dataplatform Cluster's Kubeconfig by ID

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'dataplatform' and env var 'IONOS_API_URL' (default "https://api.ionos.com/dataplatform")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Version MaintenanceWindow DatacenterId State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dataplatform cluster kubeconfig --cluster-id <cluster-id>
```

