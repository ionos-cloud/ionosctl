---
description: "Update a Virtual Data Center's name or description"
---

# DatacenterUpdate

## Usage

```text
ionosctl compute datacenter update [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `update` command:

```text
[u up]
```

## Description

Update the editable properties of an existing Virtual Data Center: its `--name` and `--description`.

Only these two fields can be changed. The VDC's region (`location`) is fixed at creation and is rejected by the API in update requests - to move workloads to another region you must create a new VDC there and recreate the resources. Renaming a VDC does not touch the resources inside it.

Pass only the flags you want to change; unspecified fields are left untouched.

Use `--wait` (`-w`) to block until the VDC is back in the AVAILABLE state.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location CpuFamily IPv6CidrBlock State Description Version Features SecAuthProtection]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -d, --description string     New free-text description for the VDC. Omit to leave unchanged
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            New human-friendly name for the VDC. Omit to leave unchanged
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Rename a VDC
ionosctl compute datacenter update --datacenter-id DATACENTER_ID --name "eu-prod"

# Change only the description and show the result
ionosctl compute datacenter update --datacenter-id DATACENTER_ID --description "Production workloads, EU" --cols "DatacenterId,Description"
```

