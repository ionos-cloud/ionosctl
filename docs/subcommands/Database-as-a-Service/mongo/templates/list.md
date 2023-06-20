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
[mongodb mdb m]
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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [TemplateId Name Edition Cores StorageSize Ram] (default [TemplateId,Name,Edition,Cores,StorageSize,Ram])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --offset int32        Skip a certain number of results
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo templates list
```

