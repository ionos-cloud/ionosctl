---
description: "Delete a private image"
---

# ImageDelete

## Usage

```text
ionosctl compute image delete [flags]
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

Delete one of your PRIVATE images, or all of them with --all. PUBLIC (IONOS-provided) images cannot be deleted and are always skipped.

IMPORTANT: for an image you uploaded via FTP, this API call only sets the image size to 0B — it does NOT remove the underlying file from the FTP server. To fully remove an FTP-uploaded image you must contact IONOS support.

Required values to run command:

* Image Id

## Options

```text
  -a, --all               Delete every PRIVATE (non-public) image on the contract. Public images are skipped because the API forbids deleting them
  -u, --api-url string    Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [ImageId Name ImageAliases Location LicenceType ImageType CloudInit CreatedDate Size Description Public CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType RequiredFeatures]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int         Level of detail for response objects (default 1)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -i, --image-id string   The unique Image Id (required)
      --limit int         Maximum number of items to return per request (default 50)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl compute image delete --image-id IMAGE_ID
```

