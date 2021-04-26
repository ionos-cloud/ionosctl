---
description: Get a Label from a Snapshot
---

# GetLabel

## Usage

```text
ionosctl snapshot get-label [flags]
```

## Description

Use this command to get information about a specified Label from a Snapshot.

Required values to run command:

* Snapshot Id
* Label Key

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                Force command to execute without user input
  -h, --help                 help for get-label
      --label-key string     The unique Label Key [Required flag]
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique Snapshot Id [Required flag]
```

## Examples

```text
 ionosctl snapshot get-label --snapshot-id df7f4ad9-b942-4e79-939d-d1c10fb6fbff --label-key test
Key    Value
test   testsnapshot
```

