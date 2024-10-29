---
description: "Retrieve all distributions using pagination and optional filters"
---

# CdnDistributionList

## Usage

```text
ionosctl cdn distribution list [flags]
```

## Aliases

For `distribution` command:

```text
[ds]
```

For `list` command:

```text
[ls]
```

## Description

Retrieve all distributions using pagination and optional filters

## Options

```text
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [Id Domain CertificateId State]
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --filter.domain string   Filter used to fetch only the records that contain specified domain.
      --filter.state string    Filter used to fetch only the records that contain specified state.. Can be one of: AVAILABLE, BUSY, FAILED, UNKNOWN
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --max-results int32      The maximum number of elements to return
      --no-headers             Don't print table headers when table output is used
      --offset int32           The first element (of the total list of elements) to include in the response. Use together with limit for pagination
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -v, --verbose                Print step-by-step process when running command
```

## Examples

```text
ionosctl cdn ds list
```

