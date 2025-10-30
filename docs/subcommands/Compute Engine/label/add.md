---
description: "Add a Label to a Resource"
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
* Resource Id: Datacenter Id, Server Id, Volume Id, IpBlock Id, Image ID, or Snapshot Id
* Label Key
* Label Value

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [URN Key Value ResourceType ResourceId] (default [URN,Key,Value,ResourceType,ResourceId])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --image-id string        The unique Image Id(note: only private images supported)
      --ipblock-id string      The unique IpBlock Id
      --label-key string       The unique Label Key (required)
      --label-value string     The unique Label Value (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Type of resource to add labels to. Can be one of: datacenter, volume, server, snapshot, ipblock, image (required)
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label add --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID  --label-key LABEL_KEY --label-value LABEL_VALUE

ionosctl label add --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY --label-value LABEL_VALUE
```

