---
description: Remove a Label from a Resource
---

# LabelRemove

## Usage

```text
ionosctl label remove [flags]
```

## Aliases

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a Label from a Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id
* Label Key

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [Key Value] (default [Key,Value])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -f, --force                  Force command to execute without user input
  -h, --help                   help for remove
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Resource Type
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label remove --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY
```

