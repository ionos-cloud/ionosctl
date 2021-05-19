---
description: Create a Data Center
---

# DatacenterCreate

## Usage

```text
ionosctl datacenter create [flags]
```

## Description

Use this command to create a Virtual Data Center. You can specify the name, description or location for the object.

Virtual Data Centers (VDCs) are the foundation of the IONOS platform. VDCs act as logical containers for all other objects you will be creating, e.g. servers. You can provision as many Data Centers as you want. Data Centers have their own private network and are logically segmented from each other to create isolation.

You can wait for the Request to be executed using `--wait-for-request` option.

## Options

```text
  -u, --api-url string                  Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings                    Columns to be printed in the standard output (default [DatacenterId,Name,Location,Features,State])
  -c, --config string                   Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-description string   Description of the Data Center
      --datacenter-location string      Location for the Data Center (default "de/txl")
      --datacenter-name string          Name of the Data Center
      --force                           Force command to execute without user input
  -h, --help                            help for create
  -o, --output string                   Desired output format [text|json] (default "text")
  -q, --quiet                           Quiet output
      --timeout int                     Timeout option for Request for Data Center creation [seconds] (default 60)
      --wait-for-request                Wait for the Request for Data Center creation to be executed
```

## Examples

```text
ionosctl datacenter create --datacenter-name demoDatacenter --datacenter-location us/las
DatacenterId                           Name             Location
f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d   demoDatacenter   us/las
RequestId: 98ab8148-96c4-4091-90e8-9ee2b8a172f4
Status: Command datacenter create has been successfully executed

ionosctl datacenter create --datacenter-name demoDatacenter --datacenter-location gb/lhr --wait-for-state 
1.2s Waiting for request... DONE
DatacenterId                           Name             Location
8e543958-04f5-4872-bbf3-b28d46393ac7   demoDatacenter   gb/lhr
RequestId: 2401b498-8afb-4728-a22a-d2b26f5e31c3
Status: Command datacenter create & wait have been successfully executed
```

