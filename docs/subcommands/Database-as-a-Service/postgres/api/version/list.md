---
description: "List DBaaS PostgreSQL API Versions"
---

# DbaasPostgresApiVersionList

## Usage

```text
ionosctl dbaas postgres api-version list [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
```

For `api-version` command:

```text
[api info]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve all available DBaaS PostgreSQL API versions.

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Version SwaggerUrl] (default [Version,SwaggerUrl])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        Pagination limit: Maximum number of items to return per request (default 50)
      --no-headers       Don't print table headers when table output is used
      --offset int       Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas postgres api-version list
```

