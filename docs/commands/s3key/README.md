---
description: S3Key Operations
---

# S3key

## Usage

```text
ionosctl s3key [command]
```

## Description

The sub-commands of `ionosctl s3key` allow you to see information, to list, get, create, update, delete Users S3Keys. To view details about Users, check the `ionosctl user` commands.

## Options

```text
  -u, --api-url string                           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols --cols="S3KeyId,Active,SecretKey"   Columns to be printed in the standard output. You can also print SecretKey, using --cols="S3KeyId,Active,SecretKey" (default [S3KeyId,Active])
  -c, --config string                            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --force                                    Force command to execute without user input
  -h, --help                                     help for s3key
  -o, --output string                            Desired output format [text|json] (default "text")
  -q, --quiet                                    Quiet output
      --user-id string                           The unique User Id [Required flag]
```

## Related commands

| Command | Description |
| :--- | :--- |
| [ionosctl s3key create](create.md) | Create a S3Key for a User |
| [ionosctl s3key delete](delete.md) | Delete a S3Key |
| [ionosctl s3key get](get.md) | Get a User S3Key |
| [ionosctl s3key list](list.md) | List User S3Keys |
| [ionosctl s3key update](update.md) | Update a S3Key |

