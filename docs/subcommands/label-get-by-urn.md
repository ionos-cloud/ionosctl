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
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [Key Value] (default [Key,Value])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id
  -f, --force                  Force command to execute without user input
  -h, --help                   help for get-by-urn
      --ipblock-id string      The unique IpBlock Id
      --label-urn string       URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>] (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --resource-type string   Resource Type
      --server-id string       The unique Server Id
      --snapshot-id string     The unique Snapshot Id
      --volume-id string       The unique Volume Id
```

## Examples

```text
ionosctl label get-by-urn --label-urn "urn:label:server:SERVER_ID:test"
```

