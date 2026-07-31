---
description: "Create a MongoDB user in a cluster"
---

# DbaasMongoUserCreate

## Usage

```text
ionosctl dbaas mongo user create [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `user` command:

```text
[u]
```

For `create` command:

```text
[c]
```

## Description

Create a database user in a MongoDB cluster and grant it one or more roles.

Roles are given as --roles DATABASE=ROLE pairs, comma-separated: each pair grants ROLE on the named DATABASE (MongoDB authorization is per-database). The DATABASE need not exist yet. Valid role names: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor, enableSharding. For the "...AnyDatabase"/clusterMonitor roles the database is conventionally 'admin'.

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster the user is created in (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           Maximum number of items to return per request (default 50)
  -n, --name string         Username the new user authenticates with (required)
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     Password the new user authenticates with (required)
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -r, --roles string        Comma-separated DATABASE=ROLE grants, e.g. mydb=readWrite,admin=clusterMonitor. Valid roles: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor, enableSharding (required)
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# A user with read/write on one application database
ionosctl dbaas mongo user create --cluster-id <cluster-id> --name appuser --password <password> --roles mydb=readWrite

# A monitoring user plus admin rights on two databases
ionosctl dbaas mongo user create --cluster-id <cluster-id> --name ops --password <password> --roles admin=clusterMonitor,orders=dbAdmin,billing=dbAdmin
```

