---
description: "Create a new Token"
---

# TokenGenerate

## Usage

```text
ionosctl token generate [flags]
```

## Aliases

For `generate` command:

```text
[create]
```

## Description

Use this command to generate a new Token. Only the JSON Web Token, associated with user credentials, will be displayed.

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
      --cols strings     Set of columns to be printed on output 
                         Available columns: [TokenId CreatedDate ExpirationDate Href]
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --contract int     Users with multiple contracts can provide the contract number, for which the token is generated
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl token generate
```

