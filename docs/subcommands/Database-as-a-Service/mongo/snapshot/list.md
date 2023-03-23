---
description: List the snapshots of your Mongo Cluster
---

# DbaasMongoSnapshotList

## Usage

```text
ionosctl dbaas mongo snapshot list [flags]
```

## Aliases

For `mongo` command:

```text
[mongodb mdb m]
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
  -u, --api-url string      Override default host url (default "https://api.ionos.com")
  -i, --cluster-id string   The unique ID of the cluster (required)
      --cols strings        Set of columns to be printed on output 
                            Available columns: [SnapshotId CreationTime Size Version]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --no-headers          When using text output, don't print headers
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo cluster snapshot ls --cluster-id <cluster-id>
```

