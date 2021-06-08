---
description: Delete a FlowLog from a NIC
---

# FlowlogDelete

## Usage

```text
ionosctl flowlog delete [flags]
```

## Aliases

For `flowlog` command:
```text
[fl]
```

For `delete` command:
```text
[d]
```

## Description

Use this command to delete a specified FlowLog from a NIC.

You can wait for the Request to be executed using `--wait-for-request` option. You can force the command to execute without user input using `--force` option.

Required values to run command:

* Data Center Id
* Server Id
* Nic Id
* FlowLog Id

## Options

```text
  -u, --api-url string         Override default API endpoint (default "https://api.ionos.com/cloudapi/v6")
      --cols strings           Set of columns to be printed on output 
                               Available columns: [FlowLogId Name Action Direction Bucket State] (default [FlowLogId,Name,Action,Direction,Bucket,State])
  -c, --config string          Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --datacenter-id string   The unique Data Center Id (required)
  -i, --flowlog-id string      The unique FlowLog Id (required)
  -f, --force                  Force command to execute without user input
  -h, --help                   help for delete
      --nic-id string          The unique NIC Id (required)
  -o, --output string          Desired output format [text|json] (default "text")
  -q, --quiet                  Quiet output
      --server-id string       The unique Server Id (required)
  -t, --timeout int            Timeout option for Request for FlowLog deletion [seconds] (default 60)
  -w, --wait-for-request       Wait for Request for FlowLog deletion to be executed
```

## Examples

```text
ionosctl flowlog delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --flowlog-id FLOWLOG_ID -f -w
```

