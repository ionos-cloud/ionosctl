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
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [Version SwaggerUrl] (default [Version,SwaggerUrl])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       When using text output, don't print headers
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres api-version list
```

