---
description: Blacklists the User, disabling them
---

# UserDelete

## Usage

```text
ionosctl user delete [flags]
```

## Aliases

For `user` command:
```text
[u]
```

## Description

This command blacklists the User, disabling them. The User is not completely purged, therefore if you anticipate needing to create a User with the same name in the future, we suggest renaming the User before you delete it.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -F, --format strings   Collection of fields to be printed on output (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -h, --help             help for delete
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
      --user-id string   The unique User Id (required)
```

## Examples

```text
ionosctl user delete --user-id 2470f439-1d73-42f8-90a9-f78cf2776c74 --force 
RequestId: a2f6e7fa-6030-4267-950e-1e0886316475
Status: Command user delete has been successfully executed
```

