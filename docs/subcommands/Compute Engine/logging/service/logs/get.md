---
description: "Retrieve a log from a logging pipeline"
---

# LoggingServiceLogsGet

## Usage

```text
ionosctl logging-service logs get [flags]
```

## Description

Retrieve a log from a logging pipeline

## Options

```text
  -u, --api-url string       Override default host url (default "logging.de-txl.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Tag Source Protocol Public Destinations] (default [Tag,Source,Protocol,Public,Destinations])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --log-tag string       The tag of the pipeline log that you want to retrieve (required)
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string   The ID of the logging pipeline (required)
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service logs get --pipeline-id ID --log-tag TAG
```

