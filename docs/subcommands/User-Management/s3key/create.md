---
description: "Create a S3Key for a User"
---

# UserS3keyCreate

## Usage

```text
ionosctl user s3key create [flags]
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

For `create` command:

```text
[c]
```

## Description

Use this command to create a S3Key for a particular User.

Note: A maximum of five S3 keys may be created for any given user.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string     Override default host URL. If set, this will be preferred over the config file override. If unset, the default will only be used as a fallback (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [S3KeyId Active SecretKey] (default [S3KeyId,Active,SecretKey])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -t, --timeout int        Timeout option for Request for User S3Key creation [seconds] (default 60)
      --user-id string     The unique User Id (required)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait-for-request   Wait for the Request for User S3Key creation to be executed
```

## Examples

```text
ionosctl user s3key create --user-id USER_ID
```

