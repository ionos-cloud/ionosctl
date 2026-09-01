package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/psql/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCmd() *core.Command {
	ctx := context.TODO()
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "cluster",
			Aliases: []string{"c"},
			Short:   "Create and manage PostgreSQL clusters",
			Long: `Manage PostgreSQL clusters, the top-level DBaaS PostgreSQL resource.

A cluster is a group of PostgreSQL instances (one master plus n-1 read-standby replicas, 1-5 total) running a chosen PostgreSQL version in a single physical location. It is reachable only over a private LAN in one of your virtual datacenters (VDC), addressed by the CIDR you assign. Compute (cores, RAM), storage (size and type), the synchronization mode, and the maintenance window are all configured per cluster.

Provisioning is asynchronous: after ` + "`create`" + `, the cluster moves through BUSY before reaching AVAILABLE. Databases and non-initial users can only be added once the cluster is AVAILABLE. Backups are taken automatically and can be used to restore this cluster in place (` + "`restore`" + `) or to clone a new one (` + "`create --backup-id`" + `).`,
			TraverseChildren: true,
		},
	}

	clusterCmd.AddColsFlag(allClusterCols)

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List PostgreSQL clusters",
		LongDesc:   "Retrieve all PostgreSQL clusters provisioned under your account. Use `--name` to keep only clusters whose display name contains the given substring (case-insensitive).",
		Example:    listClusterExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunClusterList,
		InitClient: true,
	})
	list.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Keep only clusters whose display name contains this substring (case-insensitive). Omit to list all clusters")
	/*
		Get Command
	*/
	get := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a PostgreSQL cluster",
		Example:    getClusterExample,
		LongDesc:   "Retrieve full details of a single PostgreSQL cluster by its ID, including its connection, sizing, synchronization mode, maintenance window and current lifecycle state.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterGet,
		InitClient: true,
	})
	get.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a PostgreSQL cluster",
		LongDesc: `Create a new PostgreSQL cluster.

The cluster attaches to a private LAN inside one of your virtual datacenters (VDC), so at minimum you must supply the network connection (--datacenter-id, --lan-id, --cidr) and the credentials for the initial database user (--db-username, --db-password). Every other property has a sensible default (PostgreSQL 15, 1 instance, 2 cores, 4GB RAM, 20GB HDD storage, ASYNCHRONOUS replication).

--location is the physical region the instances live in and cannot be changed later; if omitted it is inherited from the datacenter's location.

Sizing constraints: --instances 1-5 (1 master + n-1 read-standbys), --cores min 1, --ram min 4GB (must be a multiple of 1024MB), --storage-size 10GB-2TB.

To seed the cluster from an existing backup instead of an empty database (a clone), add --backup-id, and optionally --recovery-time to replay to a point in time within that backup's recovery window.

Provisioning is asynchronous: the cluster returns immediately in a non-AVAILABLE state. Wait for AVAILABLE (e.g. via 'cluster get') before creating additional databases or users.

Required values to run command:

* Datacenter Id
* Lan Id
* CIDR (IP and subnet)
* Credentials for the database user: Username and Password`,
		Example:    createClusterExample,
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(constants.FlagVersion, constants.FlagVersionShortPsql, "15", "The major PostgreSQL engine version to run (e.g. 13, 14, 15, 16). See 'dbaas postgres version list' for the versions currently offered")
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagVersion, "it will be removed in a future release.")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 1, "Total number of PostgreSQL instances: 1 master plus n-1 read-standby replicas. 1 means a single standalone instance (no HA). Minimum: 1, Maximum: 5")
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagInstances, "it will be removed in a future release.")
	create.AddIntFlag(constants.FlagCores, "", 2, "Number of CPU cores allocated to each instance. Minimum: 1")
	create.AddStringFlag(constants.FlagRam, "", "4GB", "Memory per instance. Must be a multiple of 1024MB and at least 4GB. Default unit is MB if none is given. e.g. --ram 4096, --ram 4096MB, --ram 4GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagBackupLocation, constants.FlagBackupLocationShortPsql, "", "Object Storage (S3) region where automated backups are stored: de, eu-south-2, eu-central-2, eu-central-3, eu-central-4, us-central-1. Defaults to a region derived from the cluster location. Cannot be changed after creation")
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagBackupLocation, "it will be removed in a future release.")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"de", "eu-south-2", "eu-central-2", "eu-central-3", "eu-central-4", "us-central-1"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagSyncMode, constants.FlagSyncModeShort, "ASYNCHRONOUS", "Replication mode between master and standbys. ASYNCHRONOUS: fastest, standbys may lag and a failover can lose the last transactions. STRICTLY_SYNCHRONOUS: a write is only acknowledged once a standby has it, safest but slower. SYNCHRONOUS is deprecated; prefer ASYNCHRONOUS or STRICTLY_SYNCHRONOUS")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagSyncMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageSize, "", "20GB", "Disk storage per instance. Default unit is MB if none is given. Minimum 10GB, maximum 2TB. e.g.: --storage-size 20480, --storage-size 20480MB, --storage-size 20GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageType, "", "HDD", "Disk performance tier: HDD (spinning disk, cheapest), SSD_STANDARD (general-purpose SSD), SSD_PREMIUM (highest IOPS). SSD is deprecated and treated as SSD_PREMIUM. Cannot be changed after creation")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD_PREMIUM", "SSD_STANDARD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLocation, "", "", "Physical region where the instances are provisioned (e.g. de/fra, de/txl, gb/lhr, us/las). Cannot be modified after creation. If omitted, the datacenter's location is used, so it must match the datacenter's region")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagName, constants.FlagNameShort, "UnnamedCluster", "Human-friendly display name for the cluster (does not have to be unique)")
	create.AddUUIDFlag(constants.FlagDatacenterId, "", "", "ID of the virtual datacenter (VDC) hosting the LAN the cluster attaches to. Its region also sets the default --location", core.RequiredFlagOption())
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagDatacenterId, "it will be removed in a future release.")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLanId, constants.FlagLanIdShortPsql, "", "Numeric ID of the LAN, within the chosen datacenter, that the cluster connects to. The cluster is reachable only from this private LAN", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(create.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCidr, constants.FlagCidrShortPsql, "", "IP address and subnet the master reserves on the LAN, in CIDR notation (e.g. 192.168.1.100/24). Must not overlap the reserved ranges 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagBackupId, constants.FlagBackupIdShortPsql, "", "Seed the new cluster from this backup instead of an empty database (clone). The backup's PostgreSQL version must be compatible with --version. See 'dbaas postgres backup list'")
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagBackupId, "it will be removed in a future release.")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "Only with --backup-id: replay the backup up to this point in time (ISO 8601 / RFC3339, e.g. 2024-01-15T10:00:00Z) for point-in-time recovery. Must fall within the backup's recovery window. If empty, the backup is applied in full (latest available point)")
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagRecoveryTime, "it will be removed in a future release.")

	create.AddStringFlag(constants.FlagDbUsername, constants.FlagDbUsernameShortPsql, "", "Username of the initial PostgreSQL superuser-equivalent role created with the cluster. Reserved system names such as postgres, admin and standby are not allowed", core.RequiredFlagOption())
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagDbUsername, "it will be removed in a future release.")

	create.AddStringFlag(constants.FlagDbPassword, constants.FlagDbPasswordShortPsql, "", "Password for the initial user. Minimum 10 characters", core.RequiredFlagOption())
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagDbPassword, "it will be removed in a future release.")

	create.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, "", "Start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window during which the service may apply updates. Set together with --maintenance-day. If omitted, a window is assigned automatically")
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagMaintenanceTime, "it will be removed in a future release.")
	create.AddStringFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, "", "Day of the week (e.g. Monday) for the weekly 4-hour maintenance window. Set together with --maintenance-time")
	create.Command.Flags().MarkShorthandDeprecated(constants.FlagMaintenanceDay, "it will be removed in a future release.")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a PostgreSQL cluster",
		LongDesc: `Update attributes of an existing PostgreSQL cluster. Only the flags you pass are changed; everything else is left untouched (a PATCH).

You can scale compute (--cores, --ram) and grow storage (--storage-size; storage can only be increased, not shrunk), rename the cluster (--name), upgrade the engine (--version), adjust the maintenance window (--maintenance-day and --maintenance-time must be given together), or change the network connection (--datacenter-id, --lan-id, --cidr).

--storage-type and --location cannot be changed after creation. Use --remove-connection to detach the cluster from its LAN entirely.

Required values to run command:

* Cluster Id`,
		Example:    updateClusterExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagVersion, constants.FlagVersionShortPsql, "", "Upgrade the PostgreSQL engine to this major version (e.g. 16). Only forward upgrades are supported; the cluster is briefly unavailable during the upgrade")
	update.Command.Flags().MarkShorthandDeprecated(constants.FlagVersion, "it will be removed in a future release.")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(constants.FlagRemoveConnection, "", false, "Detach the cluster from its LAN, removing the network connection entirely. Mutually exclusive with setting --datacenter-id/--lan-id/--cidr")
	update.AddUUIDFlag(constants.FlagDatacenterId, "", "", "Move the cluster's connection to this virtual datacenter. It must be in the same location as the current datacenter")
	update.Command.Flags().MarkShorthandDeprecated(constants.FlagDatacenterId, "it will be removed in a future release.")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagLanId, constants.FlagLanIdShortPsql, "", "Move the cluster's connection to this LAN (within the target datacenter)")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(update.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagCidr, constants.FlagCidrShortPsql, "", "New IP address and subnet for the cluster on the LAN, in CIDR notation (e.g. 192.168.1.100/24). Must not overlap the reserved ranges 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24")
	update.Command.Flags().MarkShorthandDeprecated(constants.FlagCidr, "it will be removed in a future release.")
	update.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 0, "New total number of instances (1 master + n-1 standbys). Maximum: 5. Leave at 0 to keep the current count")
	update.Command.Flags().MarkShorthandDeprecated(constants.FlagInstances, "it will be removed in a future release.")
	update.AddIntFlag(constants.FlagCores, "", 0, "New number of CPU cores per instance. Leave at 0 to keep the current value")
	update.AddStringFlag(constants.FlagRam, "", "", "New memory per instance. Must be a multiple of 1024MB and at least 4GB. Default unit is MB. e.g. --ram 4096, --ram 4096MB, --ram 4GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagStorageSize, "", "", "New storage per instance. Storage can only be increased, never decreased. Default unit is MB. e.g.: --storage-size 20480, --storage-size 20480MB, --storage-size 20GB")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "New human-friendly display name for the cluster")
	update.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, "", "New start time (UTC, HH:MM:SS, e.g. 16:30:59) of the weekly 4-hour maintenance window. Must be set together with --maintenance-day")
	update.Command.Flags().MarkShorthandDeprecated(constants.FlagMaintenanceTime, "it will be removed in a future release.")
	update.AddStringFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, "", "New day of the week (e.g. Monday) for the weekly 4-hour maintenance window. Must be set together with --maintenance-time")
	update.Command.Flags().MarkShorthandDeprecated(constants.FlagMaintenanceDay, "it will be removed in a future release.")
	_ = update.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Restore Command
	*/
	restoreCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
		Resource:  "cluster",
		Verb:      "restore",
		Aliases:   []string{"r"},
		ShortDesc: "Restore a PostgreSQL cluster in place from a backup",
		LongDesc: `Trigger an in-place restore of an existing PostgreSQL cluster from one of its backups. This overwrites the cluster's current data with the state captured in the backup, so it is a destructive operation on the live cluster.

A backup is not a single frozen moment but a continuous recovery WINDOW. By default the backup is replayed in full (the latest available point). Pass --recovery-time to roll back to an earlier point in time inside that window (point-in-time recovery); the timestamp must fall within the backup's recovery window (see 'dbaas postgres backup get', field EarliestRecoveryTargetTime).

To instead create a NEW cluster from a backup while leaving the current one intact, use 'cluster create --backup-id'.

Required values to run command:

* Cluster Id
* Backup Id`,
		Example:    restoreClusterExample,
		PreCmdRun:  PreRunClusterBackupIds,
		CmdRun:     RunClusterRestore,
		InitClient: true,
	})
	restoreCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(constants.FlagBackupId, "", "", "ID of the backup to restore from. Completion lists only backups belonging to the chosen cluster", core.RequiredFlagOption())
	_ = restoreCmd.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIdsForCluster(viper.GetString(core.GetFlagName(restoreCmd.NS, constants.FlagClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	restoreCmd.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "Replay the backup up to this point in time (ISO 8601 / RFC3339, e.g. 2024-01-15T10:00:00Z) for point-in-time recovery. Must fall within the backup's recovery window. If empty, the backup is applied in full")
	restoreCmd.Command.Flags().MarkShorthandDeprecated(constants.FlagRecoveryTime, "it will be removed in a future release.")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-postgres",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a PostgreSQL cluster",
		LongDesc: `Delete a PostgreSQL cluster and all of its data. This is irreversible; the cluster's databases and users are destroyed. Any retained backups are governed by the service's backup retention, not by this command.

Delete a single cluster with --cluster-id, or delete many at once with --all (optionally narrowed by --name substring match). Use ` + "`--wait` (`-w`)" + ` to block until deletion completes.

Required values to run command:

* Cluster Id`,
		Example:    deleteClusterExample,
		PreCmdRun:  PreRunClusterDelete,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all clusters in the account (subject to --name filtering). Prompts for confirmation per cluster unless --force is set")
	deleteCmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Restrict --all to clusters whose display name contains this substring (case-insensitive, not an exact match). Only valid together with --all")

	clusterCmd.AddCommand(ClusterBackupCmd())

	return clusterCmd
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func PreRunClusterDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagClusterId}, []string{constants.ArgAll})
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		if !viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) {
			return errors.New("error: name flag can to be used with the --all flag")
		}
	}
	return nil
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr, constants.FlagDbUsername, constants.FlagDbPassword)
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores)) < 1 {
		return errors.New("cores must be set to minimum: 1")
	}
	if viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances)) < 1 || viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances)) > 5 {
		return errors.New("instances must be set to minimum: 1, maximum: 5")
	}
	return nil
}

func PreRunClusterBackupIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagBackupId)
}

func RunClusterList(c *core.CommandConfig) error {
	c.Verbose("Getting Clusters...")

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		c.Verbose("Filtering after Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	}

	req := client.Must().PostgresClient.ClustersApi.ClustersGet(context.Background())

	nameFilter := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	if nameFilter != "" {
		req.FilterName(nameFilter)
	}

	clusters, _, err := req.Execute()
	if err != nil {
		return fmt.Errorf("could not list clusters: %w", err)
	}

	return c.Printer(allClusterCols).Prefix("items").Print(clusters)
}

func RunClusterGet(c *core.CommandConfig) error {
	c.Verbose(constants.ClusterId, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	c.Verbose("Getting Cluster...")

	cluster, _, err := client.Must().PostgresClient.ClustersApi.ClustersFindById(
		context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))).Execute()
	if err != nil {
		return fmt.Errorf("could not get cluster: %w", err)
	}

	return c.Printer(allClusterCols).Print(cluster)
}

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}

	c.Verbose("Creating Cluster...")

	cluster, _, err := client.Must().PostgresClient.ClustersApi.ClustersPost(context.Background()).
		CreateClusterRequest(input).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allClusterCols).Print(cluster)
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	c.Verbose(constants.ClusterId, clusterId)

	input, err := getPatchClusterRequest(c)
	if err != nil {
		return err
	}

	c.Verbose("Updating Cluster...")

	item, _, err := client.Must().PostgresClient.ClustersApi.
		ClustersPatch(context.Background(), clusterId).
		PatchClusterRequest(input).
		Execute()
	if err != nil {
		return err
	}

	return c.Printer(allClusterCols).Print(item)
}

func RunClusterRestore(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))

	c.Verbose(constants.ClusterId, clusterId)
	c.Verbose("Backup ID: %v", backupId)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("restore cluster with id: %v from backup: %v", clusterId, backupId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	input := psql.CreateRestoreRequest{
		BackupId: backupId,
	}

	if viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) != "" {
		c.Verbose("Setting RecoveryTargetTime [RFC3339 format]: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))

		recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
		if err != nil {
			return err
		}

		input.SetRecoveryTargetTime(recoveryTargetTime)
	}

	c.Verbose("Restoring Cluster from Backup...")

	_, err := client.Must().PostgresClient.RestoresApi.ClusterRestorePost(context.Background(), clusterId).
		CreateRestoreRequest(input).Execute()

	if err != nil {
		return err
	}
	c.Msg("PostgreSQL Cluster successfully restored")
	return nil
}

func RunClusterDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
		if err := ClusterDeleteAll(c); err != nil {
			return err
		}
		return nil
	}

	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	c.Verbose(constants.ClusterId, clusterId)

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete cluster with id: %v", clusterId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	_, _, err := client.Must().PostgresClient.ClustersApi.ClustersDelete(context.Background(), clusterId).Execute()
	if err != nil {
		return err
	}
	return nil
}

func ClusterDeleteAll(c *core.CommandConfig) error {

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		c.Verbose("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	}

	req := client.Must().PostgresClient.ClustersApi.ClustersGet(context.Background())
	if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) && viper.GetString(fn) != "" {
		req = req.FilterName(viper.GetString(fn))
	}
	clusters, _, err := req.Execute()
	if err != nil {
		return err
	}

	dataOk, ok := clusters.GetItemsOk()
	if !ok || dataOk == nil {
		return fmt.Errorf("could not get items of Clusters")
	}

	if len(dataOk) <= 0 {
		return fmt.Errorf("no Clusters found")
	}

	var multiErr error
	for _, cluster := range dataOk {
		idOk, ok := cluster.GetIdOk()
		if !ok || idOk == nil {
			continue
		}

		clusterId := *idOk
		propertiesOk, ok := cluster.GetPropertiesOk()
		clusterName, ok := propertiesOk.GetDisplayNameOk()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete cluster with id: %v, name: %v", clusterId, *clusterName), viper.GetBool(constants.ArgForce)) {
			multiErr = errors.Join(multiErr, fmt.Errorf(confirm.UserDenied))
			continue
		}

		_, _, err := client.Must().PostgresClient.ClustersApi.ClustersDelete(context.Background(), *idOk).Execute()
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *idOk, err))
			continue
		}

	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

func getCreateClusterRequest(c *core.CommandConfig) (psql.CreateClusterRequest, error) {
	inputCluster := psql.CreateClusterRequest{}
	input := psql.CreateClusterProperties{}

	// Setting Attributes
	pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	c.Verbose("PostgresVersion: %v", pgsqlVersion)
	input.SetPostgresVersion(pgsqlVersion)

	syncMode := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, constants.FlagSyncMode)))
	c.Verbose("SynchronizationMode: %v", syncMode)
	input.SetSynchronizationMode(psql.SynchronizationMode(syncMode))

	replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
	c.Verbose("Instances: %v", replicas)
	input.SetInstances(replicas)

	cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
	c.Verbose("Cores: %v", cpuCoreCount)
	input.SetCores(cpuCoreCount)

	// Convert Ram
	size, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils2.MegaBytes)
	if err != nil {
		return inputCluster, err
	}
	input.SetRam(int32(size))
	c.Verbose("Ram: %v[MB]", int32(size))

	// Convert StorageSize
	storageSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils2.MegaBytes)
	if err != nil {
		return inputCluster, err
	}
	input.SetStorageSize(int32(storageSize))
	c.Verbose("StorageSize: %v[MB]", int32(storageSize))
	storageType := strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageType)))
	// "HDD" "SSD" "SSD Standard" "SSD Premium". "SSD" is deprecated and equivalent to "SSD Premium"
	if storageType == "SSD_PREMIUM" || storageType == "SSD PREMIUM" {
		storageType = "SSD Premium"
	}
	if storageType == "SSD_STANDARD" || storageType == "SSD STANDARD" {
		storageType = "SSD Standard"
	}
	c.Verbose("StorageType: %v", storageType)
	input.SetStorageType(psql.StorageType(storageType))

	if viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)) != "" {
		c.Verbose("Backup Location: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)))
		input.SetBackupLocation(viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLocation)) {
		location := viper.GetString(core.GetFlagName(c.NS, constants.FlagLocation))
		c.Verbose("Location: %v", location)
		input.SetLocation(location)
	} else {
		c.Verbose("Getting Location from VDC...")
		vdc, _, err := client.Must().CloudClient.DataCentersApi.
			DatacentersFindById(
				context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId)),
			).Execute()
		if err != nil {
			return inputCluster, err
		}

		if properties, ok := vdc.GetPropertiesOk(); ok && properties != nil {
			if location, ok := properties.GetLocationOk(); ok && location != nil {
				c.Verbose("Location: %v", *location)
				input.SetLocation(*location)
			}
		}
	}

	displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	c.Verbose("DisplayName: %v", displayName)
	input.SetDisplayName(displayName)

	dbuser := psql.DBUser{}
	username := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbUsername))
	c.Verbose("DBUser - Username: %v", username)
	dbuser.SetUsername(username)

	password := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbPassword))
	c.Verbose("DBUser - Password: %v", password)
	dbuser.SetPassword(password)

	input.SetCredentials(dbuser)

	vdcConnection := psql.Connection{}
	vdcId := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
	c.Verbose("Connection - DatacenterId: %v", vdcId)
	vdcConnection.SetDatacenterId(vdcId)

	lanId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
	c.Verbose("Connection - LanId: %v", lanId)
	vdcConnection.SetLanId(lanId)

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
		ip := viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
		c.Verbose("Connection - Cidr: %v", ip)
		vdcConnection.SetCidr(ip)
	}

	input.SetConnections([]psql.Connection{vdcConnection})

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) ||
		viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		maintenanceWindow := psql.MaintenanceWindow{}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
			c.Verbose("MaintenanceWindow - Time: %v", maintenanceTime)
			maintenanceWindow.SetTime(maintenanceTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			c.Verbose("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay)
			maintenanceWindow.SetDayOfTheWeek(psql.DayOfTheWeek(maintenanceDay))
		}

		input.SetMaintenanceWindow(maintenanceWindow)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) ||
		viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
		createRestoreRequest := psql.CreateRestoreRequest{}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
			recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
			if err != nil {
				return inputCluster, err
			}

			c.Verbose("From Backup - RecoveryTargetTime [RFC3339 format]: %v", recoveryTargetTime)
			createRestoreRequest.SetRecoveryTargetTime(recoveryTargetTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) {
			backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))
			c.Verbose("From Backup - BackupId: %v", backupId)
			createRestoreRequest.SetBackupId(backupId)
		}

		input.SetFromBackup(createRestoreRequest)
	}

	inputCluster.SetProperties(input)
	return inputCluster, nil
}

func getPatchClusterRequest(c *core.CommandConfig) (psql.PatchClusterRequest, error) {
	inputCluster := psql.PatchClusterRequest{}
	input := psql.PatchClusterProperties{}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCores)) {
		cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
		c.Verbose("Cores: %v", cpuCoreCount)
		input.SetCores(cpuCoreCount)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagInstances)) {
		replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
		c.Verbose("Instances: %v", replicas)
		input.SetInstances(replicas)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRam)) {
		// Convert Ram
		size, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils2.MegaBytes)
		if err != nil {
			return inputCluster, err
		}

		input.SetRam(int32(size))
		c.Verbose("Ram: %vMB", int32(size))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagStorageSize)) {
		// Convert StorageSize
		storageSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils2.MegaBytes)
		if err != nil {
			return inputCluster, err
		}

		input.SetStorageSize(int32(storageSize))
		c.Verbose("StorageSize: %vMB", storageSize)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagVersion)) {
		pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
		c.Verbose("PostgresVersion: %v", pgsqlVersion)
		input.SetPostgresVersion(pgsqlVersion)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagName)) {
		displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
		c.Verbose("DisplayName: %v", displayName)
		input.SetDisplayName(displayName)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		maintenanceWindow := psql.MaintenanceWindow{}
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
			c.Verbose("MaintenanceTime: %v", maintenanceTime)
			maintenanceWindow.SetTime(maintenanceTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			c.Verbose("MaintenanceDayOfWeek: %v", maintenanceDay)
			maintenanceWindow.SetDayOfTheWeek(psql.DayOfTheWeek(maintenanceDay))
		}

		input.SetMaintenanceWindow(maintenanceWindow)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatacenterId)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagLanId)) || viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
		connection, err := getConnectionFromCluster(c, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
		if err != nil {
			return inputCluster, err
		}

		c.Verbose(getConnectionMessage(connection))

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagDatacenterId)) {
			lanId := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
			c.Verbose("Updated Datacenter Id: %v", lanId)
			connection.SetDatacenterId(lanId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLanId)) {
			lanId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
			c.Verbose("Updated Lan Id: %v", lanId)
			connection.SetLanId(lanId)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
			cidrId := viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
			c.Verbose("Updated Cidr: %v", cidrId)
			connection.SetCidr(cidrId)
		}

		input.SetConnections([]psql.Connection{connection})
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.FlagRemoveConnection)) {
		connection, err := getConnectionFromCluster(c, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
		if err != nil {
			return inputCluster, err
		}

		c.Verbose("Removing Connection with: %v...", getConnectionMessage(connection))
		input.SetConnections([]psql.Connection{})

	}

	inputCluster.SetProperties(input)
	return inputCluster, nil
}

func getConnectionFromCluster(c *core.CommandConfig, clusterId string) (psql.Connection, error) {
	if c != nil {
		oldCluster, _, err := client.Must().PostgresClient.ClustersApi.
			ClustersFindById(context.Background(), clusterId).Execute()
		if err != nil {
			return psql.Connection{}, err
		}

		c.Verbose("Getting properties from cluster with Id: %v", clusterId)
		if propertiesOk, ok := oldCluster.GetPropertiesOk(); ok && propertiesOk != nil {
			c.Verbose("Getting connection..")

			if connectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && connectionsOk != nil {
				for _, connectionOk := range connectionsOk {
					return connectionOk, nil
				}
			} else {
				return psql.Connection{}, errors.New("no connections found")
			}
		}
	}

	return psql.Connection{}, nil
}

func getConnectionMessage(connection psql.Connection) string {
	var msg string

	if datacenterOk, ok := connection.GetDatacenterIdOk(); ok && datacenterOk != nil {
		msg = fmt.Sprintf("DatacenterId: %v", *datacenterOk)
	}

	if lanOk, ok := connection.GetLanIdOk(); ok && lanOk != nil {
		msg = fmt.Sprintf("%v, LanId: %v", msg, *lanOk)
	}

	if cidrOk, ok := connection.GetCidrOk(); ok && cidrOk != nil {
		msg = fmt.Sprintf("%v, Cidr: %v", msg, *cidrOk)
	}

	return msg
}

// Output Printing

var allClusterCols = []table.Column{
	{Name: "ClusterId", JSONPath: "id", Default: true},
	{Name: "DisplayName", JSONPath: "properties.displayName", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "DatacenterId", JSONPath: "properties.connections.0.datacenterId", Default: true},
	{Name: "LanId", JSONPath: "properties.connections.0.lanId", Default: true},
	{Name: "Cidr", JSONPath: "properties.connections.0.cidr", Default: true},
	{Name: "Instances", JSONPath: "properties.instances", Default: true},
	{Name: "State", JSONPath: "metadata.state", Default: true},
	{Name: "PostgresVersion", JSONPath: "properties.postgresVersion"},
	{Name: "RAM", Format: func(item map[string]any) any {
		v := table.Navigate(item, "properties.ram")
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return v
		}
		return fmt.Sprintf("%d GB", convbytes.Convert(int64(f), convbytes.MB, convbytes.GB))
	}},
	{Name: "Cores", JSONPath: "properties.cores"},
	{Name: "StorageSize", Format: func(item map[string]any) any {
		v := table.Navigate(item, "properties.storageSize")
		if v == nil {
			return nil
		}
		f, ok := v.(float64)
		if !ok {
			return v
		}
		return fmt.Sprintf("%d GB", convbytes.Convert(int64(f), convbytes.MB, convbytes.GB))
	}},
	{Name: "StorageType", JSONPath: "properties.storageType"},
	{Name: "MaintenanceWindow", Format: func(item map[string]any) any {
		day, _ := table.Navigate(item, "properties.maintenanceWindow.dayOfTheWeek").(string)
		t, _ := table.Navigate(item, "properties.maintenanceWindow.time").(string)
		if day == "" && t == "" {
			return nil
		}
		return fmt.Sprintf("%s %s", day, t)
	}},
	{Name: "SynchronizationMode", JSONPath: "properties.synchronizationMode"},
	{Name: "BackupLocation", JSONPath: "properties.backupLocation"},
}
