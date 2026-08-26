---
description: "Delete an auto-certificate by ID, or all auto-certificates"
---

# CertmanagerAutocertificateDelete

## Usage

```text
ionosctl certmanager autocertificate delete [flags]
```

## Aliases

For `certmanager` command:

```text
[cert certs certificate-manager certificates certificate]
```

For `autocertificate` command:

```text
[autocert autocerts auto autocertificates]
```

For `delete` command:

```text
[del d]
```

## Description

Delete an auto-certificate. Pass --autocertificate-id to delete one, or --all to delete every auto-certificate in the account.

Deleting stops future auto-renewals for that certificate. Ensure no product still relies on it before removing it.

## Options

```text
  -a, --all                         Delete every auto-certificate in the account. Use instead of --autocertificate-id
  -u, --api-url string              Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
  -i, --autocertificate-id string   The ID (UUID) of the auto-certificate to delete. Required unless --all is set (required)
      --cols strings                Set of columns to be printed on output 
                                    Available columns: [Id Provider CommonName KeyAlgorithm Name AlternativeNames State]
  -c, --config string               Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int                   Level of detail for response objects (default 1)
  -F, --filters strings             Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                       Force command to execute without user input
  -h, --help                        Print usage
      --limit int                   Maximum number of items to return per request (default 50)
  -l, --location string             Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
      --no-headers                  Don't print table headers when table output is used
      --offset int                  Number of items to skip before starting to collect the results
      --order-by string             Property to order the results by
  -o, --output string               Desired output format [text|json|api-json] (default "text")
      --query string                JMESPath query string to filter the output
  -q, --quiet                       Quiet output
  -t, --timeout int                 Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count               Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                        Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl certmanager autocertificate delete --autocertificate-id ID
```

