---
description: Create a FlowLog
---

# FlowlogCreate

## Usage

```text
ionosctl flowlog create [flags]
```

## Aliases

For `flowlog` command:
```text
[fl]
```

For `create` command:
```text
[c]
```

## Description

Use this command to create a new FlowLog to the specified NIC.

You can wait for the Request to be executed using `--wait-for-request` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id 
* Name
* Direction
* Action
* Target S3 Bucket Name

## Options

```text
  -a, --action string          Specifies the traffic Action pattern (required)
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --bucket-name string     S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -d, --direction string       Specifies the traffic Direction pattern (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for create
  -n, --name string            The name for the FlowLog (required)
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for FlowLog creation [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for FlowLog creation to be executed
```

## Examples

```text
ionosctl flowlog create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name NAME --action ACTION --direction DIRECTION --bucket-name BUCKET_NAME
```

