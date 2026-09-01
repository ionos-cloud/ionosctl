---
description: "Force a hard reboot (power-cycle) of a Server"
---

# ServerReboot

## Usage

```text
ionosctl compute server reboot [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `reboot` command:

```text
[r]
```

## Description

Use this command to force a HARD reboot of the Server: the equivalent of pulling the power and turning the machine back on. This does NOT ask the guest OS to shut down cleanly, so unsaved in-guest state can be lost. Use it only when the machine is unresponsive or you specifically need a power-cycle; for a graceful restart, reboot from inside the operating system instead.

Unlike `server stop`, a reboot keeps the compute resources allocated (billing continues) and the server stays powered on afterwards; it does not deallocate the machine.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ServerId Name Type AvailabilityZone Cores RAM CpuFamily VmState State DatacenterId TemplateId BootCdromId BootVolumeId NicMultiQueue EnabledFeatures]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
  -i, --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                   Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute server reboot --datacenter-id DATACENTER_ID --server-id SERVER_ID
```

