---
description: "List Kubernetes Versions"
---

# K8sVersionList

## Usage

```text
ionosctl k8s version list [flags]
```

## Aliases

For `version` command:

```text
[v]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve all available Kubernetes versions.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose count    Print step-by-step process when running command
```

## Examples

```text
ionosctl k8s version list
```

