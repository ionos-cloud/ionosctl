---
description: "Delete an image"
---

# ImageDelete

## Usage

```text
ionosctl image delete [flags]
```

## Aliases

For `image` command:

```text
[img]
```

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specified Image.

Required values to run command:

* Image Id

## Options

```text
  -a, --all                Delete all non-public images
  -u, --api-url string     Override default host URL. Preferred over the config file override 'cloud'|'compute' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -i, --image-id string    The unique Image Id (required)
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout option for Request for Image update [seconds] (default 60)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait-for-request   Wait for the Request for Image update to be executed
```

## Examples

```text
ionosctl image delete --image-id IMAGE_ID
```

