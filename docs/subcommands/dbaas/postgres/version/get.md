---
description: Get DBaaS PostgreSQLVersions for a Cluster
---

# DbaasPostgresVersionGet

## Usage

```text
ionosctl dbaas postgres version get [flags]
```

## Aliases

For `postgres` command:

```text
[pg]
```

For `version` command:

```text
[v]
```

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve a list of all PostgreSQL versions available for a specified Cluster.

Required values to run command:

* Cluster Id

## Options

```text
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --no-headers          When using text output, don't print headers
```

## Examples

```text
ionosctl dbaas postgres version get --cluster-id CLUSTER_ID
```

