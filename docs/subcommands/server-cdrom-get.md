---
description: Get a specific attached CD-ROM from a Server
---

# ServerCdromGet

## Usage

```text
ionosctl server cdrom get [flags]
```

## Description

Use this command to retrieve information about an attached CD-ROM on Server.

Required values to run command:

* Data Center Id
* Server Id
* Cdrom Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cdrom-id string        The unique Cdrom Id (required)
      --cols strings           Columns to be printed in the standard output (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for get
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
```

## Examples

```text
ionosctl server cdrom get --datacenter-id 4fd7996d-2b08-4c04-9c47-d9d884ee179a --server-id f7438b0c-2f36-4bec-892f-af027930b81e --cdrom-id 80c63662-49a0-11ea-94e0-525400f64d8d 
ImageId                                Name                              ImageAliases                       Location   LicenceType   ImageType   CloudInit
80c63662-49a0-11ea-94e0-525400f64d8d   CentOS-8.1.1911-x86_64-boot.iso   [centos:latest_iso centos:8_iso]   us/las     LINUX         CDROM       NONE
```

