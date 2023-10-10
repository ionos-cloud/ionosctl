---
description: "Update a NIC"
---

# NicUpdate

## Usage

```text
ionosctl nic update [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update the configuration of a specified NIC. Some restrictions are in place: The primary address of a NIC connected to a Load Balancer can only be changed by changing the IP of the Load Balancer. You can also add additional reserved, public IPs to the NIC.

The user can specify and assign private IPs manually. Valid IP addresses for private networks are 10.0.0.0/8, 172.16.0.0/12 or 192.168.0.0/16.

The value for firewallActive can be toggled between true and false to enable or disable the firewall. When the firewall is enabled, incoming traffic is filtered by all the firewall rules configured on the NIC. When the firewall is disabled, then all incoming traffic is routed directly to the NIC and any configured firewall rules are ignored.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* NIC Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock] (default [NicId,Name,Dhcp,LanId,Ips,IPv6Ips,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                   Boolean value that indicates if the NIC is using DHCP (true) or not (false). E.g.: --dhcp=true, --dhcp=false (default true)
      --dhcpv6                 Set to false if you wish to disable DHCPv6 on the NIC. E.g.: --dhcpv6=true, --dhcpv6=false (default true)
      --firewall-active        Activate or deactivate the Firewall. E.g.: --firewall-active=true, --firewall-active=false
      --firewall-type string   The type of Firewall Rules that will be allowed on the NIC (default "INGRESS")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            IPs assigned to the NIC
      --ipv6-cidr string       The /80 IPv6 Cidr as defined in RFC 4291. It needs to be within the LAN IPv6 Cidr Block range. (default "disable")
      --ipv6-ips strings       IPv6 IPs assigned to the NIC. They need to be within the NIC IPv6 Cidr Block.
      --lan-id int             The LAN ID the NIC sits on (default 1)
  -n, --name string            The name of the NIC
  -i, --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for NIC update [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for NIC update to be executed
```

## Examples

```text
ionosctl nic update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --wait-for-request
```

