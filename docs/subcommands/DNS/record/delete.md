---
description: "Delete a record"
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

To delete a specific record from a specific zone:
ionosctl dns r del --zone ZONE --record RECORD
Here, ZONE is the ID or name of the DNS zone from where you want to delete a record, and RECORD is the ID or full name of the DNS record you want to delete.

To delete all records, optionally filtering by partial name and zone:
ionosctl dns r delete --all [--record PARTIAL_NAME] [--zone ZONE]
Here, --all deletes all DNS records. You can also filter the records to delete by providing a PARTIAL_NAME that matches part of the name of the records you want to delete. Additionally, you can specify a ZONE to restrict the deletion to a specific DNS zone.

To delete a record by partial name, specifying the zone:
ionosctl dns r delete --record PARTIAL_NAME --zone ZONE
Here, PARTIAL_NAME is a part of the name of the DNS record you want to delete. If multiple records match the partial name, an error will be thrown: you will need to narrow down to a single record

## Options

```text
  -a, --all                Delete all records. You can optionally filter the deleted records using --zone (full name / ID) and --record (partial name)
  -u, --api-url string     Override default host URL (default "https://dns.de-fra.ionos.com")
      --cols strings       Set of columns to be printed on output 
                           Available columns: [Id Name Content Type Enabled FQDN ZoneId ZoneName State]
  -c, --config string      Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force              Force command to execute without user input
  -h, --help               Print usage
  -l, --location string    Location of the resource to operate on. Can be one of: de/fra
      --no-headers         Don't print table headers when table output is used
  -o, --output string      Desired output format [text|json|api-json] (default "text")
  -q, --quiet              Quiet output
  -r, --record string      The ID, or full name of the DNS record. Required together with --zone. Can also provide partial names, but must narrow down to a single record result if not using --all. If using it, will however delete all records that match.
  -t, --timeout duration   Timeout for waiting for resource to reach desired state (default 1m0s)
  -v, --verbose            Print step-by-step process when running command
  -w, --wait               Polls the request continuously until the operation is completed
  -z, --zone string        The full name or ID of the zone of the containing the target record. If --all is set this is applied as a filter - limiting to records within this zone
```

## Examples

```text
ionosctl dns r del --zone ZONE --record RECORD
ionosctl dns r delete --all [--record PARTIAL_NAME] [--zone ZONE]
ionosctl dns r delete --record PARTIAL_NAME --zone ZONE
```

