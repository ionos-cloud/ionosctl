---
description: Get a User S3Key
---

# UserS3KeyGet

## Usage

```text
ionosctl user s3key get [flags]
```

## Aliases

For `user` command:

```text
[u]
```

For `s3key` command:

```text
[k s3k]
```

For `get` command:

```text
[g]
```

## Description

Use this command to get information about a specified S3Key from a specified User.

Required values to run command:

* User Id
* S3Key Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [S3KeyId Active] (default [S3KeyId,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              help for get
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -i, --s3key-id string   The unique User S3Key Id (required)
      --user-id string    The unique User Id (required)
```

## Examples

```text
ionosctl user s3key get --user-id USER_ID --s3key-id S3KEY_ID
```

