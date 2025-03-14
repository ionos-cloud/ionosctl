---
description: "Get a Label using URN"
---

# LabelGetByUrn

## Usage

```text
ionosctl label get-by-urn [flags]
```

## Description

Use this command to get information about a specified Label using its URN. A URN is used for uniqueness of a Label and composed using `urn:label:<resource_type>:<resource_uuid>:<key>`.

Required values to run command:

* Label URN

## Options

```text
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [URN Key Value ResourceType ResourceId] (default [URN,Key,Value,ResourceType,ResourceId])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --label-urn string   URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>] (required)
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout in seconds for polling the request (default 60)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait               Polls the request continuously until the operation is completed 
```

## Examples

```text
ionosctl label get-by-urn --label-urn "urn:label:server:SERVER_ID:test"
```

