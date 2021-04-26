---
description: Remove a Label from a Volume
---

# RemoveLabel

## Usage

```text
ionosctl volume remove-label [flags]
```

## Description

Use this command to remove/delete a specified Label from a Volume.

Required values to run command:

* Data Center Id
* Volume Id
* Label Key

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [VolumeId,Name,Size,Type,LicenceType,State,Image])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id [Required flag]
      --force                  Force command to execute without user input
  -h, --help                   help for remove-label
      --label-key string       The unique Label Key [Required flag]
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --volume-id string       The unique Volume Id [Required flag]
```

## Examples

```text
ionosctl volume remove-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --volume-id 5d23eee2-45e5-44fe-96fe-e15aba2c48f5 --label-key test
```

