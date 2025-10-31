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
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32       Controls the detail depth of the response objects. Max depth is 10.
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -i, --image-id string   The unique Image Id (required)
      --limit int         pagination limit: Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        pagination offset: Number of items to skip before starting to collect the results
  -o, --output string     Desired output format [text|json|api-json] (default "text")
  -q, --quiet             Quiet output
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl image get --image-id IMAGE_ID
```

