---
description: Get a Location
---

# LocationGet

## Usage

```text
ionosctl location get [flags]
```

## Aliases

For `location` command:

```text
[loc]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specific Location from a Region.

Required values to run command:

* Location Id

## Options

```text
  -D, --depth int32          Controls the detail depth of the response objects. Max depth is 10.
  -i, --location-id string   The unique Location Id (required)
      --no-headers           When using text output, don't print headers
```

## Examples

```text
ionosctl location get --location-id LOCATION_ID
```

