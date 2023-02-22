---
description: Delete a MongoDB user
---

# MongoUserDelete

## Usage

```text
mongo user delete [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `user` command:

```text
[u]
```

For `delete` command:

```text
[g]
```

## Description

Delete a MongoDB user

## Options

```text
  -i, --cluster-id string   The unique ID of the cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -d, --database string     The authentication database
  -h, --help                help for delete
      --no-headers          When using text output, don't print headers
  -u, --user string         The authentication username
```

## Examples

```text
ionosctl dbaas mongo user delete
```

