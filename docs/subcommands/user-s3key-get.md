---
description: Get a User S3Key
---

# UserS3keyGet

## Usage

```text
ionosctl user s3key get [flags]
```

## Description

Use this command to get information about a specified S3Key from a specified User.

Required values to run command:

* User Id
* S3Key Id

## Options

```text
  -u, --api-url string                             Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string                              Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                                      Force command to execute without user input
  -F, --format --cols="S3KeyId,Active,SecretKey"   Columns to be printed in the standard output. You can also print SecretKey, using --cols="S3KeyId,Active,SecretKey" (default [S3KeyId,Active])
  -h, --help                                       help for get
  -o, --output string                              Desired output format [text|json] (default "text")
  -q, --quiet                                      Quiet output
      --s3key-id string                            The unique User S3Key Id (required)
      --user-id string                             The unique User Id (required)
```

## Examples

```text
ionosctl user s3key get --user-id 013188d4-af9a-4207-b495-de36cb2dc344 --s3key-id 00a29d110b48daa3a18b 
S3KeyId                Active
00a29d110b48daa3a18b   false
```

