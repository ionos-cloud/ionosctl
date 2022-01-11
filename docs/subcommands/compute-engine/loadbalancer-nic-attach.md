---
description: Attach a NIC to a Load Balancer
---

# LoadbalancerNicAttach

## Usage

```text
ionosctl loadbalancer nic attach [flags]
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
  -u, --api-url string           Override default host url (default "https://api.ionos.com")
      --cols strings             Set of columns to be printed on output 
                                 Available columns: [NicId Name Dhcp LanId Ips State FirewallActive Mac] (default [NicId,Name,Dhcp,LanId,Ips,State])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string     The unique Data Center Id (required)
  -f, --force                    Force command to execute without user input
  -h, --help                     Print usage
      --loadbalancer-id string   The unique Load Balancer Id (required)
  -i, --nic-id string            The unique NIC Id (required)
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --server-id string         The unique Server Id on which NIC is build on. Not required, but it helps in autocompletion
  -t, --timeout int              Timeout option for Request for NIC attachment [seconds] (default 60)
  -v, --verbose                  Print step-by-step process when running command
  -w, --wait-for-request         Wait for the Request for NIC attachment to be executed
```

## Examples

```text
ionosctl loadbalancer nic attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --loadbalancer-id LOADBALANCER_ID
```

