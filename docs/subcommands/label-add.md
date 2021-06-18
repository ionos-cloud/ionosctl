---
description: Add a Label to a Resource
---

# LabelAdd

## Usage

```text
ionosctl label add [flags]
```

## Aliases

For `add` command:

```text
[a]
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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [Key Value] (default [Key,Value])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -f, --force                  Force command to execute without user input
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
ionosctl label add --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID  --label-key LABEL_KEY --label-value LABEL_VALUE

ionosctl label add --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY --label-value LABEL_VALUE
```

