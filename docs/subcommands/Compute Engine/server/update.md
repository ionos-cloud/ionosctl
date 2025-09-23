---
description: "Update a Server"
---

# ServerUpdate

## Usage

```text
ionosctl server update [flags]
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

Use this command to update a specified Server from a Virtual Data Center.

You can set the RAM size in the following ways:

* providing only the value, e.g.`--ram 256` equals 256MB.
* providing both the value and the unit, e.g.`--ram 1GB`.

The amount of memory for the Server must be specified in multiples of 256. The default unit is MB. Minimum: 256MB. Maximum: it depends on your contract limit.

Note: For CUBE Servers, only Name attribute can be updated.

You can wait for the Request to be executed using `--wait-for-request` option. You can also wait for Server to be in AVAILABLE state using `--wait-for-state` option. It is recommended to use both options together for this command.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string             Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -z, --availability-zone string   Availability zone of the Server
      --cdrom-id string            The unique Cdrom Id for the BootCdrom. The Cdrom needs to be already attached to the Server
      --cols strings               Set of columns to be printed on output 
                                   Available columns: [ServerId DatacenterId Name AvailabilityZone Cores RAM CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId NicMultiQueue] (default [ServerId,Name,Type,AvailabilityZone,Cores,RAM,CpuFamily,VmState,State])
  -c, --config string              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                  The total number of cores for the Server, e.g. 4. Maximum: depends on contract resource limits (default 2)
      --cpu-family string          CPU Family of the Server
      --datacenter-id string       The unique Data Center Id (required)
  -D, --depth int32                Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                      Force command to execute without user input
  -h, --help                       Print usage
  -n, --name string                Name of the Server
      --nic-multi-queue            Enable NIC Multi Queue to improve NIC throughput; changing this setting restarts the server. Not supported for CUBEs
      --no-headers                 Don't print table headers when table output is used
  -o, --output string              Desired output format [text|json|api-json] (default "text")
  -q, --quiet                      Quiet output
      --ram string                 The amount of memory for the Server. Size must be specified in multiples of 256. e.g. --ram 256 or --ram 256MB
  -i, --server-id string           The unique Server Id (required)
  -t, --timeout int                Timeout option for Request for Server update/for Server to be in AVAILABLE state [seconds] (default 60)
  -v, --verbose                    Print step-by-step process when running command
      --volume-id string           The unique Volume Id for the BootVolume. The Volume needs to be already attached to the Server
  -w, --wait-for-request           Wait for the Request for Server update to be executed
  -W, --wait-for-state             Wait for the updated Server to be in AVAILABLE state
```

## Examples

```text
ionosctl server update --datacenter-id DATACENTER_ID --server-id SERVER_ID --cores 4
```

