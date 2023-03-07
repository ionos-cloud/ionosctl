---
description: Retrieves a list of MongoDB users.
---

# MongoUserList

## Usage

```text
mongo user list [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `list` command:

```text
[l ls]
```

## Description

Retrieves a list of MongoDB users.

## Options

```text
  -a, --all                 List all users, across all clusters
  -i, --cluster-id string   
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Username CreatedBy Roles]
  -h, --help                help for list
  -M, --max-results int32   The maximum number of elements to return
      --no-headers          When using text output, don't print headers
      --offset int32        Skip a certain number of results
```

## Examples

```text
ionosctl dbaas mongo user list
```

