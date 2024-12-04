---
description: "Create DBaaS MariaDB clusters"
---

# DbaasMariadbClusterCreate

## Usage

```text
ionosctl dbaas mariadb cluster create [flags]
```

## Aliases

For `mariadb` command:

```text
[maria mar ma]
```

For `cluster` command:

```text
[c]
```

For `create` command:

```text
[c]
```

## Description

Create DBaaS MariaDB clusters

## Options

```text
  -u, --api-url string            Override default host URL (default "https://mariadb.de-txl.ionos.com")
      --cidr string               The IP and subnet for your cluster. All IPs must be in a /24 network (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name DNS Instances Version State Cores RAM StorageSize MaintenanceDay MaintenanceTime] (default [ClusterId,Name,DNS,Instances,Version,State])
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.json")
      --cores int32               Core count (default 1)
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster (required)
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances of the cluster (one primary and n-1 secondaries) (default 1)
      --lan-id string             The numeric LAN ID with which you connect your cluster (required)
  -l, --location string           Location of the resource to operate on. Can be one of: de/txl, de/fra, es/vit, fr/par, gb/lhr, us/ewr, us/las, us/mci
      --maintenance-day string    Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
      --maintenance-time string   Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59. Defaults to a random day during Mon-Fri, during the hours 10:00-16:00 (default "Random (Mon-Fri 10:00-16:00)")
  -n, --name string               The name of your cluster (required)
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
      --password string           The password (required)
  -q, --quiet                     Quiet output
      --ram string                RAM size. e.g.: --ram 4GB. Minimum of 4GB. The maximum RAM size is determined by your contract limit (default "4GB")
      --storage-size string       The size of the Storage in GB. e.g.: --size 10 or --size 10GB. The maximum Volume size is determined by your contract limit (default "10")
      --user string               The initial username (required)
  -v, --verbose                   Print step-by-step process when running command
      --version string            The MariaDB version of your cluster (required) (default "10.6")
```

## Examples

```text
i db mariadb cluster create --name NAME --version VERSION --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --user USER --password PASSWORD 
```

