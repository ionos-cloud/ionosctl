---
description: Add a Label on a Snapshot
---

# AddLabel

## Usage

```text
ionosctl snapshot add-label [flags]
```

## Description

Use this command to create a Label on Snapshot. You must specify the key and the value for the Label.

Required values to run command: 

* Snapshot Id 
* Label Key
* Label Value

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings         Columns to be printed in the standard output (default [SnapshotId,Name,LicenceType,Size,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                 help for add-label
      --ignore-stdin         Force command to execute without user input
      --label-key string     The unique Label Key [Required flag]
      --label-value string   The unique Label Value [Required flag]
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
      --snapshot-id string   The unique Snapshot Id [Required flag]
```

## Examples

```text
ionosctl snapshot add-label --snapshot-id df7f4ad9-b942-4e79-939d-d1c10fb6fbff --label-key test --label-value testsnapshot
Key    Value
test   testsnapshot
```

