---
description: "Retrieve pipelines"
---

# MonitoringPipelineList

## Usage

```text
ionosctl monitoring pipeline list [flags]
```

## Aliases

For `pipeline` command:

```text
[p pipe]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve pipelines

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'monitoring' and env var 'IONOS_API_URL' (default "https://monitoring.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name GrafanaEndpoint HttpEndpoint Status]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/bhx, gb/lhr, fr/par, us/mci (default "de/fra")
      --no-headers        Don't print table headers when table output is used
      --offset int        Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl monitoring pipeline list --location de/txl
```

