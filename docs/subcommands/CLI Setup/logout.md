---
description: Convenience command for deletion of config file credentials
---

# Logout

## Usage

```text
ionosctl logout [flags]
```

## Description

Convenience command for deletion of config file credentials

## Options

```text
  -u, --api-url string   Override default host url (default "https://api.ionos.com")
  -c, --config string    Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force            Force command to execute without user input
  -h, --help             Print usage
  -o, --output string    Desired output format [text|json] (default "text")
  -q, --quiet            Quiet output
  -v, --verbose          Print step-by-step process when running command
```

## Examples

```text
ionosctl login --user $IONOS_USERNAME --password $IONOS_PASSWORD

ionosctl login --token $IONOS_TOKEN

ionosctl login
Enter your username:
USERNAME
Enter your password:
```

