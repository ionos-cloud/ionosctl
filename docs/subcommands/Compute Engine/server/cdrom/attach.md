---
description: "Attach a CD-ROM to a Server"
---

# ServerCdromAttach

## Usage

```text
ionosctl server cdrom attach [flags]
```

## Aliases

For `server` command:

```text
[s svr]
```

For `cdrom` command:

```text
[cd]
```

For `attach` command:

```text
[a]
```

## Description

Use this command to attach a CD-ROM to an existing Server.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id

## Options

```text
  -u, --api-url string         Override default host URL. Preferred over the config file override 'cloud' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
  -i, --cdrom-id string        The unique Cdrom Id (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId ExposeSerial RequireLegacyBios ApplicationType] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int              Level of detail for response objects (default 1)
  -F, --filters strings        Limit results to results containing the specified filter:KEY1=VALUE1,KEY2=VALUE2
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
      --limit int              Maximum number of items to return per request (default 50)
      --no-headers             Don't print table headers when table output is used
      --offset int             Number of items to skip before starting to collect the results
      --order-by string        Property to order the results by
  -o, --output string          Desired output format [text|json|api-json] (default "text")
      --query string           JMESPath query string to filter the output
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for Cdrom attachment [seconds] (default 60)
  -v, --verbose count          Increase verbosity level [-v, -vv, -vvv]
  -w, --wait-for-request       Wait for the Request for CD-ROM attachment to be executed
```

## Examples

```text
ionosctl server cdrom attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID --wait-for-request
```

