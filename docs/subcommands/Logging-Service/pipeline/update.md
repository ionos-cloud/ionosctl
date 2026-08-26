---
description: "Update a logging pipeline"
---

# LoggingServicePipelineUpdate

## Usage

```text
ionosctl logging-service pipeline update [flags]
```

## Aliases

For `logging-service` command:

```text
[log-svc]
```

For `pipeline` command:

```text
[p pipelines]
```

## Description

Replace a pipeline's properties (its name and the full set of logs) from a JSON file. The JSON is sent as a patch, so the 'logs' array you provide REPLACES the existing logs - include every log you want to keep, not just the changed ones. Use --json-properties-example to print a template with the expected shape.

To tweak a single log stream in place instead, prefer 'ionosctl logging-service logs update', which edits one log by its tag without rewriting the whole pipeline.

The pipeline's ingestion/Grafana addresses and its key are unaffected by an update.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'logging' and env var 'IONOS_API_URL' (default "https://logging.%s.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id Name GrafanaAddress TCPAddress HTTPAddress CreatedDate State]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
  -F, --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example   If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/txl, de/fra, gb/lhr, fr/par, es/vit, us/mci, gb/bhx. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string        The ID of the logging pipeline to update (required)
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -t, --timeout int               Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                      Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Replace the pipeline's name and logs from a file
ionosctl logging-service pipeline update --location de/txl --pipeline-id ID --json-properties ./pipeline.json

# Print a JSON template to edit
ionosctl logging-service pipeline update --json-properties-example
```

