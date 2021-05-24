---
description: List Data Centers
---

# DatacenterList

## Usage

```text
ionosctl datacenter list [flags]
```

## Aliases

For `datacenter` command:
```text
[d dc]
```

For `list` command:
```text
[l ls]
```

## Description

Use this command to retrieve a complete list of Virtual Data Centers provisioned under your account.

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [DatacenterId Name Location State Description Version Features SecAuthProtection] (default [DatacenterId,Name,Location,Features,State])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             help for list
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
```

## Examples

```text
ionosctl datacenter list 
DatacenterId                           Name             Location
ff279ffd-ac61-4e5d-ba5e-058296c77774   demoDatacenter   us/las

ionosctl datacenter list --cols "DatacenterId,Name,Location,Version"
DatacenterId                           Name             Location   Version
ff279ffd-ac61-4e5d-ba5e-058296c77774   demoDatacenter   us/las     1
```

