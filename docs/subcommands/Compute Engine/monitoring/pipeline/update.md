---
description: "Partially modify a pipeline's properties. This command uses a combination of GET and PUT to simulate a PATCH operation"
---

# MonitoringPipelineUpdate

## Usage

```text
ionosctl monitoring pipeline update [flags]
```

## Aliases

For `pipeline` command:

```text
[p pipe]
```

For `update` command:

```text
[u]
```

## Description

Partially modify a pipeline's properties. This command uses a combination of GET and PUT to simulate a PATCH operation

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'monitoring' and env var 'IONOS_API_URL' (default "https://monitoring.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name GrafanaEndpoint HttpEndpoint Status]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra, de/txl, es/vit, gb/bhx, gb/lhr, fr/par, us/mci (default "de/fra")
  -n, --name string          The new name of the Monitoring Pipeline (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --pipeline-id string   The ID of the monitoring pipeline (required)
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl monitoring pipeline update --location de/txl --pipeline-id ID --name name
```

