---
description: "Create a NIC"
---

# NicCreate

## Usage

```text
ionosctl nic create [flags]
```

## Aliases

For `nic` command:

```text
[n]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create/add a new NIC to the target Server. You can specify the name, ips, dhcp and Lan Id the NIC will sit on. If the Lan Id does not exist it will be created.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run a command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Name Dhcp LanId Ips IPv6Ips State FirewallActive FirewallType DeviceNumber PciSlot Mac DHCPv6 IPv6CidrBlock] (default [NicId,Name,Dhcp,LanId,Ips,IPv6Ips,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
      --dhcp                   Set to false if you wish to disable DHCP on the NIC. E.g.: --dhcp=true, --dhcp=false (default true)
      --dhcpv6                 Set to false if you wish to disable DHCPv6 on the NIC. E.g.: --dhcpv6=true, --dhcpv6=false (default true)
      --firewall-active        Activate or deactivate the Firewall. E.g.: --firewall-active=true, --firewall-active=false
      --firewall-type string   The type of Firewall Rules that will be allowed on the NIC (default "INGRESS")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ips strings            IPs assigned to the NIC. This can be a collection
      --ipv6-cidr string       The /80 IPv6 Cidr as defined in RFC 4291. It needs to be within the LAN IPv6 Cidr Block range. (default "disable")
      --ipv6-ips strings       IPv6 IPs assigned to the NIC. They need to be within the NIC IPv6 Cidr Block.
      --lan-id int             The LAN ID the NIC will sit on. If the LAN ID does not exist it will be created (default 1)
  -n, --name string            The name of the NIC (default "Internet Access")
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for NIC creation [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for the Request for NIC creation to be executed
```

## Examples

```text
ionosctl nic create --datacenter-id DATACENTER_ID --server-id SERVER_ID --name NAME
```

