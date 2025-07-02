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

Use this command to parse a Token and find out Token ID, User ID, Contract Number, Role.
If you want to view the privileges associated with the token, you must set the --privileges flag. When this flag is set, the command will output a list of privileges instead of the default output.

Required values to run:

* Token

## Options

```text
  -u, --api-url string   Override default host URL. Preferred over the config file override 'auth' and env var 'IONOS_API_URL' (default "https://api.ionos.com/auth/v1")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
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

