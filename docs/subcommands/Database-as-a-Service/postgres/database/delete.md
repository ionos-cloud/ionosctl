---
description: "Delete database"
---

# DbaasPostgresDatabaseDelete

## Usage

```text
ionosctl dbaas postgres database delete [flags]
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

Delete the specified database from the given cluster

## Options

```text
  -i, --cluster-id string   The ID of the Postgres cluster
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --database string     The name of the database
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas postgres database delete --cluster-id <cluster-id> --database <database>
```

