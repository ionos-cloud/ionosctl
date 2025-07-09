---
description: "Create MongoDB users."
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

Create MongoDB users.

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string   
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -n, --name string         The authentication username (required)
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -p, --password string     The authentication password (required)
  -q, --quiet               Quiet output
  -r, --roles string        User's role for each db. DB1=Role1,DB2=Role2. Roles: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor, enableSharding (required)
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo user create --cluster-id CLUSTER_ID --name USERNAME --password PASSWORD --roles DATABASE=ROLE
```

