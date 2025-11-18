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
<<<<<<< HEAD
  -u, --api-url string   Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers       Don't print table headers when table output is used
      --offset int       Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
      --query string     JMESPath query string to filter the output
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
=======
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
>>>>>>> 8e970fd7 (remove deprecated 'D' for 'datacenter-id' only on psql)
```

## Examples

```text
ionosctl k8s version get
```

