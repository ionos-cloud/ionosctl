---
description: "Delete a Server"
---

# ServerDelete

## Usage

```text
ionosctl server delete [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Server from a Virtual Data Center.

NOTE: This will not automatically remove the storage Volumes attached to a Server.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -a, --all                    Delete all Servers form a virtual Datacenter.
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ServerId DatacenterId Name AvailabilityZone Cores RAM CpuFamily VmState State TemplateId Type BootCdromId BootVolumeId NicMultiQueue] (default [ServerId,Name,Type,AvailabilityZone,Cores,RAM,CpuFamily,VmState,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Server deletion [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for Server deletion to be executed
```

## Examples

```text
ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID

ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --force
```

