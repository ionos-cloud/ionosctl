---
description: List User S3Keys
---

# UserS3keyList

## Usage

```text
ionosctl user s3key list [flags]
```

## Description

Use this command to get a list of S3Keys of a specified User.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string                           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols --cols="S3KeyId,Active,SecretKey"   Columns to be printed in the standard output. You can also print SecretKey, using --cols="S3KeyId,Active,SecretKey" (default [S3KeyId,Active])
  -c, --config string                            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                                    Force command to execute without user input
  -h, --help                                     help for list
  -o, --output string                            Desired output format [text|json] (default "text")
  -q, --quiet                                    Quiet output
      --user-id string                           The unique User Id (required)
```

## Examples

```text
ionosctl user s3key list --user-id 013188d4-af9a-4207-b495-de36cb2dc344 
S3KeyId                Active
00a29d110b48daa3a18b   false
```

