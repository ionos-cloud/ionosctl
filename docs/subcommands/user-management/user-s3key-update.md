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
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [S3KeyId Active] (default [S3KeyId,Active])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -D, --depth int          Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
      --s3key-active       Enable or disable an User S3Key. E.g.: --s3key-active=true, --s3key-active=false
  -i, --s3key-id string    The unique User S3Key Id (required)
  -t, --timeout int        Timeout option for Request for User S3Key update [seconds] (default 60)
      --user-id string     The unique User Id (required)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait-for-request   Wait for the Request for User S3Key update to be executed
```

## Examples

```text
ionosctl user s3key update --user-id USER_ID --s3key-id S3KEY_ID --s3key-active=false
```

