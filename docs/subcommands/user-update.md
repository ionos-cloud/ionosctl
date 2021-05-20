---
description: Update a User
---

# UserUpdate

## Usage

```text
ionosctl user update [flags]
```

## Description

Use this command to update details about a specific User including their privileges.

Note: The password attribute is immutable. It is not allowed in update requests. It is recommended that the new User log into the DCD and change their password.

Required values to run command:

* User Id

## Options

```text
      --administrator       Assigns the User to have administrative rights
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --email string        The email for the User
      --first-name string   The firstname for the User
      --force               Force command to execute without user input
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User
  -h, --help                help for update
      --last-name string    The lastname for the User
  -o, --output string       Desired output format [text|json] (default "text")
  -q, --quiet               Quiet output
      --user-id string      The unique User Id (required)
```

## Examples

```text
ionosctl user update --user-id 2470f439-1d73-42f8-90a9-f78cf2776c74 --user-administrator=true
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
2470f439-1d73-42f8-90a9-f78cf2776c74   test1       test1      testrandom12@ionos.com   true            false          false           a74101e7c1948450432d5b6512f2712c   true
RequestId: 439f79fc-5bfc-43da-92f3-0d804ebb28ac
Status: Command user update has been successfully executed
```

