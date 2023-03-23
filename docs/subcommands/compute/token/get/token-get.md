---
description: Get a specified Token
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
      --contract int      Users with multiple contracts must provide the contract number, for which the token information is displayed
      --no-headers        When using text output, don't print headers
  -i, --token-id string   The unique Key ID of a Token (required)
```

## Examples

```text
ionosctl token get --token-id TOKEN_ID
```

