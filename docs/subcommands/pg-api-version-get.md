---
description: Get current version of DBaaS PostgreSQL API
---

# PgApiVersionGet

## Usage

```text
ionosctl pg api-version get [flags]
```

## Aliases

For `pg` command:

```text
[postgres]
```

For `api-version` command:

```text
[api info]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get the current version of DBaaS PostgreSQL API.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Name SwaggerUrl] (default [Name,SwaggerUrl])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl pg api-version get
```

