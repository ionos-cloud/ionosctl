---
description: "Remove IP Failover group from a LAN"
---

# IpfailoverRemove

## Usage

```text
ionosctl compute ipfailover remove [flags]
```

## Aliases

For `ipfailover` command:

```text
[ipf]
```

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a NIC from an IP failover group on a LAN. Identify the entry to remove by its --ip and --nic-id. Removing the master NIC leaves the remaining member NICs in the group; removing the last NIC effectively dissolves the group.

Use --all to clear every IP failover entry from the LAN at once.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address

## Options

```text
  -a, --all                    Remove all IP failover entries from the LAN. Mutually exclusive with --nic-id/--ip
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Ip]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique ID of the Virtual Data Center that holds the LAN and servers (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ip ip                  The floating IP of the failover entry to remove. Together with --nic-id this identifies the specific entry (required)
      --lan-id string          The unique ID of the LAN the failover group lives on (required)
      --limit int              Maximum number of items to return per request (default 50)
      --nic-id string          The unique ID of the NIC to remove from the failover group (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique ID of the server that owns the NIC being removed (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Remove a single NIC from a failover group
ionosctl compute ipfailover remove --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --ip "203.0.113.10"

# Remove all IP failover entries from a LAN
ionosctl compute ipfailover remove --datacenter-id DATACENTER_ID --lan-id LAN_ID --all
```

