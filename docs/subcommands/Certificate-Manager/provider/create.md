---
description: "Register an ACME certificate provider"
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

Register an ACME certificate authority (for example Let's Encrypt) that auto-certificates can then use to be issued and auto-renewed.

Required: --name (a label), --email (the ACME account contact address), and --server (the CA's ACME directory URL, e.g. https://acme-v02.api.letsencrypt.org/directory).

External account binding (EAB) is optional and links your IONOS account to a pre-registered account at the CA. If you use it, --key-id and --key-secret must be supplied together.

## Options

```text
  -u, --api-url string      Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'cert' and env var 'IONOS_API_URL' (default "https://certificate-manager.%s.ionos.com")
      --cols strings        Set of columns to be printed on output 
                            Available columns: [Id Name Email Server KeyId KeySecret State]
  -c, --config string       Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int           Level of detail for response objects (default 1)
      --email string        Contact email registered with the ACME account; the CA uses it for expiry and policy notices. Required
  -F, --filters strings     Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force               Force command to execute without user input
  -h, --help                Print usage
      --key-id string       External account binding (EAB) key ID issued by the CA. Optional; must be given together with --key-secret
      --key-secret string   External account binding (EAB) HMAC key secret issued by the CA. Optional; must be given together with --key-id. Write-only
      --limit int           Maximum number of items to return per request (default 50)
  -l, --location string     Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
  -n, --name string         A label to identify this provider (e.g. "Let's Encrypt"). Required
      --no-headers          Don't print table headers when table output is used
      --offset int          Number of items to skip before starting to collect the results
      --order-by string     Property to order the results by
  -o, --output string       Desired output format [text|json|api-json] (default "text")
      --query string        JMESPath query string to filter the output
  -q, --quiet               Quiet output
      --server string       The CA's ACME directory URL, e.g. https://acme-v02.api.letsencrypt.org/directory. Required
  -t, --timeout int         Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count       Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
# Register Let's Encrypt (production directory)
ionosctl certmanager provider create --name letsencrypt --email admin@example.com --server https://acme-v02.api.letsencrypt.org/directory

# Register a CA that requires external account binding (EAB)
ionosctl certmanager provider create --name my-ca --email admin@example.com --server https://acme.my-ca.example/directory --key-id my-eab-key-id --key-secret my-eab-secret
```

