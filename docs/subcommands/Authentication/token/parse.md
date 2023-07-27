---
description: "Parse the contents of a Token"
---

# TokenParse

## Usage

```text
ionosctl token parse [flags]
```

## Aliases

For `parse` command:

```text
[p]
```

## Description

Use this command to parse a Token and find out Token ID, User ID, Contract Number, Role and Privileges (separate).

Required values to run: 

* Token

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -p, --privileges       Use to see the privileges that the user using this Token benefits from
  -q, --quiet            Quiet output
  -t, --token string     The contents of a Token (required)
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl token parse --token TOKEN

ionosctl token parse --privileges --token TOKEN
```
