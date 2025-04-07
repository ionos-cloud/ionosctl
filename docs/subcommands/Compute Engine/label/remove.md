---
description: "Remove a Label from a Resource"
---

# LabelRemove

## Usage

```text
ionosctl label remove [flags]
```

## Aliases

For `remove` command:

```text
[delete del r rm]
```

## Description

Use this command to remove a Label from a Resource.

Required values to run command:

* Resource Type
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id, Image ID, or Snapshot Id
* Label Key

## Options

```text
  -a, --all                    Remove all Labels
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [URN Key Value ResourceType ResourceId] (default [URN,Key,Value,ResourceType,ResourceId])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --image-id string        The unique Image Id(note: only private images supported)
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Type of resource to remove labels from. Can be one of: datacenter, volume, server, snapshot, ipblock, image (required)
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
  -v, --verbose                Print step-by-step process when running command
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label remove --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY
```

