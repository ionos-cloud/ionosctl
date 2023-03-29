---
description: List Mongo Clusters
---

# DbaasMongoClusterList

## Usage

```text
ionosctl dbaas mongo cluster list [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
```

For `cluster` command:

```text
[c]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to retrieve a list of Mongo Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.

## Options

```text
  -M, --max-results int32   The maximum number of elements to return
  -n, --name string         Response filter to list only the Mongo Clusters that contain the specified name in the DisplayName field. The value is case insensitive
      --offset int32        Skip a certain number of results
```

## Examples

```text
ionosctl dbaas mongo cluster list
```

