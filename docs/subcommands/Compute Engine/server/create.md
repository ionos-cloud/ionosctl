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

Use this command to create an ENTERPRISE or CUBE Server in a specified Virtual Data Center.

* For ENTERPRISE Servers:

You need to set the number of cores for the Server and the amount of memory for the Server to be set. The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit. You can set the RAM size in the following ways:

* providing only the value, e.g.`--ram 256` equals 256MB.
* providing both the value and the unit, e.g.`--ram 1GB`.

To see which CPU Family are available in which location, use `ionosctl location` commands.

Required values to create a Server of type ENTERPRISE:

* Data Center Id
* Cores
* RAM

* For CUBE Servers:

Servers of type CUBE will be created with a Direct Attached Storage with the size set from the Template. To see more details about the available Templates, use `ionosctl template` commands.

Required values to create a Server of type CUBE:

* Data Center Id
* Type
* Template Id

By default, Licence Type for Direct Attached Storage is set to LINUX. You can set it using the `--licence-type` option or set an Image Id. For Image Id, it is needed to set a password or SSH keys.

You can wait for the Request to be executed using `--wait-for-request` option. You can also wait for Server to be in AVAILABLE state using `--wait-for-state` option. It is recommended to use both options together for this command.

## Options

```text
  -u, --api-url string             Override default host url (default "https://api.ionos.com")
  -z, --availability-zone string   Availability zone of the Server (default "AUTO")
      --bus string                 [CUBE Server] The bus type of the Direct Attached Storage (default "VIRTIO")
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ServerId DatacenterId Name AvailabilityZone Cores Ram CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId] (default [ServerId,Name,Type,AvailabilityZone,Cores,Ram,CpuFamily,VmState,State])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int                  The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits (required) (default 2)
      --cpu-family string          CPU Family for the Server. For CUBE Servers, the CPU Family is INTEL_SKYLAKE (default "AMD_OPTERON")
      --datacenter-id string       The unique Data Center Id (required)
  -D, --depth int32                Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
  -a, --image-alias string         [CUBE Server] The Image Alias to use instead of Image Id for the Direct Attached Storage
      --image-id string            [CUBE Server] The Image Id or snapshot Id to be used as for the Direct Attached Storage
  -l, --licence-type string        [CUBE Server] Licence Type of the Direct Attached Storage (default "LINUX")
  -n, --name string                Name of the Server (default "Unnamed Server")
  -o, --output string              Desired output format [text|json] (default "text")
  -p, --password string            [CUBE Server] Initial image password to be set for installed OS. Works with public Images only. Not modifiable. Password rules allows all characters from a-z, A-Z, 0-9
  -q, --quiet                      Quiet output
      --ram string                 The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB (required)
  -k, --ssh-key-paths strings      [CUBE Server] Absolute paths for the SSH Keys of the Direct Attached Storage
      --template-id string         [CUBE Server] The unique Template Id (required)
  -t, --timeout int                Timeout option for Request for Server creation/for Server to be in AVAILABLE state [seconds] (default 60)
      --type string                Type usages for the Server (default "ENTERPRISE")
  -v, --verbose                    Print step-by-step process when running command
  -N, --volume-name string         [CUBE Server] Name of the Direct Attached Storage (default "Unnamed Direct Attached Storage")
  -w, --wait-for-request           Wait for the Request for Server creation to be executed
  -W, --wait-for-state             Wait for new Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server create --datacenter-id DATACENTER_ID --cores 2 --ram 512MB -w -W

ionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID --licence-type LICENCE_TYPE -w -W

ionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID --image-id IMAGE_ID --password IMAGE_PASSWORD -w -W
```

