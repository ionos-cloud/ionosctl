---
description: List Images
---

# ImageList

## Usage

```text
ionosctl image list [flags]
```

## Aliases

For `image` command:
```text
[img]
```

## Description

Use this command to get a list of available public Images. Use flags to retrieve a list of sorted images by location, licence type, type or size.

## Options

```text
  -u, --api-url string        Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                 Force command to execute without user input
  -F, --format strings        Collection of fields to be printed on output (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit])
  -h, --help                  help for list
      --licence-type string   The licence type of the Image
  -l, --location string       The location of the Image
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --size float32          The size of the Image
      --type string           The type of the Image
```

## Examples

```text
ionosctl image list --location us/las --type HDD
ImageId                                Name                                 Location   Size   LicenceType   ImageType
8991cf6c-8706-11eb-a1d6-72dfddd36b99   windows-2012-r2-server-2021-03       us/las     14     WINDOWS       HDD
7ab978cb-870a-11eb-a1d6-72dfddd36b99   windows-2016-server-2021-03          us/las     14     WINDOWS2016   HDD
aca2279d-870b-11eb-a1d6-72dfddd36b99   windows-2019-server-2021-03          us/las     15     WINDOWS2016   HDD
33915e02-9291-11eb-b68e-9ad3ea4b1420   CentOS-7-server-2021-04-01           us/las     4      LINUX         HDD
bc81846b-929c-11eb-b68e-9ad3ea4b1420   Debian-testing-server-2021-04-01     us/las     2      LINUX         HDD
8eca30f5-92a0-11eb-b68e-9ad3ea4b1420   Ubuntu-16.04-LTS-server-2021-04-01   us/las     3      LINUX         HDD
f1316493-92a4-11eb-b68e-9ad3ea4b1420   Ubuntu-18.04-LTS-server-2021-04-01   us/las     3      LINUX         HDD
a86dd807-9297-11eb-b68e-9ad3ea4b1420   Debian-10-server-2021-04-01          us/las     2      LINUX         HDD
4bf9c128-929a-11eb-b68e-9ad3ea4b1420   Debian-9-server-2021-04-01           us/las     2      LINUX         HDD
378ecd34-9295-11eb-b68e-9ad3ea4b1420   CentOS-8-server-2021-04-01           us/las     4      LINUX         HDD
02e06d34-92aa-11eb-b68e-9ad3ea4b1420   Ubuntu-20.04-LTS-server-2021-04-01   us/las     3      LINUX         HDD
81adeb01-3379-11eb-a681-1e659523cb7b   CentOS-6-server-2020-12-01           us/las     2      LINUX         HDD
6f9bae91-3386-11eb-a681-1e659523cb7b   Debian-8-server-2020-12-01           us/las     2      LINUX         HDD
8fc5f591-338e-11eb-a681-1e659523cb7b   Ubuntu-19.10-server-2020-12-01       us/las     3      LINUX         HDD
```

