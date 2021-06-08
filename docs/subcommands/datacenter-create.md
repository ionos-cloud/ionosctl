---
description: Create a Data Center
---

# DatacenterCreate

## Usage

```text
ionosctl datacenter create [flags]
```

## Aliases

For `datacenter` command:
```text
[d dc]
```

For `create` command:
```text
[c]
```

## Description

Use this command to create a Virtual Data Center. You can specify the name, description or location for the object.

Virtual Data Centers (VDCs) are the foundation of the IONOS platform. VDCs act as logical containers for all other objects you will be creating, e.g. servers. You can provision as many Data Centers as you want. Data Centers have their own private network and are logically segmented from each other to create isolation.

You can wait for the Request to be executed using `--wait-for-request` option.

## Options

```text
  -u, --api-url string       Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection] (default [DatacenterId,Name,Location,CpuFamily,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -d, --description string   Description of the Data Center
  -f, --force                Force command to execute without user input
  -h, --help                 help for create
  -l, --location string      Location for the Data Center (default "de/txl")
  -n, --name string          Name of the Data Center
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for Request for Data Center creation [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Data Center creation to be executed
```

## Examples

```text
ionosctl datacenter create --name NAME --location LOCATION_ID

ionosctl datacenter create --name NAME --location LOCATION_ID --wait-for-request
```

