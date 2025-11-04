---
description: "Remove IP Failover group from a LAN"
---

# IpfailoverRemove

## Usage

```text
ionosctl ipfailover remove [flags]
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

Use this command to remove an IP Failover group from a LAN.

Required values to run command:

* Data Center Id
* Lan Id
* Server Id
* Nic Id
* IP address

## Options

```text
  -a, --all                    Remove all IP Failovers.
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Ip] (default [NicId,Ip])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ip ip                  Allocated IP (required)
      --lan-id string          The unique LAN Id (required)
      --limit int              Pagination limit: Maximum number of items to return per request (default 50)
      --nic-id string          The unique NIC Id (required)
      --no-headers             Don't print table headers when table output is used
      --offset int             Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for IP Failover deletion [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for IP Failover deletion to be executed
```

## Examples

```text
ionosctl ipfailover remove --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --ip "x.x.x.x"
```

