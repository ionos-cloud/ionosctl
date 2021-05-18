---
description: Attach a CD-ROM to a Server
---

# ServerCdromAttach

## Usage

```text
ionosctl server cdrom attach [flags]
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
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v5")
      --cdrom-id string        The unique Cdrom Id (required)
      --cols strings           Columns to be printed in the standard output (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
      --force                  Force command to execute without user input
  -h, --help                   help for attach
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
      --timeout int            Timeout option for Request for Cdrom attachment [seconds] (default 60)
      --wait-for-request       Wait for the Request for CD-ROM attachment to be executed
```

## Examples

```text
ionosctl server cdrom attach --datacenter-id 4fd7996d-2b08-4c04-9c47-d9d884ee179a --server-id f7438b0c-2f36-4bec-892f-af027930b81e --cdrom-id 99d43e40-49a0-11ea-94e0-525400f64d8d --wait-for-request 
1.4s Waiting for request... FAILED                                                                                                                                                                         
Error: FAILED [VDC-5-622] Image 99d43e40-49a0-11ea-94e0-525400f64d8d is from location us/ewr but Virtual Data Center is from another location us/las

ionosctl server cdrom attach --datacenter-id 4fd7996d-2b08-4c04-9c47-d9d884ee179a --server-id f7438b0c-2f36-4bec-892f-af027930b81e --cdrom-id 80c63662-49a0-11ea-94e0-525400f64d8d --wait-for-request 
13s Waiting for request..  DONE                                                                                                                                                                            
ImageId                                Name                              ImageAliases                       Location   LicenceType   ImageType   CloudInit
80c63662-49a0-11ea-94e0-525400f64d8d   CentOS-8.1.1911-x86_64-boot.iso   [centos:latest_iso centos:8_iso]   us/las     LINUX         CDROM       NONE
RequestId: 3f63b766-b27c-42a4-b421-6e2cfb57a877
Status: Command cdrom attach & wait have been successfully executed
```

