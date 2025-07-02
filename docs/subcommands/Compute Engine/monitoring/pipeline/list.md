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
  -u, --api-url string    Override default host URL (default "https://monitoring.de-fra.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name GrafanaEndpoint HttpEndpoint Status]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -l, --location string   The name of the location for monitoring pipeline 
      --no-headers        Don't print table headers when table output is used
      --order-by string   The field to order the results by. If not provided, the results will be ordered by the default field.
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl monitoring pipeline list --location de-fra
```

