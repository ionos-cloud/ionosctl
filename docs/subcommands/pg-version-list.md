---
description: List DBaaS PostgreSQL Versions
---

# PgVersionList

## Usage

```text
ionosctl pg version list [flags]
```

## Aliases

For `pg` command:

```text
[postgres]
```

For `version` command:

```text
[v]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve all available DBaaS PostgreSQL versions.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [PostgresVersions] (default [PostgresVersions])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl pg version list
```

