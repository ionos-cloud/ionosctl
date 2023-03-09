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

For `delete` command:

```text
[g]
```

## Description

Delete a MongoDB user

## Options

```text
  -a, --all                 Delete all users in a cluster
  -i, --cluster-id string   The unique ID of the cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -d, --database string     The authentication database
  -f, --force               Skip y/n checks
  -h, --help                help for delete
      --name string         The authentication username
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas mongo user delete
```

