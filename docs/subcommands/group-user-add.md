---
description: Add User to a Group
---

# GroupUserAdd

## Usage

```text
ionosctl group user add [flags]
```

## Aliases

For `group` command:
```text
[g]
```

For `user` command:
```text
[u]
```

## Description

Use this command to add an existing User to a specific Group.

Required values to run command:

* Group Id
* User Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              help for add
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -i, --user-id string    The unique User Id (required)
```

## Examples

```text
ionosctl group user add --group-id 45ba215b-6897-40b6-879c-cbadb527cefd --user-id 62599641-aa2d-4ecc-bdc4-118f5f39f23d 
UserId                                 Firstname   Lastname   Email                    S3CanonicalUserId                  Administrator   ForceSecAuth   SecAuthActive   Active
62599641-aa2d-4ecc-bdc4-118f5f39f23d   test        test       testrandom53@gmail.com   f670112b3e74038b51db78d5836d7854   false           false          false           true
RequestId: 296f4d86-629c-44f4-bacc-0fefb2356029
Status: Command user add has been successfully executed
```

