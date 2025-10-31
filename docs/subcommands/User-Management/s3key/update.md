---
description: "Update a S3Key"
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
  -u, --api-url string     Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [S3KeyId Active SecretKey] (default [S3KeyId,Active,SecretKey])
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int32        Controls the detail depth of the response objects. Max depth is 10.
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
      --limit int          pagination limit: Maximum number of items to return per request (default 50)
      --no-headers         Don't print table headers when table output is used
      --offset int         pagination offset: Number of items to skip before starting to collect the results
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
      --s3key-active       Enable or disable an User S3Key. E.g.: --s3key-active=true, --s3key-active=false
  -i, --s3key-id string    The unique User S3Key Id (required)
  -t, --timeout int        Timeout option for Request for User S3Key update [seconds] (default 60)
      --user-id string     The unique User Id (required)
  -v, --verbose count      Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request   Wait for the Request for User S3Key update to be executed
```

## Examples

```text
ionosctl user s3key update --user-id USER_ID --s3key-id S3KEY_ID --s3key-active=false
```

