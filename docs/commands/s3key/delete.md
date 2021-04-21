---
description: Delete a S3Key
---

# Delete

## Usage

```text
ionosctl s3key delete [flags]
```

## Description

Use this command to delete a specific S3Key of an User.

Required values to run command:

* User Id
* S3Key Id

## Options

```text
  -u, --api-url string                           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols --cols="S3KeyId,Active,SecretKey"   Columns to be printed in the standard output. You can also print SecretKey, using --cols="S3KeyId,Active,SecretKey" (default [S3KeyId,Active])
  -c, --config string                            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                                     help for delete
      --ignore-stdin                             Force command to execute without user input
  -o, --output string                            Desired output format [text|json] (default "text")
  -q, --quiet                                    Quiet output
      --s3key-id string                          The unique User S3Key Id [Required flag]
      --timeout int                              Timeout option for a S3Key to be deleted [seconds] (default 60)
      --user-id string                           The unique User Id [Required flag]
      --wait                                     Wait for S3Key to be deleted
```

## Examples

```text
ionosctl s3key delete --user-id 013188d4-af9a-4207-b495-de36cb2dc344 --s3key-id 75f4319cbf3f6d538da7 --wait 
Warning: Are you sure you want to delete S3Key (y/N) ? 
y
RequestId: 1529f8b7-08bb-4321-a996-08865660dee8
Status: Command s3key delete and request have been successfully executed
```

