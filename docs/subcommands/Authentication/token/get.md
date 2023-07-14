---
description: "Get a specified Token"
---

# TokenGet

## Usage

```text
ionosctl token get [flags]
```

## Aliases

For `get` command:

```text
[g]
```

## Description

Use this command to retrieve details about a Token by using its ID.

Required values to run command:

* Token Id

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --contract int      Users with multiple contracts must provide the contract number, for which the token information is displayed
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --no-headers        When using text output, don't print headers
  -o, --output string     Desired output format [text|json] (default "text")
  -q, --quiet             Quiet output
  -t, --token string      The contents of a Token (required)
  -i, --token-id string   The unique Key ID of a Token (required)
  -v, --verbose           Print step-by-step process when running command
```

## Examples

```text
ionosctl token get --token-id TOKEN_ID
ionosctl token get --token TOKEN
```

