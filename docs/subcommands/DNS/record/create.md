---
description: "Create a DNS record"
---

# DnsRecordCreate

## Usage

```text
ionosctl dns record create [flags]
```

## Aliases

For `record` command:

```text
[r]
```

For `create` command:

```text
[c post]
```

## Description

Create a DNS record inside a zone.

Three things define a record: --type, --name and --content. --name is the host under the zone ('www' for www.example.com, the zone name itself or an empty string '' for the apex, '*' for a wildcard). --content is the record's data and its meaning depends on --type:

  A       IPv4 address            e.g. 1.2.3.4
  AAAA    IPv6 address            e.g. 2001:db8::1
  CNAME   target hostname         e.g. www.example.com
  ALIAS   target hostname (apex)  e.g. example.com
  MX      mail server hostname    e.g. mail.example.com   (set --priority)
  NS      name server hostname    e.g. ns1.example.com
  TXT     free text               e.g. "v=spf1 -all"
  SRV     "weight port target"    e.g. "5 5060 sip.example.com"  (set --priority)
  CAA     flags tag "value"       e.g. 0 issue "letsencrypt.org"

--priority is required for MX, SRV and URI and ignored otherwise. --ttl sets the cache lifetime in seconds (60-604800, default 3600). Records are --enabled by default.

## Options

```text
  -u, --api-url string    Override default host URL. If contains placeholder, location will be embedded. Preferred over the config file override 'dns' and env var 'IONOS_API_URL' (default "https://dns.%s.ionos.com")
      --cols strings      Set of columns to be printed on output 
                          Available columns: [Id Name Content Type Enabled FQDN ZoneId ZoneName State]
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --content string    Record data, interpreted per --type: an A record takes an IPv4 (1.2.3.4), AAAA an IPv6, CNAME/MX/NS a hostname, TXT free text. See this command's --help for the full per-type table (required)
  -D, --depth int         Level of detail for response objects (default 1)
      --enabled           Whether the record answers lookups. true = live; false = kept but not served (default true) (default true)
  -F, --filters strings   Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
      --limit int         Maximum number of items to return per request (default 50)
  -l, --location string   Location of the resource to operate on. When unset, list commands query all locations. Can be one of: de/fra. A facility inside one of these metro regions (e.g. de/fra/1) is also accepted and served by its metro region's endpoint (default "de/fra")
  -n, --name string       Host under the zone this record answers for, e.g. 'www'. For the apex, use the zone name itself (an empty name is also accepted). Use '*' for a wildcard matching non-existent names (some shells need it escaped as '\*') (required)
      --no-headers        Don't print table headers when table output is used
      --offset int        Number of items to skip before starting to collect the results
      --order-by string   Property to order the results by
  -o, --output string     Desired output format [text|json|api-json] (default "text")
      --priority int32    Preference value 0-65535, lower wins. Required for MX, SRV and URI records; ignored for all other types
      --query string      JMESPath query string to filter the output
  -q, --quiet             Quiet output
  -t, --timeout int       Timeout in seconds for --wait and other wait operations (default 600)
      --ttl int32         How long (seconds) resolvers may cache this record before re-querying; 60-604800 (default 3600 = 1h) (default 3600)
      --type string       Record type; decides how --content is interpreted (A=IPv4, AAAA=IPv6, CNAME/MX/NS=hostname, TXT=text, …). Can be one of: A, AAAA, CNAME, ALIAS, MX, NS, SRV, TXT, CAA, SSHFP, TLSA, SMIMEA, DS, HTTPS, SVCB, OPENPGPKEY, CERT, URI, RP, LOC (required) (default "AAAA")
  -v, --verbose count     Increase verbosity level [-v, -vv, -vvv]
  -w, --wait              Wait for the resource to reach AVAILABLE state after the command completes. No-op for list commands
  -z, --zone string       The ID or name of the DNS zone (required)
```

## Examples

```text
ionosctl dns record create --zone example.com --type A --name www --content 1.2.3.4
ionosctl dns record create --zone example.com --type MX --name example.com --content mail.example.com --priority 10 --ttl 300
```

