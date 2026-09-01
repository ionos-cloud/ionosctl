---
description: "Delete a Virtual Data Center and everything inside it"
---

# DatacenterDelete

## Usage

```text
ionosctl compute datacenter delete [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `delete` command:

```text
[d]
```

## Description

Delete a Virtual Data Center. This removes the VDC object AND every resource it contains - servers, volumes (and their data), LANs, NICs, firewall rules - in a single cascading operation.

NOTE: This is a highly destructive, irreversible operation. Deleted volumes cannot be recovered unless you have snapshots or backups stored elsewhere. Use with extreme caution.

You must identify the VDC either with `--datacenter-id` (delete one) or with `--all` (delete every VDC on the account) - the two are mutually exclusive. Combine with `--force` (`-f`) to skip the interactive confirmation prompt (useful in scripts) and `--wait` (`-w`) to block until deletion completes.

Required values to run command:

* Data Center Id (or --all)

## Options

```text
  -a, --all                    Delete every Virtual Data Center on the account (and all resources inside them). Mutually exclusive with --datacenter-id
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [DatacenterId Name Location CpuFamily IPv6CidrBlock State Description Version Features SecAuthProtection]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -i, --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
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
# Delete a single VDC (prompts for confirmation)
ionosctl compute datacenter delete --datacenter-id DATACENTER_ID

# Delete a VDC non-interactively and wait for it to finish
ionosctl compute datacenter delete --datacenter-id DATACENTER_ID --force --wait
```

