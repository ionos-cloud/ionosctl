---
description: "List databases"
---

# DbaasPostgresDatabaseList

## Usage

```text
ionosctl dbaas postgres database list [flags]
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

For `list` command:

```text
[ls]
```

## Description

List databases in the given cluster

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
  -i, --cluster-id string   The ID of the Postgres cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Owner ClusterId]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas postgres database list
```

