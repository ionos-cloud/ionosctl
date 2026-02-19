---
description: "Attach a NIC to a Load Balancer"
---

# LoadbalancerNicAttach

## Usage

```text
ionosctl compute loadbalancer nic attach [flags]
```

## Aliases

For `loadbalancer` command:

```text
[lb]
```

For `nic` command:

```text
[n]
```

For `attach` command:

```text
[a]
```

## Description

Use this command to associate a NIC to a Load Balancer, enabling the NIC to participate in load-balancing.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Load Balancer Id
* NIC Id

## Options

```text
  -u, --api-url string           Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock] (default [NicId,Name,Dhcp,LanId,Ips,IPv6Ips,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string     The unique Data Center Id (required)
  -D, --depth int                Level of detail for response objects (default 1)
  -F, --filters strings          Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --limit int                Maximum number of items to return per request (default 50)
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -i, --nic-id string            The unique NIC Id (required)
      --no-headers               Don't print table headers when table output is used
      --offset int               Number of items to skip before starting to collect the results
      --order-by string          Property to order the results by
  -o, --output string            Desired output format [text|json|api-json] (default "text")
      --query string             JMESPath query string to filter the output
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion
  -t, --timeout int              Timeout option for Request for NIC attachment [seconds] (default 60)
  -v, --verbose count            Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request         Wait for the Request for NIC attachment to be executed
```

## Examples

```text
ionosctl compute loadbalancer nic attach --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --nic-id NIC_ID
```

