---
description: List available CPU Architecture from a Location
---

# LocationCpuList

## Usage

```text
ionosctl location cpu list [flags]
```

## Aliases

For `location` command:

```text
[loc]
```

For `list` command:

```text
[l ls]
```

## Description

Use this command to get information about available CPU Architectures from a specific Location.

Required values to run command:

* Location Id

## Options

```text
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10. (default 1)
      --location-id string   The unique Location Id (required)
  -M, --max-results int32    The maximum number of elements to return
      --no-headers           When using text output, don't print headers
```

## Examples

```text
ionosctl location cpu list --location-id LOCATION_ID
```

