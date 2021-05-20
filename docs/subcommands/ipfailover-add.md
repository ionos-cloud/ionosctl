---
description: Add IP Failover group to a LAN
---

# IpfailoverAdd

## Usage

```text
ionosctl ipfailover add [flags]
```

## Description

Use this command to add an IP Failover group to a LAN. 

Successfully setting up an IP Failover group requires three steps:

* Add a reserved IP address to a NIC that will become the IP Failover master.
* Use `ionosctl ipfailover add` command to enable IP Failover by providing the relevant IP and NIC Id values.
* Add the same reserved IP address to any other NICs that are a member of the same LAN. Those NICs will become IP Failover members.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -f, --force                  Force command to execute without user input
  -F, --format strings         Collection of fields to be printed on output (default [NicId,Ip])
  -h, --help                   help for add
      --ip string              IP address to be added to IP Failover Group (required)
      --lan-id string          The unique LAN Id (required)
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for IP Failover creation [seconds] (default 60)
  -w, --wait-for-request       Wait for the Request for IP Failover creation to be executed
```

## Examples

```text
ionosctl ipfailover add --datacenter-id 2c08a329-dbe3-427a-8ef9-897e620fef3d --server-id 11c8ac02-224b-4bd0-833c-196719860fc1 --lan-id 1 --nic-id 5662f39c-b7cb-4840-b6ab-ae43cd0202cc --ip "x.x.x.x"
NicId                                  Ip
5662f39c-b7cb-4840-b6ab-ae43cd0202cc   x.x.x.x
```

