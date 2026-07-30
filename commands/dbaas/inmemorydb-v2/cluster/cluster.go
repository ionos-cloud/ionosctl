package cluster

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	inmemorydb "github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// sha256HexRe matches a value that is already a SHA-256 hash (64 hex chars).
var sha256HexRe = regexp.MustCompile("^[a-fA-F0-9]{64}$")

// buildHashedPassword returns a HashedPassword for the given input. The In-Memory
// DB API only accepts hashed passwords, so a plaintext value is hashed with
// SHA-256 client-side; a value that already looks like a SHA-256 hash is passed
// through unchanged.
func buildHashedPassword(password string) inmemorydb.HashedPassword {
	if sha256HexRe.MatchString(password) {
		return inmemorydb.HashedPassword{Hash: password, Algorithm: "SHA-256"}
	}
	sum := sha256.Sum256([]byte(password))
	return inmemorydb.HashedPassword{Hash: hex.EncodeToString(sum[:]), Algorithm: "SHA-256"}
}

// applyCredentialsFromFlags rebuilds a cluster's credentials from the --password
// (and optional --user) flags. The API never returns the password on a GET, so a
// fetched cluster carries a password with an empty algorithm that a PUT would
// reject; update and restore must always re-send valid credentials. The existing
// username is kept unless --user overrides it.
func applyCredentialsFromFlags(c *core.CommandConfig, props *inmemorydb.Cluster) error {
	credentials := inmemorydb.ClusterCredentials{}
	if props.Credentials != nil {
		credentials = *props.Credentials
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgUser)) {
		credentials.Username = viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	}
	if credentials.Username == "" {
		return fmt.Errorf("could not determine username from the existing cluster; pass --%s", constants.ArgUser)
	}
	credentials.Password = buildHashedPassword(viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword)))
	props.Credentials = &credentials
	c.Verbose("Credentials - Username: %v", credentials.Username)
	return nil
}

var clusterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.name", Default: true},
	{Name: "DnsName", JSONPath: "metadata.dnsName", Default: true},
	{Name: "Version", JSONPath: "properties.version", Default: true},
	{Name: "Instances", JSONPath: "properties.instances.count"},
	{Name: "Cores", JSONPath: "properties.instances.cores"},
	{Name: "Ram", JSONPath: "properties.instances.ram"},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "EvictionPolicy", JSONPath: "properties.evictionPolicy", Default: true},
	{Name: "PersistenceMode", JSONPath: "properties.persistenceMode"},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "SnapshotLocation", JSONPath: "properties.snapshot.location"},
	{Name: "RetentionDays", JSONPath: "properties.snapshot.retentionDays"},
	{Name: "MaintenanceDay", JSONPath: "properties.maintenanceWindow.dayOfTheWeek"},
	{Name: "MaintenanceTime", JSONPath: "properties.maintenanceWindow.time"},
	{Name: "LogsEnabled", JSONPath: "properties.logsEnabled"},
	{Name: "MetricsEnabled", JSONPath: "properties.metricsEnabled"},
	{Name: "DatacenterId", JSONPath: "properties.connection.datacenterId"},
	{Name: "LanId", JSONPath: "properties.connection.lanId"},
	{Name: "Cidr", JSONPath: "properties.connection.primaryInstanceAddress"},
	{Name: "Username", JSONPath: "properties.credentials.username"},
	{Name: "StatusMessage", JSONPath: "metadata.statusMessage"},
}

func ClusterCmd() *core.Command {
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "In-Memory DB Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas in-memory-db-v2 cluster` allow you to manage the In-Memory DB Clusters under your account.",
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
