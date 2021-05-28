---
description: Get a specified Image
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
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              help for get
  -i, --image-id string   The unique Image Id (required)
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
```

## Examples

```text
ionosctl image get --image-id IMAGE_ID
```

