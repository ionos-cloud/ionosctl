package cluster

import (
	"context"
	"errors"
	"fmt"
	"time"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres-v2/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	psqlv2 "github.com/ionos-cloud/sdk-go-dbaas-psql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCreateCmd() *core.Command {
	ctx := context.TODO()
	create := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "dbaas-postgres-v2",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a PostgreSQL Cluster",
		LongDesc: `Use this command to create a new PostgreSQL Cluster. You must set the unique ID of the Datacenter, the unique ID of the LAN, and IP and subnet. If the other options are not set, the default values will be used. Regarding the location field, if it is not manually set, it will be used the location of the Datacenter.

Required values to run command:

* Datacenter Id
* Lan Id
* CIDR (IP and subnet)
* Credentials for the database user: Username and Password`,
		Example:    "ionosctl dbaas postgres cluster create --datacenter-id <datacenter-id> --lan-id <lan-id> --cidr <cidr> --db-username <username> --db-password <password>",
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(constants.FlagVersion, constants.FlagVersionShortPsql, "15", "The PostgreSQL version of your Cluster")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PostgresVersions(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddIntFlag(constants.FlagInstances, constants.FlagInstancesShortPsql, 1, "The number of instances in your cluster (one master and n-1 standbys). Minimum: 1. Maximum: 5")
	create.AddIntFlag(constants.FlagCores, "", 2, "The number of CPU cores per instance. Minimum: 1")
	create.AddStringFlag(constants.FlagRam, "", "4GB", "The amount of memory per instance. Size must be specified in multiples of 1024. The default unit is MB. Minimum: 4GB. e.g. --ram 4096, --ram 4096MB, --ram 4GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagBackupLocation, constants.FlagBackupLocationShortPsql, "", "The S3 location where the backups will be stored")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"de", "eu-south-2", "eu-central-2"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagSyncMode, constants.FlagSyncModeShort, "ASYNCHRONOUS", "Synchronization Mode. Represents different modes of replication")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagSyncMode, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"ASYNCHRONOUS", "SYNCHRONOUS", "STRICTLY_SYNCHRONOUS"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageSize, "", "20GB", "The amount of storage per instance. The default unit is MB. e.g.: --size 20480 or --size 20480MB or --size 20GB")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"2048MB", "10GB", "20GB", "50GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagStorageType, "", "HDD", "The storage type used in your cluster: HDD, SSD, SSD_PREMIUM, SSD_STANDARD. (Value \"SSD\" is deprecated. Use the equivalent \"SSD_PREMIUM\" instead)")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "SSD", "SSD_PREMIUM", "SSD_STANDARD"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLocation, "", "", "The physical location where the cluster will be created. It cannot be modified after datacenter creation. If not set, it will be used Datacenter's location")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagName, constants.FlagNameShort, "UnnamedCluster", "The friendly name of your cluster")
	create.AddUUIDFlag(constants.FlagDatacenterId, "", "", "The unique ID of the Datacenter to connect to your cluster", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagLanId, constants.FlagLanIdShortPsql, "", "The unique ID of the LAN to connect your cluster to", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(viper.GetString(core.GetFlagName(create.NS, constants.FlagDatacenterId))), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagCidr, constants.FlagCidrShortPsql, "", "The IP and subnet for the cluster. Note the following unavailable IP ranges: 10.233.64.0/18, 10.233.0.0/18, 10.233.114.0/24. e.g.: 192.168.1.100/24", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagBackupId, constants.FlagBackupIdShortPsql, "", "The unique ID of the backup you want to restore")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagBackupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BackupsIds(), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(constants.FlagRecoveryTime, constants.FlagRecoveryTimeShortPsql, "", "If this value is supplied as ISO 8601 timestamp, the backup will be replayed up until the given timestamp. If empty, the backup will be applied completely")

	create.AddStringFlag(constants.FlagDbUsername, constants.FlagDbUsernameShortPsql, "", "Username for the initial postgres user. Some system usernames are restricted (e.g. postgres, admin, standby)", core.RequiredFlagOption())
	create.AddStringFlag(constants.FlagDbPassword, constants.FlagDbPasswordShortPsql, "", "Password for the initial postgres user", core.RequiredFlagOption())

	create.AddStringFlag(constants.FlagMaintenanceTime, constants.FlagMaintenanceTimeShortPsql, "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59")
	create.AddStringFlag(constants.FlagMaintenanceDay, constants.FlagMaintenanceDayShortPsql, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur")
	_ = create.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(constants.ArgWaitForState, constants.ArgWaitForStateShort, constants.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	create.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	create.AddStringSliceFlag(constants.ArgCols, "", defaultClusterCols, tabheaders.ColsMessage(allClusterCols))
	_ = create.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	return create
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

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Creating Cluster..."))

	cluster, _, err := client.Must().PostgresClientV2.ClustersApi.ClustersPost(context.Background()).ClusterCreate(input).Execute()
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if id, ok := cluster.GetIdOk(); ok && id != nil {
			if err = waitfor.WaitForState(c, waiter.ClusterStateInterrogator, *id); err != nil {
				return err
			}

			if cluster, _, err = client.Must().PostgresClientV2.ClustersApi.
				ClustersFindById(context.Background(), *id).Execute(); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new Cluster Id")
		}
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasPostgresV2Cluster, cluster,
		tabheaders.GetHeaders(allClusterCols, defaultClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func getCreateClusterRequest(c *core.CommandConfig) (psqlv2.ClusterCreate, error) {
	inputCluster := psqlv2.ClusterCreate{}
	input := psqlv2.ClusterCreateProperties{}

	// Setting Attributes
	pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, constants.FlagVersion))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("PostgresVersion: %v", pgsqlVersion))
	input.SetVersion(pgsqlVersion)

	instanceConfig := psqlv2.InstanceConfiguration{}
	replicas := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagInstances))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Instances: %v", replicas))
	instanceConfig.Count = replicas

	cpuCoreCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Cores: %v", cpuCoreCount))
	instanceConfig.Cores = cpuCoreCount

	// Convert Ram
	size, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils2.MegaBytes)
	if err != nil {
		return inputCluster, err
	}
	instanceConfig.Ram = int32(size)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Ram: %v[MB]", int32(size)))

	// Convert StorageSize
	storageSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils2.MegaBytes)
	if err != nil {
		return inputCluster, err
	}
	instanceConfig.SetStorageSize(int32(storageSize))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("StorageSize: %v[MB]", int32(storageSize)))

	input.SetInstances(instanceConfig)

	if viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Backup Location: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation))))
		input.SetBackupLocation(viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupLocation)))
	}

	displayName := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DisplayName: %v", displayName))
	input.SetName(displayName)

	dbuser := psqlv2.PostgresUser{}
	username := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbUsername))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DBUser - Username: %v", username))
	dbuser.SetUsername(username)

	password := viper.GetString(core.GetFlagName(c.NS, constants.FlagDbPassword))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("DBUser - Password: %v", password))
	dbuser.SetPassword(password)

	input.SetCredentials(dbuser)

	vdcConnection := psqlv2.PostgresClusterConnection{}
	vdcId := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Connection - DatacenterId: %v", vdcId))
	vdcConnection.SetDatacenterId(vdcId)

	lanId := viper.GetString(core.GetFlagName(c.NS, constants.FlagLanId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Connection - LanId: %v", lanId))
	vdcConnection.SetLanId(lanId)

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCidr)) {
		ip := viper.GetString(core.GetFlagName(c.NS, constants.FlagCidr))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Connection - Cidr: %v", ip))
		vdcConnection.SetPrimaryInstanceAddress(ip)
	}

	input.SetConnection(vdcConnection)

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) ||
		viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
		maintenanceWindow := psqlv2.MaintenanceWindow{}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceTime))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceWindow - Time: %v", maintenanceTime))
			maintenanceWindow.SetTime(maintenanceTime)
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, constants.FlagMaintenanceDay))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay))
			maintenanceWindow.SetDayOfTheWeek(psqlv2.DayOfTheWeek(maintenanceDay))
		}

		input.SetMaintenanceWindow(maintenanceWindow)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) ||
		viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
		restoreFromBackup := psqlv2.NewPostgresClusterFromBackup()

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagRecoveryTime)) {
			recoveryTargetTime, err := time.Parse(time.RFC3339, viper.GetString(core.GetFlagName(c.NS, constants.FlagRecoveryTime)))
			if err != nil {
				return inputCluster, err
			}

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("From Backup - RecoveryTargetTime [RFC3339 format]: %v", recoveryTargetTime))
			targetTime := psqlv2.IonosTime{Time: recoveryTargetTime}
			restoreFromBackup.RecoveryTargetDatetime = &targetTime
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagBackupId)) {
			backupId := viper.GetString(core.GetFlagName(c.NS, constants.FlagBackupId))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("From Backup - BackupId: %v", backupId))
			restoreFromBackup.SourceBackupId = &backupId
		}

		input.RestoreFromBackup = restoreFromBackup
	}

	inputCluster.SetProperties(input)
	return inputCluster, nil
}
