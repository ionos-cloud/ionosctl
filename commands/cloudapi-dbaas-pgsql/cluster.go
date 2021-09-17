package cloudapi_dbaas_pgsql

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-dbaas-pgsql/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources"
	sdkgo "github.com/ionos-cloud/sdk-go-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ClusterCmd() *core.Command {
	ctx := context.TODO()
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas-pgsql cluster` allow you to create, list, get, update and delete PostgreSQL Clusters.",
			TraverseChildren: true,
		},
	}
	globalFlags := clusterCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(clusterCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = clusterCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "cluster",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Clusters",
		LongDesc:   "Use this command to retrieve a list of PostgreSQL Clusters provisioned under your account.",
		Example:    listClusterExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunClusterList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "cluster",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Cluster",
		Example:    getClusterExample,
		LongDesc:   "Use this command to retrieve details about a PostgreSQL Cluster by using its ID.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapidbaaspgsql.ArgClusterId, cloudapidbaaspgsql.ArgIdShort, "", cloudapidbaaspgsql.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapidbaaspgsql.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "cluster",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a PostgreSQL Cluster",
		LongDesc: `Use this command to create a new PostgreSQL Cluster. You can specify the name, description or location for the object.

Virtual Clusters are the foundation of the IONOS platform. VDCs act as logical containers for all other objects you will be creating, e.g. servers. You can provision as many Clusters as you want. Clusters have their own private network and are logically segmented from each other to create isolation.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.`,
		Example:    createClusterExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Cluster creation to be executed")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Cluster creation [seconds]")

	//	/*
	//		Update Command
	//	*/
	//	update := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
	//		Namespace: "cluster",
	//		Resource:  "cluster",
	//		Verb:      "update",
	//		Aliases:   []string{"u", "up"},
	//		ShortDesc: "Update a Cluster",
	//		LongDesc: `Use this command to change a Virtual Cluster's name, description.
	//
	//You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.
	//
	//Required values to run command:
	//
	//* Cluster Id`,
	//		Example:    updateClusterExample,
	//		PreCmdRun:  PreRunClusterId,
	//		CmdRun:     RunClusterUpdate,
	//		InitClient: true,
	//	})
	//	update.AddStringFlag(cloudapidbaaspgsql.ArgClusterId, cloudapidbaaspgsql.ArgIdShort, "", cloudapidbaaspgsql.ClusterId, core.RequiredFlagOption())
	//	_ = update.Command.RegisterFlagCompletionFunc(cloudapidbaaspgsql.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	//		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	//	})
	//	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Cluster update to be executed")
	//	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Cluster update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "cluster",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a PostgreSQL Cluster",
		LongDesc: `Use this command to delete a specified Virtual Cluster from your account. This will remove all objects within the VDC and remove the VDC object itself.

NOTE: This is a highly destructive operation which should be used with extreme caution!

Required values to run command:

* Cluster Id`,
		Example:    deleteClusterExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapidbaaspgsql.ArgClusterId, cloudapidbaaspgsql.ArgIdShort, "", cloudapidbaaspgsql.ClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapidbaaspgsql.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Cluster deletion")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Cluster deletion [seconds]")

	return clusterCmd
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapidbaaspgsql.ArgClusterId)
}

func RunClusterList(c *core.CommandConfig) error {
	clusters, _, err := c.CloudApiDbaasPgsqlServices.Clusters().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(nil, c, getClusters(clusters)))
}

func RunClusterGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Cluster with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgClusterId)))
	dc, _, err := c.CloudApiDbaasPgsqlServices.Clusters().Get(viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(nil, c, []resources.Cluster{*dc}))
}

func RunClusterCreate(c *core.CommandConfig) error {
	pgsqlVersion := viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgPostgresVersion))
	replicas := float32(viper.GetFloat64(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgReplicas)))
	cpuCoreCount := float32(viper.GetFloat64(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgCpuCoreCount)))
	ramSize := viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgRamSize))
	storageSize := viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgStorageSize))
	storageType := sdkgo.StorageType(viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgStorageType)))
	location := viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgLocation))
	displayName := viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgName))
	backupEnabled := viper.GetBool(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgBackupEnabled))
	input := resources.CreateClusterRequest{
		CreateClusterRequest: sdkgo.CreateClusterRequest{
			PostgresVersion:   &pgsqlVersion,
			Replicas:          &replicas,
			CpuCoreCount:      &cpuCoreCount,
			RamSize:           &ramSize,
			StorageSize:       &storageSize,
			StorageType:       &storageType,
			VdcConnections:    nil,
			Location:          &location,
			DisplayName:       &displayName,
			BackupEnabled:     &backupEnabled,
			MaintenanceWindow: nil,
			Credentials:       nil,
		},
	}
	dc, resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Create(input, "", "")
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(resp, c, []resources.Cluster{*dc}))
}

func RunClusterDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete cluster"); err != nil {
		return err
	}
	id := viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgClusterId))
	c.Printer.Verbose("Deleting Cluster with ID: %v...", id)
	resp, err := c.CloudApiDbaasPgsqlServices.Clusters().Delete(id)
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(resp, c, nil))
}

// Output Printing

var (
	defaultClusterCols = []string{"ClusterId", "DisplayName", "Location", "BackupEnabled", "LifecycleStatus"}
	allClusterCols     = []string{"ClusterId", "DisplayName", "Location", "BackupEnabled", "LifecycleStatus", "PostgresVersion", "Replicas", "RamSize", "CpuCoreCount",
		"StorageSize", "StorageType", "VdcConnections", "MaintenanceWindow"}
)

type ClusterPrint struct {
	ClusterId         string  `json:"ClusterId,omitempty"`
	Location          string  `json:"Location,omitempty"`
	BackupEnabled     bool    `json:"BackupEnabled,omitempty"`
	LifecycleStatus   string  `json:"LifecycleStatus,omitempty"`
	DisplayName       string  `json:"DisplayName,omitempty"`
	PostgresVersion   string  `json:"PostgresVersion,omitempty"`
	Replicas          float32 `json:"Replicas,omitempty"`
	RamSize           string  `json:"RamSize,omitempty"`
	CpuCoreCount      float32 `json:"CpuCoreCount,omitempty"`
	StorageSize       string  `json:"StorageSize,omitempty"`
	StorageType       string  `json:"StorageType,omitempty"`
	VdcConnections    string  `json:"VdcConnections,omitempty"`
	MaintenanceWindow string  `json:"MaintenanceWindow,omitempty"`
}

func getClusterPrint(resp *resources.Response, c *core.CommandConfig, dcs []resources.Cluster) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getClustersKVMaps(dcs)
			r.Columns = getClusterCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getClusterCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultClusterCols
	}

	columnsMap := map[string]string{
		"ClusterId":         "ClusterId",
		"DisplayName":       "DisplayName",
		"Location":          "Location",
		"PostgresVersion":   "PostgresVersion",
		"BackupEnabled":     "BackupEnabled",
		"LifecycleStatus":   "LifecycleStatus",
		"Replicas":          "Replicas",
		"CpuCoreCount":      "CpuCoreCount",
		"StorageSize":       "StorageSize",
		"StorageType":       "StorageType",
		"VdcConnections":    "VdcConnections",
		"MaintenanceWindow": "MaintenanceWindow",
	}
	var clusterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			clusterCols = append(clusterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return clusterCols
}

func getClusters(clusters resources.ClusterList) []resources.Cluster {
	c := make([]resources.Cluster, 0)
	if data, ok := clusters.GetDataOk(); ok && data != nil {
		for _, d := range *data {
			c = append(c, resources.Cluster{Cluster: d})
		}
	}
	return c
}

func getClustersKVMaps(clusters []resources.Cluster) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		var clusterPrint ClusterPrint
		if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
			clusterPrint.ClusterId = *idOk
		}
		if displayNameOk, ok := cluster.GetDisplayNameOk(); ok && displayNameOk != nil {
			clusterPrint.DisplayName = *displayNameOk
		}
		if locationOk, ok := cluster.GetLocationOk(); ok && locationOk != nil {
			clusterPrint.Location = *locationOk
		}
		if backupEnabledOk, ok := cluster.GetBackupEnabledOk(); ok && backupEnabledOk != nil {
			clusterPrint.BackupEnabled = *backupEnabledOk
		}
		if lifecycleStatusOk, ok := cluster.GetLifecycleStatusOk(); ok && lifecycleStatusOk != nil {
			clusterPrint.LifecycleStatus = *lifecycleStatusOk
		}
		if postgresVersionOk, ok := cluster.GetPostgresVersionOk(); ok && postgresVersionOk != nil {
			clusterPrint.PostgresVersion = *postgresVersionOk
		}
		if replicasOk, ok := cluster.GetReplicasOk(); ok && replicasOk != nil {
			clusterPrint.Replicas = *replicasOk
		}
		if ramSizeOk, ok := cluster.GetRamSizeOk(); ok && ramSizeOk != nil {
			clusterPrint.RamSize = *ramSizeOk
		}
		if cpuCoreCountOk, ok := cluster.GetCpuCoreCountOk(); ok && cpuCoreCountOk != nil {
			clusterPrint.CpuCoreCount = *cpuCoreCountOk
		}
		if storageSizeOk, ok := cluster.GetStorageSizeOk(); ok && storageSizeOk != nil {
			clusterPrint.StorageSize = *storageSizeOk
		}
		if storageTypeOk, ok := cluster.GetStorageTypeOk(); ok && storageTypeOk != nil {
			clusterPrint.StorageType = string(*storageTypeOk)
		}
		o := structs.Map(clusterPrint)
		out = append(out, o)
	}
	return out
}
