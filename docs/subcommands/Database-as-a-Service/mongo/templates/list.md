---
description: "List Mongo Templates"
---

# DbaasMongoTemplatesList

## Usage

```text
ionosctl dbaas mongo templates list [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `templates` command:

```text
[t]
```

For `list` command:

```text
[l ls]
```

## Description

Retrieves a list of valid templates. These templates can be used to create MongoDB clusters; they contain properties, such as number of cores, RAM, and the storage size.

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [TemplateId Name Edition Cores StorageSize RAM] (default [TemplateId,Name,Edition,Cores,StorageSize,RAM])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          Don't print table headers when table output is used
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo templates list
```

