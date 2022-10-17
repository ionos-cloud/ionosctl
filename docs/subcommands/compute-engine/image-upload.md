---
description: Upload an image to FTP server
---

# ImageUpload

## Usage

```text
ionosctl image upload [flags]
```

## Aliases

For `image` command:

```text
[img]
```

For `upload` command:

```text
[ftp-upload ftp upl]
```

## Description

Use this command to upload an HDD or ISO image.

Required values to run command:

* Location


## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
      --cols strings          Set of columns to be printed on output 
                              Available columns: [ImageId Name ImageAliases Location Size LicenceType ImageType Description Public CloudInit CreatedDate CreatedBy CreatedByUserId] (default [ImageId,Name,ImageAliases,Location,LicenceType,ImageType,CloudInit,CreatedDate])
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --crt-path string       (Unneeded for IONOS FTP Servers) Path to file containing server certificate. If your FTP server is self-signed, you need to add the server certificate to the list of certificate authorities trusted by the client.
  -f, --force                 Force command to execute without user input
      --ftp-url string        URL of FTP server, with %s flag if location is embedded into url (default "ftp-%s.ionos.com")
  -h, --help                  Print usage
  -i, --image strings         Slice of paths to images, absolute path or relative to ionosctl binary. (required)
  -a, --image-alias strings   Rename the uploaded images. These names should not contain any extension. By default, this is the base of the image path
  -l, --location strings      Location to upload to. Must be an array containing only fra, fkb, txl, lhr, las, ewr, vit (required)
  -o, --output string         Desired output format [text|json] (default "text")
  -q, --quiet                 Quiet output
      --skip-verify           Skip verification of server certificate, useful if using a custom ftp-url. WARNING: You can be the target of a man-in-the-middle attack!
      --timeout int           (seconds) Context Deadline. FTP connection will time out after this many seconds (default 300)
  -v, --verbose               Print step-by-step process when running command
```

## Examples

```text
ionosctl img u -i kolibri.iso -l fkb,fra,vit
```

