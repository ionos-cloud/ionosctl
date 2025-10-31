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
  -u, --api-url string   Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
      --limit int        pagination limit: Maximum number of items to return per request (default 50)
      --no-headers       Don't print table headers when table output is used
      --offset int       pagination offset: Number of items to skip before starting to collect the results
  -o, --output string    Desired output format [text|json|api-json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose count    Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl dbaas mongo api-versions
```

