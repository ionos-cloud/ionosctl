---
description: Get Kubernetes Default Version
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
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for get
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          see step by step process when running a command
```

## Examples

```text
ionosctl k8s version get
```

