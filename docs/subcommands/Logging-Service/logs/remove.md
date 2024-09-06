---
description: "Remove a log from a logging pipeline. NOTE:There needs to be at least one log in a pipeline at all times."
---

# LoggingServiceLogsRemove

## Usage

```text
ionosctl logging-service logs remove [flags]
```

## Description

Remove a log from a logging pipeline. NOTE:There needs to be at least one log in a pipeline at all times.

## Options

```text
  -u, --api-url string       Override default host url (default "logging.de-txl.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Tag Source Protocol Public Destinations] (default [Tag,Source,Protocol,Public,Destinations])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --log-tag string       The tag of the pipeline log that you want to remove (required)
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --pipeline-id string   The ID of the logging pipeline (required)
  -q, --quiet                Quiet output
  -v, --verbose              Print step-by-step process when running command
```

## Examples

```text
ionosctl logging-service logs remove --pipeline-id ID --log-tag TAG
```

