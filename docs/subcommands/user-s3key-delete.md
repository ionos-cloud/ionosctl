---
description: Delete a S3Key
---

# UserS3keyDelete

## Usage

```text
ionosctl user s3key delete [flags]
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

For `delete` command:

```text
[d]
```

## Description

Use this command to delete a specific S3Key of an User.

Required values to run command:

* User Id
* S3Key Id

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [S3KeyId Active] (default [S3KeyId,Active])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               help for delete
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -i, --s3key-id string    The unique User S3Key Id (required)
  -t, --timeout int        Timeout option for Request for User S3Key deletion [seconds] (default 60)
      --user-id string     The unique User Id (required)
  -w, --wait-for-request   Wait for Request for User S3Key deletion to be executed
```

## Examples

```text
ionosctl user s3key delete --user-id USER_ID --s3key-id S3KEY_ID --force
```

