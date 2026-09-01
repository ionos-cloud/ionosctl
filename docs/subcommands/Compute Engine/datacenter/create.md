---
description: "Create a Virtual Data Center in a chosen region"
---

# DatacenterCreate

## Usage

```text
ionosctl compute datacenter create [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `create` command:

```text
[c]
```

## Description

Create a Virtual Data Center (VDC) - the top-level, network-isolated container that will hold your compute resources (servers, volumes, LANs, NICs, firewalls).

Everything you provision afterwards is created inside a VDC and inherits its region, so the single most important choice here is `--location`: the region the VDC lives in. This is set once at creation and CANNOT be changed later - to run workloads in another region you create a separate VDC there. You can provision as many VDCs as your contract allows; each is logically segmented from the others.

The name defaults to "Unnamed Data Center" and the location defaults to `de/txl` (Berlin), so a VDC can be created with no flags at all - but it is strongly recommended to pass an explicit `--location` so you are not surprised by where your resources land.

Use `--wait` (`-w`) to block until the VDC reaches the AVAILABLE state before the command returns; without it the command returns as soon as the provisioning request is accepted.

## Options

```text
  -u, --api-url string       Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [DatacenterId Name Location CpuFamily IPv6CidrBlock State Description Version Features SecAuthProtection]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
  -d, --description string   Free-text description of the VDC's purpose. Optional and editable later
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Region the VDC and all resources inside it will live in, e.g. de/txl (Berlin), de/fra (Frankfurt), gb/lhr (London), es/vit (Logrono), fr/par (Paris), us/las (Las Vegas), us/ewr (Newark). IMMUTABLE - cannot be changed after creation. Must be enabled for your contract (default "de/txl")
  -n, --name string          Human-friendly name shown in the DCD and CLI listings. Does not need to be unique (default "Unnamed Data Center")
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a VDC in Frankfurt with an explicit name
ionosctl compute datacenter create --name "prod-vdc" --location de/fra

# Create a VDC and wait until it is AVAILABLE, with a description, showing only chosen columns
ionosctl compute datacenter create --name "prod-vdc" --description "Production workloads, EU" --location de/fra --wait --cols "DatacenterId,Name,Location,State"
```

