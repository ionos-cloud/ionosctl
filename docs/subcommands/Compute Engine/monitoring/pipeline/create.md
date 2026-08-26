---
description: "Create a monitoring pipeline"
---

# MonitoringPipelineCreate

## Usage

```text
ionosctl monitoring pipeline create [flags]
```

## Aliases

For `pipeline` command:

```text
[p pipe]
```

For `create` command:

```text
[post c]
```

## Description

Create a new monitoring pipeline in the region given by --location. The only configurable property at creation is the name; the service provisions the rest and returns:
  * HttpEndpoint    - Prometheus remote-write target. Push metrics to <HttpEndpoint>/api/v1/push with the ingest key in the APIKEY header.
  * GrafanaEndpoint - the managed Grafana base URL for querying the ingested metrics.
  * Status          - the provisioning state; the pipeline is usable once it reports as available.

The ingest key is NOT part of this command's output. It is shown only once, immediately after creation, so retrieve or rotate it with 'monitoring key create --pipeline-id <id>' before configuring agents. An account may hold up to 10 pipelines by default (adjustable via Support).

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'monitoring' and env var 'IONOS_API_URL' (default "https://monitoring.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name GrafanaEndpoint HttpEndpoint Status]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra, de/txl, es/vit, gb/bhx, gb/lhr, fr/par, us/mci. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
  -n, --name string       A human-friendly label for the pipeline, shown in listings and the DCD. It does not affect the ingest or Grafana endpoints
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a pipeline in Berlin
ionosctl monitoring pipeline create --location de/txl --name my-metrics

# Create a pipeline in Frankfurt and immediately generate its ingest key, capturing both
PIPE=$(ionosctl monitoring pipeline create --location de/fra --name prod-metrics --cols Id --no-headers)
ionosctl monitoring key create --location de/fra --pipeline-id "$PIPE"
```

