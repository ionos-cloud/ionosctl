---
description: Remove User from a Group
---

# GroupUserRemove

## Usage

```text
ionosctl group user remove [flags]
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

Use this command to remove a User from a Group.

Required values to run command:

* Group Id
* User Id

## Options

```text
  -u, --api-url string    Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -C, --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
      --group-id string   The unique Group Id (required)
  -h, --help              help for remove
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
      --user-id string    The unique User Id (required)
```

## Examples

```text
ionosctl group user remove --group-id 45ba215b-6897-40b6-879c-cbadb527cefd --user-id 62599641-aa2d-4ecc-bdc4-118f5f39f23d 
Warning: Are you sure you want to remove user from group (y/N) ? 
y
RequestId: 07e1eb6a-2618-42dd-b614-6b34359a79b3
Status: Command user remove has been successfully executed
```

