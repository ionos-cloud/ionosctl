---
description: "Get a specified Image"
---

# ImageGet

## Usage

```text
ionosctl image get [flags]
```

## Aliases

For `image` command:

```text
[img]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified Image.

Required values to run command:

* Image Id

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -i, --image-id string   The unique Image Id (required)
      --no-headers        When using text output, don't print headers
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl image get --image-id IMAGE_ID
```

