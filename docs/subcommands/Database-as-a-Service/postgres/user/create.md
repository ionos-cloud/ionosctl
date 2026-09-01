---
description: "Create a user"
---

# DbaasPostgresUserCreate

## Usage

```text
ionosctl dbaas postgres user create [flags]
```

## Aliases

For `postgres` command:

```text
[pg pgsql postgresql psql]
```

For `user` command:

```text
[usr u users]
```

## Description

Create a new PostgreSQL login role in the given cluster.

The cluster must already be AVAILABLE. The new user can subsequently be set as the owner of a database via 'dbaas postgres database create --owner'.

Required values to run command:

* Cluster Id
* User (name)
* Password

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'psql' and env var 'IONOS_API_URL' (default "https://api.ionos.com/databases/postgresql")
  -i, --cluster-id string   ID of the PostgreSQL cluster the user is created in (must be AVAILABLE)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Username System ClusterId]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     Login password for the new user. Minimum 10 characters
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
      --user string         Name of the PostgreSQL login role to create. Must not collide with a reserved system name (e.g. postgres)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas postgres user create --cluster-id CLUSTER_ID --user appuser --password 'S3cr3tPassw0rd'
```

