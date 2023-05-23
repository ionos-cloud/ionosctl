---
description: Authentication command for SDK
---

# Login

## Usage

```text
ionosctl login [flags]
```

## Aliases

For `login` command:

```text
[log auth]
```

## Description

Use this command to authenticate.
You can either use the interactive mode, or you can use "--user" and "--password" flags or "--token" flag to set the credentials.
If using username & password, this command will generate a JWT token which will be saved in the config file. Please safeguard your token.
The config file, by default, will be created at /home/avirtopeanu/.config/ionosctl/config.json. You can use another configuration file for authentication with the "--config" global option.

Note: The IONOS Cloud CLI supports also authentication with environment variables: $IONOS_USERNAME, $IONOS_PASSWORD or $IONOS_TOKEN, these override the config file token.

## Options

```text
  -u, --api-url string    Override default host url (default "https://api.ionos.com")
  -c, --config string     Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force             Force command to execute without user input
  -h, --help              Print usage
  -o, --output string     Desired output format [text|json] (default "text")
  -p, --password string   Password to authenticate
  -q, --quiet             Quiet output
  -t, --token string      Token to authenticate
      --user string       Username to authenticate
  -v, --verbose           Print step-by-step process when running command
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

