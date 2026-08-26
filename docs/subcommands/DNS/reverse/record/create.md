---
description: "Create a reverse DNS (PTR) record"
---

# DnsReverseRecordCreate

## Usage

```text
ionosctl dns reverse-record create [flags]
```

## Aliases

For `reverse-record` command:

```text
[rr]
```

For `create` command:

```text
[c post]
```

## Description

Create a reverse DNS record so a PTR lookup on --ip returns --name.

--ip must be an IPv4 or IPv6 address owned by your contract (e.g. from a reserved IP block); --name is the hostname it should resolve back to. Commonly used to give a mail server matching forward and reverse DNS.

Wiki: https://docs.ionos.com/cloud/network-services/cloud-dns/api-how-tos/create-and-manage-reverse-dns

## Options

```text
  -u, --api-url string       Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings         Set of columns to be printed on output 
                             Available columns: [Id Name IP Description]
  -c, --config string        Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
  -D, --depth int            Level of detail for response objects (default 1)
      --description string   Free-text note for your own reference; not served in DNS
  -F, --filters strings      Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                Force command to execute without user input
  -h, --help                 Print usage
      --ip string            IPv4 or IPv6 address to create the reverse record for; must be owned by your contract (required)
      --limit int            Maximum number of items to return per request (default 50)
  -l, --location string      Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
  -n, --name string          Hostname the IP should resolve back to, e.g. mail.example.com (required)
      --no-headers           Don't print table headers when table output is used
      --offset int           Number of items to skip before starting to collect the results
      --order-by string      Property to order the results by
  -o, --output string        Desired output format [text|json|api-json] (default "text")
      --query string         JMESPath query string to filter the output
  -q, --quiet                Quiet output
  -t, --timeout int          Timeout in seconds for --wait and other wait operations (default 600)
  -v, --verbose count        Increase verbosity level [-v, -vv, -vvv]
  -w, --wait                 Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
```

## Examples

```text
ionosctl dns reverse-record create --name mail.example.com --ip 5.6.7.8
```

