---
description: "Update the compute shape or boot configuration of a Server"
---

# ServerUpdate

## Usage

```text
ionosctl compute server update [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `update` command:

```text
[u up]
```

## Description

Use this command to update a Server in a Virtual Data Center. You can rescale an ENTERPRISE/VCPU server (--cores, --ram), change its --cpu-family or --availability-zone, and repoint what it boots from (--volume-id for a boot Volume, --cdrom-id for a boot CD-ROM). Both boot targets must already be attached to the Server.

RAM (--ram) must be a multiple of 256. The default unit is MB, so --ram 256 = 256MB; a unit may be given, e.g. --ram 1GB. Minimum 256MB; the maximum depends on your contract limit.

CUBE constraint: for CUBE Servers only the Name can be updated — their cores, RAM and CPU family are fixed by the instance-size template and cannot be changed here. Some changes (e.g. --nic-multi-queue, CPU family) require a server restart to take effect.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string             Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -z, --availability-zone string   Physical availability zone of the Server: AUTO (platform-chosen), ZONE_1 or ZONE_2
      --cdrom-id string            Id of the CD-ROM to set as the Server's boot device. The CD-ROM must already be attached to this Server. Use this to boot from an installer ISO instead of a Volume (--volume-id)
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ServerId Name Type AvailabilityZone Cores RAM CpuFamily VmState State DatacenterId TemplateId BootCdromId BootVolumeId NicMultiQueue EnabledFeatures]
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                  New number of CPU cores (ENTERPRISE/VCPU only; fixed by the template for CUBE/GPU). Maximum depends on your contract resource limits (default 2)
      --cpu-family string          New CPU family for the Server, e.g. INTEL_SKYLAKE, INTEL_XEON, AMD_OPTERON. Availability depends on the datacenter location; changing it requires a server restart. Not applicable to CUBE/VCPU/GPU
      --datacenter-id string       The unique Data Center Id (required)
  -D, --depth int                  Level of detail for response objects (default 1)
  -F, --filters strings            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
      --limit int                  Maximum number of items to return per request (default 50)
  -n, --name string                New display name for the Server. This is the only property updatable on CUBE Servers
      --nic-multi-queue            Enable NIC Multi Queue to improve NIC throughput; changing this setting restarts the server. Not supported for CUBEs
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Number of items to skip before starting to collect the results
      --order-by string            Property to order the results by
  -o, --output string              Desired output format [text|json|api-json] (default "text")
      --query string               JMESPath query string to filter the output
  -q, --quiet                      Quiet output
      --ram string                 New memory size, in multiples of 256. Default unit is MB (e.g. --ram 256 = 256MB); a unit may be given (e.g. --ram 1GB). Minimum 256MB, maximum per contract limit. ENTERPRISE/VCPU only
  -i, --server-id string           The unique Server Id (required)
  -t, --timeout int                Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
      --volume-id string           Id of the Volume to set as the Server's boot volume. The Volume must already be attached to this Server. Mutually exclusive with booting from a CD-ROM (--cdrom-id)
  -w, --wait                       Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# BASIC: rescale an ENTERPRISE server to 4 cores
ionosctl compute server update --datacenter-id DATACENTER_ID --server-id SERVER_ID --cores 4

# ADVANCED: rescale cores + RAM and set the boot volume to an already-attached Volume, waiting until AVAILABLE
ionosctl compute server update --datacenter-id DATACENTER_ID --server-id SERVER_ID --cores 8 --ram 16GB --volume-id VOLUME_ID --wait
```

