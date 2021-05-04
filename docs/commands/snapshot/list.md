---
description: List Snapshots
---

# List

## Usage

```text
ionosctl snapshot list [flags]
```

## Description

Use this command to get a list of Snapshots.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl snapshot list 
SnapshotId                             Name           LicenceType   Size
dc688daf-8e54-4db8-ac4a-487ad5a34e9c   testSnapshot   LINUX         10
8e0bc509-87ee-47f4-a382-302e4f7e103d   image          LINUX         10
```

