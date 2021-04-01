---
description: Get a Snapshot
---

# Get

## Usage

```text
ionosctl snapshot get [flags]
```

## Description

Use this command to get information about a specified Snapshot.

Required values to run command:
- Snapshot Id

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                 help for get
      --ignore-stdin         Force command to execute without user input
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique Snapshot Id [Required flag]
```

