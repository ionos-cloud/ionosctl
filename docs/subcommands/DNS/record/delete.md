---
description: Delete a record
---

# DnsRecordDelete

## Usage

```text
ionosctl dns record delete [flags]
```

## Aliases

For `record` command:

```text
[r]
```

For `delete` command:

```text
[del d]
```

## Description

Delete a record

## Options

```text
  -a, --all                Delete all records. Required or --zone and --record-id
  -u, --api-url string     Override default host url (default "https://api.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [Id Name Content Type Enabled FQDN State ZoneId ZoneName]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Skip yes/no confirmation
  -h, --help               Print usage
  -n, --name string        If --all is set, filter --all deletion by record name
      --no-headers         When using text output, don't print headers
  -o, --output string      Desired output format [text|json] (default "text")
  -q, --quiet              Quiet output
  -i, --record-id string   The ID or name of the DNS record. Required together with --zone or -a
  -v, --verbose            Print step-by-step process when running command
  -z, --zone string        The zone of the target record. If --all is set, filter --all deletion by limiting to records within this zone
```

## Examples

```text
ionosctl dns record delete --zone-id ZONE --record-id RECORD
ionosctl dns record delete --all [--name PARTIAL_NAME] [--zone ZONE]
ionosctl dns record delete --name PARTIAL_NAME [--zone ZONE]
```

