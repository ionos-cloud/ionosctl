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

Use this command to authenticate. You can use  `--user` and `--password` flags or you can use  `--token` flag to set the credentials.

By default, the user data after running this command will be saved in:

* macOS: `${HOME}/Library/Application Support/ionosctl/config.json`
* Linux: `${XDG_CONFIG_HOME}/ionosctl/config.json`
* Windows: `%APPDATA%\ionosctl\config.json`.

You can use another configuration file for authentication with the `--config` global option.

Note: The IONOS Cloud CLI supports also authentication with environment variables: $IONOS_USERNAME, $IONOS_PASSWORD or $IONOS_TOKEN.

## Options

```text
  -p, --password string   Password to authenticate
  -t, --token string      Token to authenticate
      --user string       Username to authenticate
```

## Examples

```text
ionosctl login --user $IONOS_USERNAME --password $IONOS_PASSWORD
Status: Authentication successful!

ionosctl login --token $IONOS_TOKEN
Status: Authentication successful!

ionosctl login 
Enter your username:
USERNAME
Enter your password:

Status: Authentication successful!
```

