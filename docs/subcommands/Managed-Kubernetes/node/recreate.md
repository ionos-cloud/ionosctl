---
description: "Recreate a Kubernetes Node"
---

# K8sNodeRecreate

## Usage

```text
ionosctl k8s node recreate [flags]
```

## Aliases

For `node` command:

```text
[n]
```

For `recreate` command:

```text
[r]
```

## Description

You can recreate a single Kubernetes Node.

Managed Kubernetes starts a process which based on the NodePool's template creates & configures a new Node, waits for status "ACTIVE", and migrates all the Pods from the faulty Node, deleting it once empty. While this operation occurs, the NodePool will have an extra billable "ACTIVE" Node.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cluster-id string    The unique K8s Cluster Id (required)
      --cols strings         Set of columns to be printed on output 
                             Available columns: [NodeId Name K8sVersion PublicIP PrivateIP State] (default [NodeId,Name,K8sVersion,PublicIP,PrivateIP,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --no-headers           Don't print table headers when table output is used
  -i, --node-id string       The unique K8s Node Id (required)
      --nodepool-id string   The unique K8s Node Pool Id (required)
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl k8s node recreate --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID
```

