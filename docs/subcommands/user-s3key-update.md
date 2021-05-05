---
description: Update a S3Key
---

# UserS3keyUpdate

## Usage

```text
ionosctl user s3key update [flags]
```

## Description

Use this command to update a specified S3Key from a particular User. This operation allows you to enable or disable a specific S3Key.
You can wait for the action to be executed using `--wait` option.
Required values to run command:
* User Id
* S3Key Id
* S3Key Active

## Options

```text
  -u, --api-url string                           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols --cols="S3KeyId,Active,SecretKey"   Columns to be printed in the standard output. You can also print SecretKey, using --cols="S3KeyId,Active,SecretKey" (default [S3KeyId,Active])
  -c, --config string                            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                                    Force command to execute without user input
  -h, --help                                     help for update
  -o, --output string                            Desired output format [text|json] (default "text")
  -q, --quiet                                    Quiet output
      --s3key-active                             Enable or disable a S3Key
      --s3key-id string                          The unique User S3Key Id (required)
      --timeout int                              Timeout option for a S3Key to be updated [seconds] (default 60)
      --user-id string                           The unique User Id (required)
      --wait                                     Wait for S3Key to be updated
```

## Examples

```text
ionosctl user s3key update --user-id 013188d4-af9a-4207-b495-de36cb2dc344 --s3key-id 75f4319cbf3f6d538da7 --s3key-active=false
S3KeyId                Active
75f4319cbf3f6d538da7   false
RequestId: 4cda4b65-f58b-492a-bf45-6f1d8fb42928
Status: Command s3key update has been successfully executed
```

