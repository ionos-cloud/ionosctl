---
description: Get a Label using URN
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
                           Available columns: [Key Value] (default [Key,Value])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --label-urn string   URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>] (required)
      --no-headers         When using text output, don't print headers
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -v, --verbose            Print step-by-step process when running command
```

## Examples

```text
ionosctl label get-by-urn --label-urn "urn:label:server:SERVER_ID:test"
```

