---
description: List Kubernetes Clusters
---

# K8sClusterList

## Usage

```text
ionosctl k8s cluster list [flags]
```

## Aliases

For `cluster` command:
```text
[c]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to get a list of existing Kubernetes Clusters.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [ClusterId Name K8sVersion State MaintenanceWindow AvailableUpgradeVersions ViableNodePoolVersions Public GatewayIp] (default [ClusterId,Name,K8sVersion,Public,State,MaintenanceWindow])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl k8s cluster list
```

