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
[pg pgsql postgresql psql]
```

For `database` command:

```text
[databases]
```

## Description

Create a new logical database in the specified cluster.

The --owner must be an existing user (role) in the same cluster and is granted ownership (full privileges) of the new database. The cluster must be AVAILABLE.

Required values to run command:

* Cluster Id
* Database (name)
* Owner

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
      --cluster-id string   ID of the PostgreSQL cluster to create the database in (must be AVAILABLE)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Owner ClusterId]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --database string     Name of the database to create (1-63 characters)
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --owner string        Name of an existing user (role) in the same cluster that will own the database and hold full privileges on it
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas postgres database create --cluster-id CLUSTER_ID --database orders --owner appuser
```

