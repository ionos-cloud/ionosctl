---
description: Get a MongoDB user
---

# MongoUserGet

## Usage

```text
mongo user get [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `get` command:

```text
[g]
```

## Description

Get a MongoDB user

## Options

```text
  -i, --cluster-id string   The unique ID of the cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -d, --database string     The authentication database
  -h, --help                help for get
      --no-headers          When using text output, don't print headers
  -u, --user string         The authentication username
```

## Examples

```text
ionosctl dbaas mongo user get
```

