---
description: "Update Certificate name"
---

# CertmanagerCertificateUpdate

## Usage

```text
ionosctl certmanager certificate update [flags]
```

## Aliases

For `certmanager` command:

```text
[cert certs certificate-manager certificates certificate]
```

For `certificate` command:

```text
[cert certificates certs]
```

For `update` command:

```text
[u]
```

## Description

Use this change a certificate's name.

## Options

```text
  -u, --api-url string            Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
  -i, --certificate-id string     Provide the specified Certificate (required)
  -n, --certificate-name string   Provide new certificate name (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [CertId DisplayName Expired NotAfter NotBefore]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                 Level of detail for response objects (default 1)
      --filters strings           Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --limit int                 Maximum number of items to return per request (default 50)
  -l, --location string           Location of the resource to operate on. Can be one of: de/fra (default "de/fra")
      --no-headers                Don't print table headers when table output is used
      --offset int                Number of items to skip before starting to collect the results
      --order-by string           Property to order the results by
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --query string              JMESPath query string to filter the output
  -q, --quiet                     Quiet output
  -v, --verbose count             Increase verbosity level [-v, -vv, -vvv]
```

## Examples

```text
ionosctl certificate-manager update --certificate-id 47c5d9cc-b613-4b76-b0cc-dc531787a422
```

