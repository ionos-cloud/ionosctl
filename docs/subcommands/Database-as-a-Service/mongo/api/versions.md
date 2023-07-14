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
[mongodb mdb m]
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
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl dbaas mongo api-versions
```

