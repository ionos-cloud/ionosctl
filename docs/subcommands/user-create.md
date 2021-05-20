---
description: Create a User under a particular contract
---

# UserCreate

## Usage

```text
ionosctl user create [flags]
```

## Description

Use this command to create a User under a particular contract. You need to specify the firstname, lastname, email and password for the new User.

Note: The password set here cannot be updated through the API currently. It is recommended that a new User log into the DCD and change their password.

Required values to run a command:

* User First Name
* User Last Name
* User Email
* User Password

## Options

```text
      --administrator       Assigns the User to have administrative rights
  -u, --api-url string      Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings        Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,S3CanonicalUserId,Administrator,ForceSecAuth,SecAuthActive,Active])
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --email string        The email for the User (required)
      --first-name string   The firstname for the User (required)
      --force               Force command to execute without user input
      --force-secure-auth   Indicates if secure (two-factor) authentication should be forced for the User
  -h, --help                help for create
      --last-name string    The lastname for the User (required)
  -o, --output string       Desired output format [text|json] (default "text")
      --password string     The password for the User (must be at least 5 characters long) (required)
  -q, --quiet               Quiet output
```

## Examples

```text
ionosctl user create --user-first-name test1 --user-last-name test1 --user-email testrandom16@gmail.com --user-password test123
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId   Active
99499053-059e-4ee6-b56f-66b0df93262d   test1       test1      testrandom16@ionos.com   false           false          false                               true
RequestId: ca349e08-5820-41ba-8252-ee4c8dd2ccdb
Status: Command user create has been successfully executed
```

