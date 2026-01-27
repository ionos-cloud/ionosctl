---
description: "Create a Server"
---

# ServerCreate

## Usage

```text
ionosctl server create [flags]
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

Use this command to create an ENTERPRISE, CUBE, VCPU or GPU Server in a specified Virtual Data Center.

1. For ENTERPRISE Servers:

You need to set the number of cores for the Server and the amount of memory for the Server to be set. The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.`--ram 256` equals 256MB.
* providing both the value and the unit, e.g.`--ram 1GB`.

To see which CPU Family are available in which location, use `ionosctl location` commands.

Required values to create a Server of type ENTERPRISE:

* Data Center Id
* Cores
* RAM

2. For CUBE Servers:

Servers of type CUBE will be created with a Direct Attached Storage with the size set from the Template. To see more details about the available Templates, use `ionosctl template` commands.

Required values to create a Server of type CUBE:

* Data Center Id
* Type
* Template Id

3. For VCPU Servers:

You need to set the number of cores for the Server and the amount of memory for the Server to be set. The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.`--ram 256` equals 256MB.
* providing both the value and the unit, e.g.`--ram 1GB`.

You cannot set the CPU Family for VCPU Servers.

Required values to create a Server of type VCPU:

* Data Center Id
* Type
* Cores
* RAM

4. For GPU Servers:

Servers of type GPU will be created with a Direct Attached Storage with the size set from the Template. To see more details about the available Templates, use `ionosctl template` commands.

GPU servers do not support the --cpu-family flag and are automatically assigned the AMD_TURIN CPU family.

Required values to create a Server of type GPU:

* Data Center Id
* Type
* Template Id

By default, Licence Type for Direct Attached Storage is set to LINUX. You can set it using the `--licence-type` option or set an Image Id. For Image Id, it is needed to set a password or SSH keys.

You can wait for the Request to be executed using `--wait-for-request` option. You can also wait for Server to be in AVAILABLE state using `--wait-for-state` option. It is recommended to use both options together for this command.

## Options

```text
  -u, --api-url string             Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -z, --availability-zone string   Availability zone of the Server (default "AUTO")
      --bus string                 [CUBE Server] The bus type of the Direct Attached Storage (default "VIRTIO")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ServerId DatacenterId Name AvailabilityZone Cores RAM CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId NicMultiQueue] (default [ServerId,Name,Type,AvailabilityZone,Cores,RAM,CpuFamily,VmState,State])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                  The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits (required) (default 2)
      --cpu-family string          CPU Family for the Server. For CUBE Servers, the CPU Family is INTEL_SKYLAKE. If the flag is not set, the CPU Family will be chosen based on the location of the Datacenter. It will always be the first CPU Family available, as returned by the API (default "AUTO")
      --datacenter-id string       The unique Data Center Id (required)
  -D, --depth int                  Level of detail for response objects (default 1)
  -F, --filters strings            Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
  -a, --image-alias string         [CUBE Server] The Image Alias to use instead of Image Id for the Direct Attached Storage
      --image-id string            [CUBE Server] The Image Id or snapshot Id to be used as for the Direct Attached Storage
  -l, --licence-type string        [CUBE Server] Licence Type of the Direct Attached Storage. Can be one of: LINUX, RHEL, WINDOWS, WINDOWS2016, WINDOWS2019, WINDOWS2022, WINDOWS2025, UNKNOWN, OTHER (default "LINUX")
      --limit int                  Maximum number of items to return per request (default 50)
  -n, --name string                Name of the Server (default "Unnamed Server")
      --nic-multi-queue            Enable NIC Multi Queue to improve NIC throughput; changing this setting restarts the server. Not supported for CUBEs
      --no-headers                 Don't print table headers when table output is used
      --offset int                 Number of items to skip before starting to collect the results
      --order-by string            Property to order the results by
  -o, --output string              Desired output format [text|json|api-json] (default "text")
  -p, --password string            [CUBE Server] Initial image password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9
      --promote-volume             For CUBE and GPU servers, promotes the attached volume to be the Boot Volume. Requires --wait-for-state
      --query string               JMESPath query string to filter the output
  -q, --quiet                      Quiet output
      --ram string                 The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB (required)
  -k, --ssh-key-paths strings      [CUBE Server] Absolute paths for the SSH Keys of the Direct Attached Storage
      --template-id string         [CUBE Server] The unique Template Id (required)
  -t, --timeout int                Timeout option for Request for Server creation/for Server to be in AVAILABLE state [seconds] (default 60)
      --type string                Type usages for the Server. Can be one of: ENTERPRISE, CUBE, VCPU, GPU (default "ENTERPRISE")
  -v, --verbose count              Increase verbosity level [-v, -vv, -vvv]
  -N, --volume-name string         [CUBE Server] Name of the Direct Attached Storage (default "Unnamed Direct Attached Storage")
  -w, --wait-for-request           Wait for the Request for Server creation to be executed
  -W, --wait-for-state             Wait for new Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server create --datacenter-id DATACENTER_ID --cores 2 --ram 512MB -w -W

ionosctl server create --datacenter-id DATACENTER_ID --type VCPU --cores CORES --ram RAM

ionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID

ionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID --licence-type LICENCE_TYPE -w -W

ionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID --image-id IMAGE_ID --password IMAGE_PASSWORD -w -W

ionosctl server create --datacenter-id DATACENTER_ID --type GPU --template-id TEMPLATE_ID
```

