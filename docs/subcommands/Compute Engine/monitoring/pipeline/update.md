---
description: "Rename a monitoring pipeline"
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

Update a pipeline's mutable properties. Today the only editable property is the name; the HTTP endpoint, Grafana endpoint, and ingest key are fixed for the life of the pipeline and are not affected here (rotate the key with 'monitoring key create').

Under the hood this reads the current pipeline and writes it back with the new values (GET + PUT), emulating a partial update, so unspecified properties are preserved.

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'monitoring' and env var 'IONOS_API_URL' (default "https://monitoring.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name GrafanaEndpoint HttpEndpoint Status]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/bhx, gb/lhr, fr/par, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string          The new human-friendly label for the pipeline. It does not affect the ingest or Grafana endpoints (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --pipeline-id string   The ID of the monitoring pipeline to update (from 'pipeline list') (required)
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a pipeline
ionosctl monitoring pipeline update --location de/txl --pipeline-id PIPELINE_ID --name new-name
```

