---
description: "Delete a Server (attached Volumes are left behind)"
---

# ServerDelete

## Usage

```text
ionosctl compute server delete [flags]
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

Use this command to permanently delete a Server from a Virtual Data Center, or every Server in a datacenter with --all.

NOTE: Deleting a Server does NOT delete its attached storage Volumes — they remain in the datacenter (and continue to be billed) so their data is preserved. Delete them separately if you no longer need them. For a CUBE server, its bundled DAS boot volume is removed together with the server.

Use `--wait` (`-w`) to wait for the resource to reach AVAILABLE state. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id

## Options

```text
  -a, --all                    Delete every Server in the given Virtual Data Center (--datacenter-id). Their attached Volumes are still left behind
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
ionosctl compute server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID
ionosctl compute server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --force
```

