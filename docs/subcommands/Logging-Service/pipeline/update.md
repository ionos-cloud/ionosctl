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

Update a logging pipeline

## Options

```text
  -u, --api-url string            Override default host url (default "logging.de-txl.ionos.com")
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [Id Name GrafanaAddress CreatedDate State] (default [Id,Name,GrafanaAddress,CreatedDate,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --json-properties string    Path to a JSON file containing the desired properties. Overrides any other properties set.
      --json-properties-example   If set, prints a complete JSON which could be used for --json-properties and exits. Hint: Pipe me to a .json file
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string        The ID of the logging pipeline (required)
  -q, --quiet                     Quiet output
  -v, --verbose                   Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service pipeline update --pipeline-id ID --json-properties PATH_TO_FILE
```

