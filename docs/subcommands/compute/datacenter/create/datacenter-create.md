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
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -d, --description string   Description of the Data Center
  -l, --location string      Location for the Data Center (default "de/txl")
  -n, --name string          Name of the Data Center (default "Unnamed Data Center")
  -t, --timeout int          Timeout option for Request for Data Center creation [seconds] (default 60)
  -w, --wait-for-request     Wait for the Request for Data Center creation to be executed
```

## Examples

```text
ionosctl datacenter create --name NAME --location LOCATION_ID
ionosctl datacenter create --name NAME --location LOCATION_ID --wait-for-request
```

