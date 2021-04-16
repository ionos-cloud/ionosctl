---
description: Get a User
---

# Get

## Usage

```text
ionosctl user get [flags]
```

## Description

Use this command to retrieve details about a specific User.

Required values to run command:

* User Id

## Options

```text
  -u, --api-url string   Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cols strings     Columns to be printed in the standard output (default [UserId,Firstname,Lastname,Email,Administrator,ForceSecAuth,SecAuthActive,S3CanonicalUserId,Active])
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -h, --help             help for get
      --ignore-stdin     Force command to execute without user input
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
      --user-id string   The unique User Id [Required flag]
```

