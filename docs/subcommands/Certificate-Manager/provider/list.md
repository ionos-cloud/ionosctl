---
description: "Retrieve Provider list"
---

# CertmanagerProviderList

## Usage

```text
ionosctl certmanager provider list [flags]
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

For `list` command:

```text
[ls]
```

## Description

Retrieve Provider list

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Email Server KeyId KeySecret State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
  -l, --location string     Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --max-results int32   Pagination limit
      --no-headers          Don't print table headers when table output is used
      --offset int32        Pagination offset
  -o, --output string       Desired output format [text|json|api-json] (default "text")
  -q, --quiet               Quiet output
  -v, --verbose             Print step-by-step process when running command
```

## Examples

```text
ionosctl certmanager provider list
```

