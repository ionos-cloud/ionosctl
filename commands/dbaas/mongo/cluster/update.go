package cluster

import (
	"context"
	"fmt"
	"net"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/templates"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/spf13/cobra"
)

func ClusterUpdateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Update a MongoDB Cluster",
		LongDesc: `
Use this command to update attributes of a MongoDB Cluster. To specify the cluster to update, use the ` + "`--cluster-id`" + ` flag and the cluster's unique ID you can get from the list command.

Every cluster can update:
* Maintenance window (day and time). To change any of these, you must specify both together (` + "`--maintenance-day`" + ` and ` + "`--maintenance-time`" + `).
* The display name of the cluster (` + "`--name`" + `).
* The MongoDB major version (` + "`--version`" + `). This can trigger a major upgrade of the cluster, so be sure to check the compatibility of your applications with the new version. Also see the notes in the [API Documentation](https://docs.ionos.com/cloud/databases/mongodb/api-howtos/modify-cluster-attributes/upgrade-the-mongodb-version).
* The backup storage location (` + "`--backup-location`" + `).

Replicaset clusters can update:
* The number of instances in the replicaset (` + "`--instances`" + `).

For enterprise edition clusters, you can also update:
* The memory for each MongoDB host system (` + "`--ram`" + `)
* The CPU Cores for each MongoDB host system (` + "`--cores`" + `)
* Storage size for each MongoDB instance (` + "`--storage-size`" + `)
* Storage type used for the Database (` + "`--storage-type`" + `)
* The number of shards (` + "`--shards`" + `). This is only possible for sharded clusters and requires a sharded_cluster type.
* The MongoDB Connector for Business Intelligence host and port (` + "`--biconnector`" + `) and whether it is enabled (` + "`--biconnector-enabled`" + `).

Business edition clusters currently cannot update their template size (which defines cores, RAM and storage size) this way. This can be done via DCD or API.

Fields which can only be updated under specific conditions:
* Network connection (CIDR, LAN, Datacenter) can only be updated if the amount of shards or instances changes and must be specified together with the new values. LAN and Datacenter must stay the same but need to be specified.
		`,
		Example: `# Rename a cluster and move its maintenance window (day+time must be given together)
ionosctl dbaas mongo cluster update --cluster-id <cluster-id> --name prod-mongo --maintenance-day Sunday --maintenance-time 03:00:00

# Scale an enterprise cluster vertically (per-instance cores/RAM/storage)
ionosctl dbaas mongo cluster update --cluster-id <cluster-id> --cores 4 --ram 8GB --storage-size 200GB`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr)
			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)

			return c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
		},
		CmdRun: func(c *core.CommandConfig) error {
			cluster := mongo.PatchClusterProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagEdition); viper.IsSet(fn) {
				cluster.Edition = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagType); viper.IsSet(fn) {
				cluster.Type = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagTemplateId); viper.IsSet(fn) {
				// Old flag kept for backwards compatibility. Behaviour fully included in --template flag
				cluster.TemplateID = pointer.From(viper.GetString(fn))
			} else {
				if fn := core.GetFlagName(c.NS, constants.FlagTemplate); viper.IsSet(fn) {
					tmplId, err := templates.Resolve(viper.GetString(fn))
					if err != nil {
						return err
					}
					cluster.TemplateID = pointer.From(tmplId)
				}
			}

			if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
				cluster.MongoDBVersion = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				cluster.DisplayName = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.IsSet(fn) {
				cluster.Instances = pointer.From(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagShards); viper.IsSet(fn) {
				cluster.Shards = pointer.From(viper.GetInt32(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.IsSet(fn) {
				if cluster.MaintenanceWindow == nil {
					cluster.MaintenanceWindow = &mongo.MaintenanceWindow{}
				}
				cluster.MaintenanceWindow.DayOfTheWeek = mongo.DayOfTheWeek(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.IsSet(fn) {
				if cluster.MaintenanceWindow == nil {
					cluster.MaintenanceWindow = &mongo.MaintenanceWindow{}
				}
				cluster.MaintenanceWindow.Time = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagCidr); viper.IsSet(fn) {
				if cluster.Connections == nil {
					cluster.Connections = make([]mongo.Connection, 1)
				}
				cluster.Connections[0].CidrList = viper.GetStringSlice(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDatacenterId); viper.IsSet(fn) {
				if cluster.Connections == nil {
					cluster.Connections = make([]mongo.Connection, 1)
				}
				cluster.Connections[0].DatacenterId = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLanId); viper.IsSet(fn) {
				if cluster.Connections == nil {
					cluster.Connections = make([]mongo.Connection, 1)
				}
				cluster.Connections[0].LanId = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, flagBackupLocation); viper.IsSet(fn) {
				if cluster.Backup == nil {
					cluster.Backup = &mongo.BackupProperties{}
				}
				cluster.Backup.Location = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, flagBiconnector); viper.IsSet(fn) {
				if cluster.BiConnector == nil {
					cluster.BiConnector = &mongo.BiConnectorProperties{}
				}
				hostAndPort := viper.GetString(fn)
				host, port, err := net.SplitHostPort(hostAndPort)
				if err != nil {
					return fmt.Errorf("failed splitting --%s %s into host and port: %w",
						flagBiconnector, hostAndPort, err)
				}
				cluster.BiConnector.Host = pointer.From(host)
				cluster.BiConnector.Port = pointer.From(port)
			}

			if fn := core.GetFlagName(c.NS, flagBiconnectorEnabled); viper.IsSet(fn) {
				if cluster.BiConnector == nil {
					cluster.BiConnector = &mongo.BiConnectorProperties{}
				}
				cluster.BiConnector.Enabled = pointer.From(viper.GetBool(fn))
			}

			// Enterprise flags
			if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.IsSet(fn) {
				cluster.Cores = pointer.From(viper.GetInt32(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageType); viper.IsSet(fn) {
				cluster.StorageType = (*mongo.StorageType)(pointer.From(viper.GetString(fn)))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.IsSet(fn) {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.MB)
				cluster.StorageSize = pointer.From(int32(sizeInt64))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagRam); viper.IsSet(fn) {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.MB)
				cluster.Ram = pointer.From(int32(sizeInt64))
			}

			clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
			createdCluster, _, err := client.Must().MongoClient.ClustersApi.ClustersPatch(context.Background(),
				clusterId).PatchClusterRequest(mongo.PatchClusterRequest{Properties: &cluster}).Execute()
			if err != nil {
				return fmt.Errorf("failed updating cluster: %w", err)
			}

			return c.Printer(allCols).Print(createdCluster)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagClusterId, constants.FlagIdShort, "", "The unique ID of the cluster to update", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterIds(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "New human-readable display name for the cluster")
	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "New replicaset size (one primary + n-1 secondaries). Valid values: 1, 3, 5, 7. Only for replicaset clusters")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagShards, "", 1, "New shard count for a sharded-cluster (enterprise only). Valid values: 2-32. Scaling shards DOWN is not supported")

	cmd.AddStringFlag(constants.FlagVersion, "", "", "Upgrade to this MongoDB major version, e.g. 6.0, 7.0. Patch versions are managed automatically. Downgrades are NOT supported; a major upgrade can affect application compatibility")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(_ *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoClusterVersions(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "New start time (UTC, HH:MM:SS) of the weekly 4-hour maintenance window, e.g. 16:30:59. Must be changed together with --maintenance-day")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "New day of the week for the weekly 4-hour maintenance window, e.g. Saturday. Must be changed together with --maintenance-time")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Enterprise-specific
	cmd.AddIntFlag(constants.FlagCores, "", 0, "New CPU core count PER instance, e.g. 4. Minimum 1. Enterprise only")
	cmd.AddStringFlag(constants.FlagRam, "", "", "New RAM PER instance. Multiple of 1024 MB (1 GB), minimum 2 GB. Accepts a unit, e.g. --ram 4GB. Enterprise only")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1024MB", "2GB", "4GB", "8GB", "12GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageType, "", "\"SSD Standard\"",
		"New disk type: HDD, 'SSD Standard' or 'SSD Premium'. Enterprise only")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "\"SSD Standard\"", "\"SSD Premium\""}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", "", "New disk size PER instance, e.g. 200GB. Storage can only be grown, not shrunk; better performance above 100 GB. Enterprise only")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"50GB", "100GB", "200GB", "400GB", "800GB", "1TB", "2TB"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Connections
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "Datacenter of the cluster's connection. Only editable while also changing the instance/shard count; the datacenter and LAN must stay the same but must be re-supplied together with --lan-id and --cidr")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "Numeric LAN ID of the cluster's connection. Must be supplied together with --datacenter-id and --cidr when changing the connection")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceFlag(constants.FlagCidr, "", nil, "New comma-separated private IPs (with /24 subnet), one per instance. Editable only when the instance/shard count changes, and must be supplied with --datacenter-id and --lan-id. Unavailable range: 10.233.114.0/24")

	cmd.AddStringFlag(flagBackupLocation, "", "", "Object Storage region where snapshots are stored, e.g. de, eu-south-2, eu-central-3. Defaults to the region nearest the cluster")
	// Biconnector
	cmd.AddStringFlag(flagBiconnector, "", "", "New host:port for the MongoDB BI Connector (SQL access for BI tools). Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015")
	cmd.AddBoolFlag(flagBiconnectorEnabled, "", false, fmt.Sprintf("Enable or disable the BI Connector. If left unset, its state is unchanged. To explicitly disable it, use --%s=false", flagBiconnectorEnabled))

	// Misc

	cmd.Command.SilenceUsage = true

	return cmd
}
