---
description: "Add a Label to a Resource"
---

# LabelAdd

## Usage

```text
ionosctl compute label add [flags]
```

## Aliases

For `add` command:

```text
[a]
```

## Description

Use this command to add (or overwrite) a Label on a specific Resource.

Pick the target with --resource-type and supply the matching id flag(s):
  * datacenter -> --datacenter-id
  * server     -> --datacenter-id and --server-id
  * volume     -> --datacenter-id and --volume-id
  * snapshot   -> --snapshot-id
  * ipblock    -> --ipblock-id
  * image      -> --image-id (private images only)

Adding a key that already exists on the resource overwrites its value (there is one value per key per resource).

Required values to run command:

* Resource Type
* Resource Id(s) as listed above
* Label Key
* Label Value

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [URN Key Value ResourceType ResourceId]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The Data Center Id. Required for --resource-type datacenter, server and volume
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --image-id string        The Image Id (private images only; public images cannot be labeled). Required for --resource-type image
      --ipblock-id string      The IpBlock Id. Required for --resource-type ipblock
      --label-key string       The label key. Unique per resource; adding an existing key overwrites its value (required)
      --label-value string     The label value to store under the key (required)
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --resource-type string   The kind of resource to label. Determines which id flag(s) are required (see command description). Can be one of: datacenter, volume, server, snapshot, ipblock, image (required)
      --server-id string       The Server Id (also needs --datacenter-id). Required for --resource-type server
      --snapshot-id string     The Snapshot Id. Required for --resource-type snapshot
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
      --volume-id string       The Volume Id (also needs --datacenter-id). Required for --resource-type volume
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Label a server
ionosctl compute label add --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID --label-key env --label-value prod

# Label a datacenter
ionosctl compute label add --resource-type datacenter --datacenter-id DATACENTER_ID --label-key team --label-value payments
```

