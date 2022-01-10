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
                             Available columns: [DatacenterId Name Location State Description Version Features SecAuthProtection] (default [DatacenterId,Name,Location,Features,State])
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -d, --description string   Description of the Data Center
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
  -l, --location string      Location for the Data Center (default "de/txl")
  -n, --name string          Name of the Data Center (default "Unnamed Data Center")
  -o, --output string        Desired output format [text|json] (default "text")
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout option for Request for Data Center creation [seconds] (default 60)
  -v, --verbose              Print step-by-step process when running command
  -w, --wait-for-request     Wait for the Request for Data Center creation to be executed
```

## Examples

```text
ionosctl datacenter create --name NAME --location LOCATION_ID

ionosctl datacenter create --name NAME --location LOCATION_ID --wait-for-request
```

