package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mariadb/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	baseReqFlags := []string{
		constants.FlagName, constants.FlagVersion, constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr,
		constants.ArgUser, constants.ArgPassword,
	}
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mariadb",
		Resource:  "cluster",
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create a DBaaS MariaDB cluster",
		LongDesc: `Create a new managed MariaDB cluster.

A cluster is a set of instances running the same MariaDB version: one primary (read-write) and, for a replicated setup, n-1 secondaries (read-only replicas kept in sync with the primary). The instance count must be odd so a quorum can elect a primary: use --instances 1 for a single standalone instance (no high availability), or 3 or 5 for a high-availability replica set. Every instance is sized identically via --cores, --ram and --storage-size.

The cluster is reached only over a private LAN inside one of your VDCs, never over the public internet. --datacenter-id, --lan-id and --cidr together define that connection; the datacenter must be in the same location (region) as the cluster, exactly one connection is allowed, and the --cidr address must sit inside the chosen LAN's subnet. After creation the cluster is reachable at the DNS name shown in 'cluster get'.

--user and --password create the initial database user; the password is write-only and is never returned by the API afterwards, so store it safely. A newly created cluster starts in state BUSY and must reach AVAILABLE before you can connect to it or run further operations (update, restore) against it.

To provision a cluster from an existing backup instead of empty, use the API's fromBackup option; in ionosctl, inspect available backups with 'ionosctl dbaas mariadb backup list'.`,
		Example: fmt.Sprintf("# Create a standalone MariaDB cluster (MariaDB 10.6, 1 instance)\ni db mariadb cluster create %s\n\n# Create a high-availability cluster (3 instances) with explicit sizing and a fixed maintenance window\nionosctl dbaas mariadb cluster create --name prod-db --version 10.11 --instances 3 --cores 4 --ram 16GB --storage-size 100GB --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr 192.168.1.100/24 --user cluster_admin --password <password> --maintenance-day Sunday --maintenance-time 03:00:00", core.FlagsUsage(baseReqFlags...)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(baseReqFlags...)
		},
		CmdRun: func(c *core.CommandConfig) error {
			cluster := mariadb.CreateClusterProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				cluster.DisplayName = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.GetString(fn) != "" {
				cluster.Instances = viper.GetInt32(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
				cluster.MariadbVersion = mariadb.MariadbVersion(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.GetString(fn) != "" {
				cluster.Cores = viper.GetInt32(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.GetString(fn) != "" {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.GB)
				cluster.StorageSize = int32(sizeInt64)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagRam); viper.GetString(fn) != "" {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.GB)
				cluster.Ram = int32(sizeInt64)
			}

			cluster.MaintenanceWindow = &mariadb.MaintenanceWindow{}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.GetString(fn) != "" {
				cluster.MaintenanceWindow.DayOfTheWeek = mariadb.DayOfTheWeek(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.GetString(fn) != "" {
				cluster.MaintenanceWindow.Time = viper.GetString(fn)
			}

			cluster.Connections = make([]mariadb.Connection, 1)
			if fn := core.GetFlagName(c.NS, constants.FlagCidr); viper.IsSet(fn) {
				cluster.Connections[0].Cidr = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDatacenterId); viper.IsSet(fn) {
				cluster.Connections[0].DatacenterId = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLanId); viper.IsSet(fn) {
				cluster.Connections[0].LanId = viper.GetString(fn)
			}

			cluster.Credentials = mariadb.DBUser{}
			if fn := core.GetFlagName(c.NS, constants.ArgUser); viper.IsSet(fn) {
				cluster.Credentials.Username = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) {
				cluster.Credentials.Password = viper.GetString(fn)
			}

			createdCluster, _, err := client.Must().MariaClient.ClustersApi.ClustersPost(context.Background()).CreateClusterRequest(
				mariadb.CreateClusterRequest{Properties: &cluster},
			).Execute()
			if err != nil {
				return fmt.Errorf("failed creating cluster: %w", err)
			}

			return c.Printer(allCols).Print(createdCluster)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Human-friendly display name for the cluster (max 63 characters). Not a DNS name and need not be unique", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "10.6", "The MariaDB server version to run. One of: 10.6, 10.11. Can later be upgraded (never downgraded) via 'cluster update'", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10.6", "10.11"}, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "Number of instances in the cluster: 1 (standalone, no high availability) or an odd number (3 or 5) for a primary + secondaries replica set. Must be odd so the replicas can elect a primary. Range 1-5")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagCores, "", 1, "Number of CPU cores per instance (minimum 1). Applies to every instance in the cluster")
	cmd.AddStringFlag(constants.FlagRam, "", "4GB", "Memory per instance, e.g. --ram 4GB. Minimum 4GB; must be a whole number of GB. The upper bound is set by your contract quota")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "12GB", "16GB", "32GB", "64GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", strconv.Itoa(cloudapiv6.DefaultVolumeSize), "Storage per instance, e.g. --storage-size 10 or --storage-size 10GB. Minimum 10GB, maximum 2000GB (2TB). Can later be increased (never shrunk) via 'cluster update'")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"10GB", "20GB", "50GB", "100GB", "1TB"}, cobra.ShellCompDirectiveNoFileComp
	})
	// Maintenance
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hour := 10 + r.Intn(7) // Random hour 10-16
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", fmt.Sprintf("%02d:00:00", hour),
		"Start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window during which IONOS may apply patches and minor upgrades. Pairs with --maintenance-day. "+
			"Defaults to a random time between 10:00-16:00")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", workingDaysOfWeek[rand.Intn(len(workingDaysOfWeek))],
		"Day of the week the weekly maintenance window opens (e.g. Monday). Pairs with --maintenance-time. "+
			"Defaults to a random weekday (Mon-Fri)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
	})
	// Connections
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "ID of the Virtual Data Center (VDC) hosting the private LAN the cluster connects to. Must be in the same location (region) as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "Numeric ID of the private LAN (inside --datacenter-id) the cluster attaches to. The cluster is reachable only over this LAN, not the public internet", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagCidr, "", "", "The cluster's IP and subnet within the LAN in CIDR notation, e.g. 192.168.1.100/24 (use a /24 network). The address must lie within the LAN's subnet. Unavailable ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCidr, completer.GetCidrCompletionFunc(cmd))
	// credentials / DBUser
	cmd.AddStringFlag(constants.ArgUser, "", "", "Username for the initial database user (1-16 chars, must start with a letter, letters/digits/underscores only). Reserved names such as mariadb, admin and standby are rejected", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.ArgPassword, "", "", "Password for the initial database user (10-63 characters). Write-only: the API never returns it afterwards, so record it securely", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
