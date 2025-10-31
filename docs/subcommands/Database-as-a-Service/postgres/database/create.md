---
description: "Create database"
---

# DbaasPostgresDatabaseCreate

## Usage

```text
ionosctl dbaas postgres database create [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
```

For `database` command:

```text
[databases]
```

## Description

Create a new database in the specified cluster

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cluster-id string   The ID of the Postgres cluster
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --database string     The name of the database
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --owner string        The owner of the database
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas postgres database create --cluster-id <cluster-id> --database <database> --owner <owner>
```

