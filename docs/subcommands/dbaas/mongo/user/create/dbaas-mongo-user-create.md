---
description: Create MongoDB users.
---

# DbaasMongoUserCreate

## Usage

```text
ionosctl dbaas mongo user create [flags]
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
  -n, --name string         The authentication username (required)
  -p, --password string     The authentication password (required)
  -r, --roles               User's role for each db. DB1=Role1,DB2=Role2. Roles: read, readWrite, readAnyDatabase, readWriteAnyDatabase, dbAdmin, dbAdminAnyDatabase, clusterMonitor
```

## Examples

```text
ionosctl dbaas mongo user create --cluster-id CLUSTER_ID --name USERNAME --password PASSWORD --database DATABASE
```

