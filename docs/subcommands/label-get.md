---
description: Get a Label
---

# LabelGet

## Usage

```text
ionosctl label get [flags]
```

## Description

Use this command to get information about a specified Label from a specified Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id
* Label Key

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings           Columns to be printed in the standard output (default [Key,Value])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
      --type string            Resource Type
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label get --resource-type datacenter --datacenter-id aa8e07a2-287a-4b45-b5e9-94761750a53c --label-key secondtest
Key          Value
secondtest   testdatacenter
```

