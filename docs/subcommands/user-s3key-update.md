---
description: Update a S3Key
---

# UserS3keyUpdate

## Usage

```text
ionosctl user s3key update [flags]
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

For `update` command:

```text
[u up]
```

## Description

Use this command to update a specified S3Key from a particular User. This operation allows you to enable or disable a specific S3Key.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* User Id
* S3Key Id
* S3Key Active

## Options

```text
  -u, --api-url string     Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [S3KeyId Active] (default [S3KeyId,Active])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               help for update
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --s3key-active       Enable or disable an User S3Key
  -i, --s3key-id string    The unique User S3Key Id (required)
  -t, --timeout int        Timeout option for Request for User S3Key update [seconds] (default 60)
      --user-id string     The unique User Id (required)
  -w, --wait-for-request   Wait for the Request for User S3Key update to be executed
```

## Examples

```text
ionosctl user s3key update --user-id USER_ID --s3key-id S3KEY_ID --s3key-active=false
```

