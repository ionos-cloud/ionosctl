---
description: "Create a LAN"
---

# LanCreate

## Usage

```text
ionosctl compute lan create [flags]
```

## Aliases

For `lan` command:

```text
[l]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a new LAN inside a Virtual Data Center (VDC). A LAN is a virtual network segment that connects the NICs attached to it within a single datacenter (`--datacenter-id`).

Decide at creation time whether the LAN is:
  * private (`--public=false`, the default): internal, datacenter-local traffic only, no direct internet route. This is the type required to later join a Cross-Connect.
  * public (`--public=true`): attached to an internet gateway; NICs on it can send/receive public internet traffic and be assigned public IPv4 addresses.

Optionally pass `--pcc` to immediately attach the LAN to a Cross-Connect (Private Cross-Connect), bridging it with private LANs in other VDCs in the same region. The LAN must be private and its IP range must not overlap the other members of the Cross-Connect.

NOTE: IP failover groups (a reserved IP that floats between servers for high availability) are NOT set here; configure them on a NIC after the LAN exists.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [LanId Name Public PccId IPv6CidrBlock State DatacenterId]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --ipv6-cidr string       The /64 IPv6 Cidr as defined in RFC 4291. It needs to be within the Datacenter IPv6 Cidr Block range. It can also be set to "AUTO" or "DISABLE". Use "DISABLE" (default) to keep the LAN IPv4-only, "AUTO" to let IONOS assign a /64 block automatically from the datacenter's range, or pass an explicit /64 block within the datacenter's IPv6 range. (default "DISABLE")
      --limit int              Maximum number of items to return per request (default 50)
  -n, --name string            A human-friendly name for the LAN. Not required to be unique (default "Unnamed LAN")
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --pcc-id string          ID of the Cross-Connect (Private Cross-Connect) to attach this LAN to, bridging it with private LANs in other VDCs of the same region. The LAN must be private (--public=false) and its IP range must not overlap the other Cross-Connect members
  -p, --public                 Whether the LAN is public. true = attached to an internet gateway so NICs can reach the internet and be assigned public IPv4 addresses; false (default) = private, internal datacenter traffic only. A Cross-Connect (--pcc) requires a private LAN. E.g.: --public=true
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Create a simple private LAN (internal traffic only)
ionosctl compute lan create --datacenter-id DATACENTER_ID --name "backend"

# Create a public LAN (internet-facing, NICs can receive public IPs)
ionosctl compute lan create --datacenter-id DATACENTER_ID --name "frontend" --public=true

# Create a private LAN and attach it to a Cross-Connect to bridge it with LANs in other VDCs (Cross-Connect requires a private LAN)
ionosctl compute lan create --datacenter-id DATACENTER_ID --name "cross-dc" --public=false --pcc PCC_ID
```

