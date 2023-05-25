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

The 'login' command allows you to authenticate with the IONOS Cloud APIs. There are three ways you can use it:
  1. Interactive mode: Just type 'ionosctl login' and you'll be prompted to enter your username and password.
  2. Use the '--user' and '--password' flags: Enter your credentials in the command.
  3. Use the '--token' flag: Provide an authentication token.

If you use a username and password, this command generates a token that's saved in the config file. Please keep this token safe. If you specify a custom '--api-url', it'll be saved to the config file when you login successfully and used for future API calls.

By default, the config file is located at /home/avirtopeanu/.config/ionosctl/config.json. If you want to use a different config file, use the '--config' global option. Changing the permissions of the config file will cause it to no longer work.

Note: The IONOS Cloud CLI supports also authentication with environment variables: $IONOS_USERNAME, $IONOS_PASSWORD or $IONOS_TOKEN, these override the config file token.

## Options

```text
  -u, --api-url string        Override default host url (default "https://api.ionos.com")
  -c, --config string         Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
  -f, --force                 Force command to execute without user input
  -h, --help                  Print usage
  -o, --output string         Desired output format [text|json] (default "text")
  -p, --password string       Password to authenticate
  -q, --quiet                 Quiet output
  -t, --token string          Token to authenticate
      --use-default-api-url   Use the default authentication URL (https://api.ionos.com) for auth checking, even if you specify a different '--api-url'
      --user string           Username to authenticate
  -v, --verbose               Print step-by-step process when running command
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

