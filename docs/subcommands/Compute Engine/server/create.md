---
description: "Create an ENTERPRISE, VCPU, CUBE or GPU Server in a Virtual Data Center"
---

# ServerCreate

## Usage

```text
ionosctl compute server create [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Server in a Virtual Data Center (--datacenter-id). The compute shape you must describe depends on the server --type, and this is the single most important relationship to get right:

  * ENTERPRISE (the default type): you size the machine explicitly. --cores and --ram are REQUIRED. --cpu-family is optional (see below).
  * VCPU: you size the machine explicitly with --cores and --ram (REQUIRED), but you CANNOT set --cpu-family — the platform selects it.
  * CUBE: a fixed bundle. --template-id is REQUIRED and defines the cores, RAM and included NVMe Direct Attached Storage (DAS) boot volume; do NOT pass --cores/--ram. A DAS boot volume is created and attached automatically.
  * GPU: a fixed bundle with an attached GPU. --template-id is REQUIRED; --cpu-family is not accepted (the AMD_TURIN family is assigned automatically). A DAS boot volume is created and attached automatically.

RAM (--ram): must be a multiple of 256. The default unit is MB, so --ram 256 means 256MB; you may also give a unit, e.g. --ram 1GB. Minimum 256MB; the maximum depends on your contract limits. Applies to ENTERPRISE and VCPU only.

CPU family (--cpu-family): for ENTERPRISE, values such as INTEL_SKYLAKE, INTEL_XEON and AMD_OPTERON are valid, but availability differs per datacenter location. Run `ionosctl compute location` to see which families a location offers. If you omit it, the first family available in the datacenter's location is chosen for you.

Boot media / OS: for CUBE and GPU the DAS boot volume's OS comes from --image-id / --image-alias (or a bare --licence-type, default LINUX, for an empty volume). When you boot from a PUBLIC image you must also set --password or --ssh-key-paths so you can log in; PRIVATE images (which already contain credentials) need neither. To boot from a CD-ROM/ISO instead, attach one after creation with `server cdrom attach`.

Confidential Computing (--confidential): creates an AMD SEV-SNP encrypted-memory VM. It is ENTERPRISE-only and requires --image-id pointing at a private, SEV-SNP-capable image; --cores and --cpu-family are NOT allowed because both are derived from the image's launch-config.json. Size its boot volume with --size and --storage-type.

Use `--wait` (`-w`) to wait for the Server to reach AVAILABLE state before the command returns.

## Options

```text
  -u, --api-url string                               Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -z, --availability-zone string                     Physical availability zone the Server is placed in. AUTO lets the platform choose; ZONE_1 / ZONE_2 pin it, letting you spread servers across zones for fault tolerance (default "AUTO")
      --bus string                                   [CUBE Server] Bus the DAS volume is exposed on. VIRTIO is faster and recommended for modern OSes; IDE maximises compatibility with older ones (default "VIRTIO")
      --cols strings                                 Set of columns to be printed on output 
                                                     Available columns: [ServerId Name Type AvailabilityZone Cores RAM CpuFamily VmState State DatacenterId TemplateId BootCdromId BootVolumeId NicMultiQueue EnabledFeatures]
      --confidential                                 Create a Confidential Computing (SEV-SNP) VM from a confidential boot image. Requires --type ENTERPRISE and --image-id (a private, SEV-SNP image). Do not set --cores or --cpu-family: both are derived from the image's launch-config.json. A boot volume is created from --image-id and attached automatically; size it with --size and --storage-type.
  -c, --config string                                Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                                    The number of CPU cores. REQUIRED for ENTERPRISE and VCPU; ignored for CUBE/GPU (fixed by --template-id) and for --confidential (derived from the image). Maximum depends on your contract resource limits (required) (default 2)
      --cpu-family ionosctl compute location         CPU family for ENTERPRISE Servers, e.g. INTEL_SKYLAKE, INTEL_XEON, AMD_OPTERON. Availability varies by datacenter location (see ionosctl compute location). Leave as AUTO to have the API pick the first family available in the datacenter's location. Not accepted for VCPU (platform-chosen), GPU (always AMD_TURIN) or --confidential (image-derived); the API also rejects it for CUBE (default "AUTO")
      --datacenter-id string                         The unique Data Center Id (required)
  -D, --depth int                                    Level of detail for response objects (default 1)
  -F, --filters strings                              Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                                        Force command to execute without user input
  -h, --help                                         Print usage
  -a, --image-alias string                           [CUBE Server] Human-friendly alias of a public image to install on the DAS boot volume (e.g. ubuntu:latest), as an alternative to --image-id. Public images require --password or --ssh-key-paths
      --image-id string                              [CUBE Server] Id of an image or snapshot to install on the DAS boot volume. Public images require --password or --ssh-key-paths; private images (and snapshots) carry their own credentials. Also used as the confidential boot image when --confidential is set
  -l, --licence-type string                          [CUBE Server] OS licence type of the DAS boot volume, used for an empty volume when no --image-id/--image-alias is given. Determines OS-specific billing/handling (e.g. LINUX, WINDOWS). Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "LINUX")
      --limit int                                    Maximum number of items to return per request (default 50)
  -n, --name string                                  Display name of the Server. Type-specific defaults are used if omitted (e.g. 'Unnamed Server', 'Unnamed Cube') (default "Unnamed Server")
      --nic-multi-queue                              Enable NIC Multi Queue to improve NIC throughput; changing this setting restarts the server. Not supported for CUBEs
      --no-headers                                   Don't print table headers when table output is used
      --offset int                                   Number of items to skip before starting to collect the results
      --order-by string                              Property to order the results by
  -o, --output string                                Desired output format [text|json|api-json] (default "text")
  -p, --password string                              [CUBE Server] Root/Administrator password to set on the installed OS. Applies to PUBLIC images only, is set once at creation (not modifiable afterwards), and accepts characters a-z, A-Z, 0-9. Provide this and/or --ssh-key-paths when booting a public image
      --promote-volume                               [CUBE/GPU Server] After creation, promote the auto-attached DAS volume to be the Server's boot volume. Requires --wait (the promotion is a follow-up PATCH once the server is AVAILABLE)
      --query string                                 JMESPath query string to filter the output
  -q, --quiet                                        Quiet output
      --ram string                                   Memory size, in multiples of 256. Default unit is MB (e.g. --ram 256 = 256MB); a unit may be given (e.g. --ram 2GB). Minimum 256MB, maximum per contract limit. REQUIRED for ENTERPRISE and VCPU; not used for CUBE/GPU (fixed by --template-id) (required)
      --size string                                  [Confidential] Size of the confidential boot volume. Default unit is GB, e.g. --size 10 or --size 10GB (default "10")
  -k, --ssh-key-paths strings                        [CUBE Server] Paths to SSH public key files to inject into the DAS boot volume's OS (public images only). Comma-separate multiple paths. An alternative or complement to --password for logging in
      --storage-type string                          [Confidential] Storage backing of the confidential boot volume. SSD tiers offer higher performance than HDD. Can be one of: HDD, SSD, SSD Standard, SSD Premium (default "HDD")
      --template-id ionosctl compute template list   [CUBE/GPU Server] Id of the instance-size Template that fixes cores, RAM and the included DAS boot volume. REQUIRED for --type CUBE and GPU. List available templates with ionosctl compute template list (required)
  -t, --timeout int                                  Timeout in seconds for --wait and other wait operations (default 600)
      --type string                                  Compute model of the Server, which decides the required sizing flags. ENTERPRISE/VCPU need --cores and --ram; CUBE/GPU need --template-id instead (see the long description). Can be one of: ENTERPRISE, CUBE, VCPU, GPU (default "ENTERPRISE")
  -v, --verbose count                                Increase verbosity level [-v, -vv, -vvv]
  -N, --volume-name string                           [CUBE Server] Display name of the included Direct Attached Storage (DAS) boot volume (default "Unnamed Direct Attached Storage")
  -w, --wait                                         Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# BASIC: an ENTERPRISE server with 2 cores and 2GB RAM (CPU family auto-selected for the datacenter's location)
ionosctl compute server create --datacenter-id DATACENTER_ID --name web-01 --cores 2 --ram 2GB

# ADVANCED: an ENTERPRISE server pinned to a CPU family, booting a public image with an SSH key, waiting until it is AVAILABLE
ionosctl compute server create --datacenter-id DATACENTER_ID --name db-01 --cores 4 --ram 8GB \
  --cpu-family INTEL_SKYLAKE --image-id IMAGE_ID --ssh-key-paths ~/.ssh/id_rsa.pub --wait

# CUBE server: sizing comes entirely from the template bundle (do not pass --cores/--ram)
ionosctl compute server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID

# VCPU server: explicit cores/RAM, no --cpu-family
ionosctl compute server create --datacenter-id DATACENTER_ID --type VCPU --cores 2 --ram 4GB

# GPU server: fixed template bundle; CPU family is set automatically
ionosctl compute server create --datacenter-id DATACENTER_ID --type GPU --template-id TEMPLATE_ID
```

