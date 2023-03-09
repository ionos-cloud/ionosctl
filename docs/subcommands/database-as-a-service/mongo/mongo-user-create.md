---
description: Create MongoDB users.
---

# MongoUserCreate

## Usage

```text
mongo user create [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `create` command:

```text
[c]
```

## Description

Create MongoDB users.

## Options

```text
  -i, --cluster-id string   
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -h, --help                help for create
  -n, --name string         The authentication username (required)
      --no-headers          When using text output, don't print headers
  -p, --password string     The authentication password (required)
  -r, --roles               User's role for each db. DB1=Role1,DB2=Role2. Roles: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor
```

## Examples

```text
ionosctl dbaas mongo user create --cluster-id CLUSTER_ID --name USERNAME --password PASSWORD --database DATABASE
```

