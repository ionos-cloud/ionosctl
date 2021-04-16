---
description: Update a User
---

# Update

## Usage

```text
ionosctl user update [flags]
```

## Description

Use this command to update details about a specific User including their privileges.

Note: The password attribute is immutable. It is not allowed in update requests. It is recommended that the new User log into the DCD and change their password.

You can wait for the action to be executed using `--wait` option.

Required values to run command:

* User First Name
* User Last Name
* User Email
* User Id

## Options

```text
  -u, --api-url string           Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings             Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,Administrator,ForceSecAuth,SecAuthActive,S3CanonicalUserId,Active])
  -c, --config string            Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help                     help for update
      --ignore-stdin             Force command to execute without user input
  -o, --output string            Desired output format [text|json] (default "text")
  -q, --quiet                    Quiet output
      --timeout int              Timeout option for User to be updated [seconds] (default 60)
      --user-administrator       Assigns the User to have administrative rights
      --user-email string        The email for the User [Required flag]
      --user-first-name string   The firstname for the User [Required flag]
      --user-id string           The unique User Id [Required flag]
      --user-last-name string    The lastname for the User [Required flag]
      --user-secure-auth         Indicates if secure (two-factor) authentication should be forced for the User
      --wait                     Wait for User attributes to be updated
```

