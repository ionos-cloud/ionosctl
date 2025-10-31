---
description: "List the snapshots of your Mongo Cluster"
---

# DbaasMongoSnapshotList

## Usage

```text
ionosctl dbaas mongo snapshot list [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `snapshot` command:

```text
[snap backup snapshots backups]
```

For `list` command:

```text
[ls]
```

## Description

List the snapshots of your Mongo Cluster

## Options

```text
  -u, --api-url string      Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [SnapshotId CreationTime Size Version]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int           pagination limit: Maximum number of items to return per request (default 50)
      --no-headers          Don't print table headers when table output is used
      --offset int          pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo cluster snapshot ls --cluster-id <cluster-id>
```

