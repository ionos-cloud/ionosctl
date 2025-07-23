---
description: "List IP Failovers groups from a LAN"
---

# IpfailoverList

## Usage

```text
ionosctl ipfailover list [flags]
```

## Aliases

For `ipfailover` command:

```text
[ipf]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get a list of IP Failovers groups from a LAN.

Required values to run command:

* Data Center Id
* Lan Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud'|'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [NicId Ip] (default [NicId,Ip])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10. (default 1)
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --lan-id string          The unique LAN Id (required)
  -M, --max-results int32      The maximum number of elements to return
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl ipfailover list --datacenter-id DATACENTER_ID --lan-id LAN_ID
```

