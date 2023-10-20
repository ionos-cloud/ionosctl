---
description: "Create a FlowLog on a NIC"
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

NOTE: Please disable the FlowLog before deleting the existing Bucket.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* Target S3 Bucket Name

## Options

```text
  -a, --action string          Specifies the traffic Action pattern (default "ALL")
  -u, --api-url string         Override default host url (default "https://api.ionos.com")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -D, --depth int32            Controls the detail depth of the response objects. Max depth is 10.
  -d, --direction string       Specifies the traffic Direction pattern (default "BIDIRECTIONAL")
  -f, --force                  Force command to execute without user input
  -h, --help                   Print usage
  -n, --name string            The name for the FlowLog (default "Unnamed FlowLog")
      --nic-id string          The unique NIC Id (required)
      --no-headers             Don't print table headers when table output is used
  -o, --output string          Desired output format [text|json|api-json] (default "text")
  -q, --quiet                  Quiet output
  -b, --s3bucket string        S3 Bucket name of an existing IONOS Cloud S3 Bucket (required)
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for FlowLog creation [seconds] (default 60)
  -v, --verbose                Print step-by-step process when running command
  -w, --wait-for-request       Wait for Request for FlowLog creation to be executed
```

## Examples

```text
ionosctl flowlog create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name NAME --action ACTION --direction DIRECTION --s3bucket BUCKET_NAME
```

