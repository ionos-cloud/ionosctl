---
description: Blacklists the User, disabling them
---

# Delete

## Usage

```text
ionosctl user delete [flags]
```

## Description

This command blacklists the User, disabling them. The User is not completely purged, therefore if you anticipate needing to create a User with the same name in the future, we suggest renaming the User before you delete it.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,Administrator,ForceSecAuth,SecAuthActive,S3CanonicalUserId,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for delete
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
      --user-id string   The unique User Id [Required flag]
```

## Examples

```text
ionosctl user delete --user-id 2470f439-1d73-42f8-90a9-f78cf2776c74 --ignore-stdin 
RequestId: a2f6e7fa-6030-4267-950e-1e0886316475
Status: Command user delete has been successfully executed
```

