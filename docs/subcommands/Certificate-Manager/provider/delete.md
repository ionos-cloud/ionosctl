---
description: "Delete an Provider"
---

# CertmanagerProviderDelete

## Usage

```text
ionosctl certmanager provider delete [flags]
```

## Aliases

For `certmanager` command:

```text
[cert certs certificate-manager certificates certificate]
```

For `provider` command:

```text
[providers]
```

For `delete` command:

```text
[del d]
```

## Description

Delete an Provider

## Options

```text
  -a, --all                  Delete all Providers. Required or -g
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name Email Server KeyId KeySecret State]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --limit int            pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers           Don't print table headers when table output is used
      --offset int           pagination offset: Number of items to skip before starting to collect the results
  -o, --output string        Desired output format [text|json|api-json] (default "text")
  -i, --provider-id string   Provide the specified Provider (required)
  -q, --quiet                Quiet output
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certmanager provider delete --provider-id ID
```

