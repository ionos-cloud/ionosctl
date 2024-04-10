---
description: "Get Mongo API swagger files"
---

# DbaasMongoApiVersions

## Usage

```text
ionosctl dbaas mongo api-versions [flags]
```

## Aliases

For `dbaas` command:

```text
[db]
```

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `api-versions` command:

```text
[versions api-version]
```

## Description

Get Mongo API swagger files

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --no-headers       Don't print table headers when table output is used
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo api-versions
```

