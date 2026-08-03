---
description: "Add IP Failover group to a LAN"
---

# IpfailoverAdd

## Usage

```text
ionosctl compute ipfailover add [flags]
```

## Aliases

For `ipfailover` command:

```text
[ipf]
```

For `add` command:

```text
[a]
```

## Description

Use this command to add a NIC to an IP failover group on a LAN, registering it with the group's floating public IP.

Setting up a working IP failover group takes three steps:

  1. Reserve a public IP block in the same region/location as the datacenter (`ionosctl ipblock create`) and assign one of its IPs to the NIC that will become the failover MASTER.
  2. Run this command with that IP and the master NIC's Id to create/enable the failover group.
  3. Assign the SAME reserved IP to the other NICs on the same LAN and run this command again for each of them. Those NICs join the group as MEMBERS (standby).

If the group does not exist yet on the LAN, the first `add` creates it; subsequent `add` calls with the same --ip extend it. The IP must belong to a reserved IP block, not an ad-hoc/DHCP address.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Ip]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique ID of the Virtual Data Center that holds the LAN and servers (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ip ip                  The floating public IP for the failover group. Must be an address from a reserved IP block in the same region/location as the datacenter. Reuse the exact same IP when adding member NICs (required)
      --lan-id string          The unique ID of the LAN the failover group lives on. All member NICs must be on this LAN (required)
      --limit int              Maximum number of items to return per request (default 50)
      --nic-id string          The unique ID of the NIC to add to the failover group. The first NIC added becomes the master; later NICs (with the same IP) become standby members (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique ID of the server that owns the NIC being added to the group (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Register the master NIC with the group's floating IP
ionosctl compute ipfailover add --datacenter-id DATACENTER_ID --server-id SERVER_ID --lan-id LAN_ID --nic-id MASTER_NIC_ID --ip "203.0.113.10"

# Add a standby member NIC (on another server) to the same group, reusing the same IP
ionosctl compute ipfailover add --datacenter-id DATACENTER_ID --server-id SERVER_ID_2 --lan-id LAN_ID --nic-id MEMBER_NIC_ID --ip "203.0.113.10"
```

