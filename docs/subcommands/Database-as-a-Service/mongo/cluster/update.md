---
description: "Update a MongoDB Cluster"
---

# DbaasMongoClusterUpdate

## Usage

```text
ionosctl dbaas mongo cluster update [flags]
```

## Aliases

For `mongo` command:

```text
[m mdb mongodb mg]
```

For `cluster` command:

```text
[c]
```

For `update` command:

```text
[u]
```

## Description


		Use this command to update attributes of a MongoDB Cluster. To specify the cluster to update, use the --cluster-id flag and the cluster's unique ID you can get from the list command.

		Every cluster can update:
		* Network connection (IP Addresses, LAN, Datacenter). You can only change the CIDR, but need to specify the LAN and Datacenter together with it (--cidr, --lan-id, --datacenter-id).
		* Maintenance window (day and time). To change any of these, you must specify both together (--maintenance-day and --maintenance-time).
		* The display name of the cluster (--name).
		* The MongoDB major version (--version). This can trigger a major upgrade of the cluster, so be sure to check the compatibility of your applications with the new version. Also see the notes in the [API Documentation](https://docs.ionos.com/cloud/databases/mongodb/api-howtos/modify-cluster-attributes/upgrade-the-mongodb-version).
		* The backup storage location (--backup-location).

		Replicaset clusters can update:
		* The number of instances in the replicaset (--instances).

		For enterprise edition clusters, you can also update:
		* The memory for each MongoDB host system (--ram)
		* The CPU Cores for each MongoDB host system (--cores)
		* Storage size for each MongoDB instance (--storage-size)
		* Storage type used for the Database (--storage-type)
		* The number of shards (--shards). This is only possible for sharded clusters and requires a sharded_cluster type.
		* The MongoDB Connector for Business Intelligence host and port (--biconnector) and whether it is enabled (--biconnector-enabled).
		
		Business edition clusters currently cannot update their template size (which defines cores, RAM and storage size) this way. This can be done via DCD or API.
		

## Options

```text
  -u, --api-url string            Override default host URL. Preferred over the config file override 'mongo' and env var 'IONOS_API_URL' (default "https://api.ionos.com")
      --backup-location string    The location where the cluster backups will be stored. If not set, the backup is stored in the backup location nearest to the cluster
      --biconnector string        The host and port where this new BI Connector is installed. The MongoDB Connector for Business Intelligence allows you to query a MongoDB database using SQL commands. Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015
      --biconnector-enabled       Enable or disable the biconnector. If left unset, no change will be made to the biconnector's status. To explicitly disable it, use --biconnector-enabled=false
      --cidr strings              The list of IPs and subnet for your cluster. All IPs must be in a /24 network. Note the following unavailable IP range: 10.233.114.0/24
  -i, --cluster-id string         The unique ID of the cluster (required)
      --cols strings              Set of columns to be printed on output 
                                  Available columns: [ClusterId Name Edition Type URL Instances Shards Health State MongoVersion MaintenanceWindow Location DatacenterId LanId Cidr TemplateId Cores RAM StorageSize StorageType]
  -c, --config string             Configuration file used for authentication (default "$XDG_CONFIG_HOME/ionosctl/config.yaml")
      --cores int                 The total number of cores for the Server, e.g. 4. (only settable for enterprise edition)
      --datacenter-id string      The datacenter to which your cluster will be connected. Must be in the same location as the cluster
  -f, --force                     Force command to execute without user input
  -h, --help                      Print usage
      --instances int32           The total number of instances of the cluster (one primary and n-1 secondaries). Minimum of 3 for business edition (default 1)
      --lan-id string             The numeric LAN ID with which you connect your cluster
      --maintenance-day string    Day for Maintenance. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: Saturday
      --maintenance-time string   Time for the Maintenance. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59
  -n, --name string               The name of your cluster
      --no-headers                Don't print table headers when table output is used
  -o, --output string             Desired output format [text|json|api-json] (default "text")
  -q, --quiet                     Quiet output
      --ram string                Custom RAM: multiples of 1024. e.g. --ram 1024 or --ram 1024MB or --ram 4GB (only settable for enterprise edition)
      --shards int32              The total number of shards in the sharded_cluster cluster. Setting this flag is only possible for enterprise clusters and requires a sharded_cluster type. Possible values: 2 - 32. Scaling down is not supported. (default 1)
      --storage-size string       Custom Storage: Greater performance for values greater than 100 GB. (only settable for enterprise edition)
      --storage-type string       Custom Storage Type. (only settable for enterprise edition) (default "\"SSD Standard\"")
  -v, --verbose                   Print step-by-step process when running command
      --version string            The MongoDB version of your cluster. This only accepts the major version, e.g. 6.0, 7.0, etc. Patch versions are set automatically. Downgrades are not supported.
```

## Examples

```text
ionosctl dbaas mongo cluster update --cluster-id <cluster-id> --version <new-version>
```

