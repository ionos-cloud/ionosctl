---
description: Delete a MongoDB user
---

# DbaasMongoUserDelete

## Usage

```text
ionosctl dbaas mongo user delete [flags]
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
  -d, --database string     The authentication database
  -f, --force               Skip y/n checks
      --name string         The authentication username
```

## Examples

```text
ionosctl dbaas mongo user delete
```

