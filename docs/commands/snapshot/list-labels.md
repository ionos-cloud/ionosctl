---
description: List Labels from a Snapshot
---

# ListLabels

## Usage

```text
ionosctl snapshot list-labels [flags]
```

## Description

Use this command to list all Labels from a specified Snapshot.

Required values to run command:

* Snapshot Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                 help for list-labels
      --ignore-stdin         Force command to execute without user input
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique Snapshot Id [Required flag]
```

## Examples

```text
ionosctl snapshot list-labels --snapshot-id df7f4ad9-b942-4e79-939d-d1c10fb6fbff
Key    Value
test   testsnapshot
```

