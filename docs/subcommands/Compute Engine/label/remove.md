---
description: "Remove a Label from a Resource"
---

# LabelRemove

## Usage

```text
ionosctl compute label remove [flags]
```

## Aliases

For `remove` command:

```text
[delete del r rm]
```

## Description

Use this command to remove a Label from a Resource. Select the resource with --resource-type plus its id flag(s) (see the same pairing as `label add`) and identify the label by --label-key.

Use --all to remove every label. With --resource-type and its id flag(s) it removes all labels on that one resource; with no --resource-type it iterates over labels of ALL resources under your account (prompting for each unless --force is given).

Required values to run command:

* Resource Type
* Resource Id(s) for that type
* Label Key

## Options

```text
  -a, --all                    Remove all labels: on the selected resource when --resource-type is given, otherwise across every labeled resource on the account
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [URN Key Value ResourceType ResourceId]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The Data Center Id. Required for --resource-type datacenter, server and volume
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --image-id string        The Image Id (private images only). Required for --resource-type image
      --ipblock-id string      The IpBlock Id. Required for --resource-type ipblock
      --label-key string       The key of the label to remove from the resource (required)
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --resource-type string   The kind of resource to remove the label from. Determines which id flag(s) are required. Can be one of: datacenter, volume, server, snapshot, ipblock, image (required)
      --server-id string       The Server Id (also needs --datacenter-id). Required for --resource-type server
      --snapshot-id string     The Snapshot Id. Required for --resource-type snapshot
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --volume-id string       The Volume Id (also needs --datacenter-id). Required for --resource-type volume
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Remove one label from a datacenter
ionosctl compute label remove --resource-type datacenter --datacenter-id DATACENTER_ID --label-key env

# Remove all labels from a server
ionosctl compute label remove --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID --all
```

