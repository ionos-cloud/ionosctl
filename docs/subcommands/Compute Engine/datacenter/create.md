---
description: "Create a Data Center"
---

# DatacenterCreate

## Usage

```text
ionosctl datacenter create [flags]
```

## Aliases

For `datacenter` command:

```text
[d dc vdc]
```

For `create` command:

```text
[c]
```

## Description

Use this command to create a Virtual Data Center. You can specify the name, description or location for the object.

Virtual Data Centers are the foundation of the IONOS platform. VDCs act as logical containers for all other objects you will be creating, e.g. servers. You can provision as many Data Centers as you want. Data Centers have their own private network and are logically segmented from each other to create isolation.

You can wait for the Request to be executed using `--wait-for-request` option.

## Options

```text
  -u, --api-url string       Override default host url (default "https://api.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [DatacenterId Name Location State Description Version Features CpuFamily SecAuthProtection IPv6CidrBlock] (default [DatacenterId,Name,Location,CpuFamily,IPv6CidrBlock,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string   Description of the Data Center
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -l, --location string      Location for the Data Center (default "de/txl")
  -n, --name string          Name of the Data Center (default "Unnamed Data Center")
      --no-headers           Don't print table headers when table output is used
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for polling the request (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait                 Polls the request continuously until the operation is completed
```

## Examples

```text
ionosctl datacenter create --name NAME --location LOCATION_ID
ionosctl datacenter create --name NAME --location LOCATION_ID --wait-for-request
```

