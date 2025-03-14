---
description: "List DBaaS PostgreSQL Versions"
---

# DbaasPostgresVersionList

## Usage

```text
ionosctl dbaas postgres version list [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
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
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -t, --timeout int      Timeout in seconds for polling the request (default 60)
  -v, --verbose          Print step-by-step process when running command
  -w, --wait             Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl dbaas postgres version list
```

