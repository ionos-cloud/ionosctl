---
description: "Update a LAN"
---

# LanUpdate

## Usage

```text
ionosctl lan update [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified LAN. You can update the name, the public option for LAN and the Pcc Id to connect the LAN to a Cross-Connect.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* LAN Id

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId IPv6CidrBlock State DatacenterId] (default [LanId,Name,Public,PccId,IPv6CidrBlock,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ipv6-cidr string       The /64 IPv6 Cidr as defined in RFC 4291. It needs to be within the Datacenter IPv6 Cidr Block range. It can also be set to "AUTO" or "DISABLE". NOTE: Using an explicit Cidr to update the resource is not fully supported yet. (default "DISABLE")
  -i, --lan-id string          The unique LAN Id (required)
  -n, --name string            The name of the LAN
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --pcc-id string          The unique Id of the Cross-Connect the LAN will connect to
      --public                 Public option for LAN. E.g.: --public=true, --public=false
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout option for Request for LAN update [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait                   Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl lan update --datacenter-id DATACENTER_ID --lan-id LAN_ID --name NAME --public=false
```

