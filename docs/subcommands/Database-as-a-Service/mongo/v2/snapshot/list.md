---
description: "List MongoDB Snapshots"
---

# DbaasMongoV2SnapshotList

## Usage

```text
ionosctl dbaas mongo-v2 snapshot list [flags]
```

## Aliases

For `mongo-v2` command:

```text
[mongodb-v2 mg-v2]
```

For `snapshot` command:

```text
[snap snaps snapshots]
```

For `list` command:

```text
[ls l]
```

## Description

Use this command to retrieve a list of MongoDB Snapshots. You can filter by cluster ID using `--cluster-id`.

Each snapshot is a recovery WINDOW: the EarliestRecoveryTargetTime and LatestRecoveryTargetTime columns bound the range you can restore to (LatestRecoveryTargetTime is empty while the window extends to the present). Pick a specific point in the window with `--recovery-time` on `cluster restore` / `cluster create --snapshot-id`.

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'mongodbv2' and env var 'IONOS_API_URL' (default "https://mongodb.%s.ionos.com")
      --cluster-id string   Response filter to list only the snapshots of the specified cluster
      --cols strings        Set of columns to be printed on output 
                            Available columns: [SnapshotId ClusterId ClusterName Version SnapshotLocation EarliestRecoveryTargetTime LatestRecoveryTargetTime TotalSize RequiredStorageSize]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --limit int32         The maximum number of elements to return (default 100)
  -l, --location string     Location of the resource to operate on. When unset, list commands query all locations. Can be one of: gb/bhx, us/ewr. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint
      --no-headers          Don't print table headers when table output is used
      --offset int32        The first element to return
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dbaas mongo-v2 snapshot list
```

