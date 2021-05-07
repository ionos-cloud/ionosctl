---
description: Create a Kubernetes Cluster
---

# K8sClusterCreate

## Usage

```text
ionosctl k8s cluster create [flags]
```

## Description

Use this command to create a new Managed Kubernetes Cluster. Regarding the name for the Kubernetes Cluster, the limit is 63 characters following the rule to begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. 

You can wait for the Cluster to be in "ACTIVE" state using `--wait-for-state` flag together with `--timeout` option.

Required values to run a command:

* K8s Cluster Name

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cluster-name string      The name for the K8s Cluster (required)
      --cluster-version string   The K8s version for the Cluster
      --cols strings             Columns to be printed in the standard output (default [ClusterId,Name,K8sVersion,State,MaintenanceWindow])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                    Force command to execute without user input
  -h, --help                     help for create
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for waiting for Cluster/Request [seconds] (default 600)
      --wait-for-request         Wait for the Request for Cluster creation to be executed
      --wait-for-state           Wait for the new Cluster to be in ACTIVE state
```

## Examples

```text
ionosctl k8s cluster create --cluster-name demoTest
ClusterId                              Name       K8sVersion  State
29d9b0c4-351d-4c9e-87e1-201cc0d49afb   demoTest   1.19.8      DEPLOYING
RequestId: 583ba6ae-dd0b-4c68-8fb2-41b3d7bc471b
Status: Command k8s cluster create has been successfully executed
```

