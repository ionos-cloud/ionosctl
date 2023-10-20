---
description: "Get a Label"
---

# LabelGet

## Usage

```text
ionosctl label get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Label from a specified Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id or Snapshot Id
* Label Key

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [URN Key Value ResourceType ResourceId] (default [URN,Key,Value,ResourceType,ResourceId])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Type of resource to get labels from. Can be one of: datacenter, volume, server, snapshot, ipblock (required)
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
  -v, --verbose                Print step-by-step process when running command
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label get --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY
```

