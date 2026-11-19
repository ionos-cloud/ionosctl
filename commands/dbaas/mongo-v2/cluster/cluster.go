package cluster

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	mongo "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editions = []string{string(mongo.MONGOEDITION_ENTERPRISE), string(mongo.MONGOEDITION_BUSINESS)}

var enterpriseOnlyFlags = []string{constants.FlagCores, constants.FlagRam, constants.FlagStorageSize, constants.FlagShards, constants.FlagBiConnectorEnabled}

// clusterProps holds the accessors shared by the ENTERPRISE and BUSINESS variants (read, PUT and create).
type clusterProps interface {
	GetName() string
	SetName(v string)
	SetDescription(v string)
	GetVersion() string
	SetVersion(v string)
	GetConnection() mongo.MongoClusterConnection
	SetConnection(v mongo.MongoClusterConnection)
	GetMaintenanceWindow() mongo.MaintenanceWindow
	SetMaintenanceWindow(v mongo.MaintenanceWindow)
	GetCredentials() mongo.MongoUser
	SetCredentials(v mongo.MongoUser)
	SetRestoreFromSnapshot(v mongo.MongoRestoreFromSnapshot)
	SetLogsEnabled(v bool)
	SetMetricsEnabled(v bool)
	GetSnapshot() mongo.ClusterSnapshot
	SetSnapshot(v mongo.ClusterSnapshot)
}

func editionProps(cluster *mongo.Cluster) (clusterProps, error) {
	switch {
	case cluster.EnterpriseCluster != nil:
		return cluster.EnterpriseCluster, nil
	case cluster.BusinessCluster != nil:
		return cluster.BusinessCluster, nil
	}
	return nil, errors.New("the API returned a cluster with neither ENTERPRISE nor BUSINESS properties")
}

// setFlags returns, as "--name", those of the given flags that were set.
func setFlags(ns string, flags []string) []string {
	var set []string
	for _, flag := range flags {
		if viper.IsSet(core.GetFlagName(ns, flag)) {
			set = append(set, "--"+flag)
		}
	}
	return set
}

func validateCredentialFlags(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword)) {
		return nil
	}
	if set := setFlags(c.NS, []string{constants.ArgUser, constants.FlagDatabase}); len(set) > 0 {
		return fmt.Errorf("--%s is required with %s: the API only accepts the initial user's credentials together with its password",
			constants.ArgPassword, strings.Join(set, " and "))
	}
	return nil
}

func clearCredentials(cluster *mongo.Cluster) {
	if cluster.EnterpriseCluster != nil {
		cluster.EnterpriseCluster.Credentials = nil
	}
	if cluster.BusinessCluster != nil {
		cluster.BusinessCluster.Credentials = nil
	}
}

func applyCredentialsFromFlags(c *core.CommandConfig, cluster *mongo.Cluster) error {
	props, err := editionProps(cluster)
	if err != nil {
		return err
	}

	if !viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword)) {
		clearCredentials(cluster)
		c.Verbose("Credentials - unchanged")
		return nil
	}

	credentials := props.GetCredentials()
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgUser)) {
		credentials.Username = viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatabase)) {
		credentials.Database = viper.GetString(core.GetFlagName(c.NS, constants.FlagDatabase))
	}
	if credentials.Username == "" || credentials.Database == "" {
		return fmt.Errorf("could not determine username/database from the existing cluster; pass --%s and --%s", constants.ArgUser, constants.FlagDatabase)
	}
	password := viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword))
	credentials.Password = &password
	props.SetCredentials(credentials)
	c.Verbose("Credentials - Username: %v, Database: %v", credentials.Username, credentials.Database)
	return nil
}

func applyOptionalFlags(c *core.CommandConfig, props clusterProps) {
	if fn := core.GetFlagName(c.NS, constants.FlagDescription); viper.IsSet(fn) {
		desc := viper.GetString(fn)
		c.Verbose("Description: %v", desc)
		props.SetDescription(desc)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagLogsEnabled); viper.IsSet(fn) {
		logsEnabled := viper.GetBool(fn)
		c.Verbose("LogsEnabled: %v", logsEnabled)
		props.SetLogsEnabled(logsEnabled)
	}
	if fn := core.GetFlagName(c.NS, constants.FlagMetricsEnabled); viper.IsSet(fn) {
		metricsEnabled := viper.GetBool(fn)
		c.Verbose("MetricsEnabled: %v", metricsEnabled)
		props.SetMetricsEnabled(metricsEnabled)
	}
}

func applyEnterpriseFeatureFlags(c *core.CommandConfig, cluster interface {
	SetSharding(v mongo.MongoSharding)
	SetBiConnectorEnabled(v bool)
}) {
	if fn := core.GetFlagName(c.NS, constants.FlagShards); viper.IsSet(fn) {
		shards := viper.GetInt32(fn)
		c.Verbose("Sharding - Shards: %v", shards)
		cluster.SetSharding(mongo.MongoSharding{Enabled: pointer.From(true), Shards: &shards})
	}
	if fn := core.GetFlagName(c.NS, constants.FlagBiConnectorEnabled); viper.IsSet(fn) {
		biConnectorEnabled := viper.GetBool(fn)
		c.Verbose("BiConnectorEnabled: %v", biConnectorEnabled)
		cluster.SetBiConnectorEnabled(biConnectorEnabled)
	}
}

func validateRetentionDays(days int32) error {
	if days < 1 || days > 365 {
		return fmt.Errorf("--%s must be between 1 and 365 (got %d)", constants.FlagRetentionDays, days)
	}
	return nil
}

func validateInstanceAddresses(addresses []string, instances int32, sharded bool) error {
	if sharded {
		if len(addresses) != 3 {
			return fmt.Errorf("a sharded cluster takes exactly 3 --%s addresses, one per config server (got %d)", constants.FlagCidr, len(addresses))
		}
		return nil
	}
	if len(addresses) != int(instances) {
		return fmt.Errorf("--%s takes one address per instance: got %d address(es) for %d instance(s)", constants.FlagCidr, len(addresses), instances)
	}
	return nil
}

func toInt32GB(raw, flag string) (int32, error) {
	size, ok := convbytes.StrToUnitOk(raw, convbytes.GB)
	if !ok {
		return 0, fmt.Errorf("invalid value for --%s: %v", flag, raw)
	}
	if size < 0 || size > math.MaxInt32 {
		return 0, fmt.Errorf("--%s value %vGB exceeds accepted int32 range: 0 - %d", flag, size, math.MaxInt32)
	}
	return int32(size), nil
}

func instanceCount(item map[string]any) any {
	if count := table.Navigate(item, "properties.instances.count"); count != nil {
		return count
	}
	return table.Navigate(item, "properties.instanceCount")
}

var clusterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "Edition", JSONPath: "properties.edition", Default: true},
	{Name: "Version", JSONPath: "properties.version", Default: true},
	{Name: "Instances", Format: instanceCount, Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "ConnectionUri", JSONPath: "metadata.connectionUri", Default: true},
	{Name: "Cores", JSONPath: "properties.instances.cores"},
	{Name: "Ram", JSONPath: "properties.instances.ram"},
	{Name: "StorageSize", JSONPath: "properties.instances.storageSize"},
	{Name: "TemplateId", JSONPath: "properties.instanceTemplateId"},
	{Name: "ShardingEnabled", JSONPath: "properties.sharding.enabled"},
	{Name: "Shards", JSONPath: "properties.sharding.shards"},
	{Name: "BiConnectorEnabled", JSONPath: "properties.biConnectorEnabled"},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "SnapshotLocation", JSONPath: "properties.snapshot.location"},
	{Name: "RetentionDays", JSONPath: "properties.snapshot.retentionDays"},
	{Name: "MaintenanceDay", JSONPath: "properties.maintenanceWindow.dayOfTheWeek"},
	{Name: "MaintenanceTime", JSONPath: "properties.maintenanceWindow.time"},
	{Name: "LogsEnabled", JSONPath: "properties.logsEnabled"},
	{Name: "MetricsEnabled", JSONPath: "properties.metricsEnabled"},
	{Name: "DatacenterId", JSONPath: "properties.connection.datacenterId"},
	{Name: "LanId", JSONPath: "properties.connection.lanId"},
	{Name: "Cidr", JSONPath: "properties.connection.instanceAddresses"},
	{Name: "Username", JSONPath: "properties.credentials.username"},
	{Name: "Database", JSONPath: "properties.credentials.database"},
	{Name: "StatusMessage", JSONPath: "metadata.statusMessage"},
}

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "cluster",
			Aliases: []string{"c"},
			Short:   "MongoDB Cluster Operations",
			Long: `The sub-commands of ` + "`ionosctl dbaas mongo-v2 cluster`" + ` allow you to manage the MongoDB Clusters under your account.

Every cluster has an edition, chosen at creation and fixed afterwards:
  - ENTERPRISE: each instance is sized explicitly (--cores, --ram, --storage-size). Supports sharding (--shards) and the MongoDB BI Connector (--bi-connector-enabled).
  - BUSINESS: each instance is sized by a predefined instance template (--template-id, see ` + "`ionosctl dbaas mongo-v2 template list`" + `), also fixed at creation.

A cluster runs one primary and n-1 secondary instances, reachable over a private LAN (--datacenter-id, --lan-id) at one IP address per instance (--cidr). Snapshots are stored in an Object Storage location and allow point-in-time recovery (see ` + "`snapshot`" + ` and ` + "`cluster restore`" + `).`,
			TraverseChildren: true,
		},
	}

	clusterCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, table.ColsMessage(clusterCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return table.AllCols(clusterCols), cobra.ShellCompDirectiveNoFileComp
		},
	)

	clusterCmd.AddCommand(ClusterListCmd())
	clusterCmd.AddCommand(ClusterCreateCmd())
	clusterCmd.AddCommand(ClusterUpdateCmd())
	clusterCmd.AddCommand(ClusterDeleteCmd())
	clusterCmd.AddCommand(ClusterGetCmd())
	clusterCmd.AddCommand(ClusterRestoreCmd())

	return clusterCmd
}
