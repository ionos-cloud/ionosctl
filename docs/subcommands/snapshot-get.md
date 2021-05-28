---
description: Get a Snapshot
---

# SnapshotGet

## Usage

```text
ionosctl snapshot get [flags]
```

## Aliases

For `snapshot` command:
```text
[ss snap]
```

For `get` command:
```text
[g]
```

## Description

Use this command to get information about a specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [SnapshotId Name LicenceType Size State] (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                Force command to execute without user input
  -h, --help                 help for get
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -i, --snapshot-id string   The unique Snapshot Id (required)
```

## Examples

```text
ionosctl snapshot get --snapshot-id SNAPSHOT_ID
```

