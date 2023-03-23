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
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
      --label-urn string   URN for the Label [urn:label:<resource_type>:<resource_uuid>:<key>] (required)
      --no-headers         When using text output, don't print headers
```

## Examples

```text
ionosctl label get-by-urn --label-urn "urn:label:server:SERVER_ID:test"
```

