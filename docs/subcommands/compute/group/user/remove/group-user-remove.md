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

For `remove` command:

```text
[r]
```

## Description

Use this command to remove a User from a Group.

Required values to run command:

* Group Id
* User Id

## Options

```text
  -a, --all               Remove all Users from a group.
      --cols strings      Set of columns to be printed on output 
                          Available columns: [UserId Firstname Lastname Email S3CanonicalUserId Administrator ForceSecAuth SecAuthActive Active] (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
      --group-id string   The unique Group Id (required)
  -i, --user-id string    The unique User Id (required)
```

## Examples

```text
ionosctl group user remove --group-id GROUP_ID --user-id USER_ID
```

