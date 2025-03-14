---
description: "Get DBaaS PostgreSQLVersions for a Cluster"
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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the Cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [PostgresVersions] (default [PostgresVersions])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          Don't print table headers when table output is used
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for polling the request (default 60)
  -v, --verbose             Print step-by-step process when running command
  -w, --wait                Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl dbaas postgres version get --cluster-id CLUSTER_ID
```

