---
description: Add a Label to a Resource
---

# LabelAdd

## Usage

```text
ionosctl label add [flags]
```

## Description

Use this command to add a Label to a specific Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id
* Label Key
* Label Value

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -f, --force                  Force command to execute without user input
  -F, --format strings         Collection of fields to be printed on output (default [Key,Value])
  -h, --help                   help for add
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
      --label-value string     The unique Label Value (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Resource Type
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label add --resource-type server --datacenter-id aa8e07a2-287a-4b45-b5e9-94761750a53c --server-id 1dc7c6a8-5ab3-4fa8-83e7-9d989bd52ffa  --label-key test --label-value testserver
Key    Value
test   testserver

ionosctl label add --resource-type datacenter --datacenter-id aa8e07a2-287a-4b45-b5e9-94761750a53c --label-key secondtest --label-value testdatacenter
Key          Value
secondtest   testdatacenter
```

