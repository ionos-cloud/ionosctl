---
description: "Create an Provider"
---

# CertmanagerProviderCreate

## Usage

```text
ionosctl certmanager provider create [flags]
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

For `create` command:

```text
[post c]
```

## Description

Create an Provider

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Email Server KeyId KeySecret State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --email string        The email address of the certificate requester
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --key-id string       The key ID of the external account binding
      --key-secret string   The key secret of the external account binding
      --limit int           Pagination limit: Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
  -n, --name string         The name of the certificate Provider
      --no-headers          Don't print table headers when table output is used
      --offset int          Pagination offset: Number of items to skip before starting to collect the results
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
      --server string       The URL of the certificate Provider
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certmanager provider create --name NAME --email EMAIL --server SERVER
```

