---
description: "Get Kubernetes Default Version"
---

# K8sVersionGet

## Usage

```text
ionosctl k8s version get [flags]
```

## Aliases

For `version` command:

```text
[v]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve the current default Kubernetes version for Clusters and NodePools.

## Options

```text
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout duration   Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait               Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl k8s version get
```

