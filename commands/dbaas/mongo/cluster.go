package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/services/dbaas-postgres"

	//"github.com/ionos-cloud/ionosctl/services/dbaas-mongo"
	"github.com/ionos-cloud/ionosctl/services/dbaas-mongo/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
)

func ClusterCmd() *core.Command {
	ctx := context.TODO()
	clusterCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "PostgreSQL Cluster Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres cluster` allow you to manage the PostgreSQL Clusters under your account.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List Mongo Clusters",
		LongDesc:  "Use this command to retrieve a list of Mongo Clusters provisioned under your account. You can filter the result based on Cluster Name using `--name` option.",
		Example:   "ionosctl dbaas mongo cluster list",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			c.Printer.Verbose("Getting Clusters...")
			clusters, _, err := c.CloudApiDbaasPgsqlServices.Clusters().List(viper.GetString(core.GetFlagName(c.NS, dbaaspg.ArgName)))
			if err != nil {
				return err
			}
			return c.Printer.Print(getClusterPrint(nil, c, getClusters(clusters)))
		},
		InitClient: true,
	})
	//list.AddStringFlag(dbaaspg.ArgName, dbaaspg.ArgNameShort, "", "Response filter to list only the PostgreSQL Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	list.AddStringSliceFlag(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	return clusterCmd
}

func getClusterCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultClusterCols
	}

	columnsMap := map[string]string{
		"ClusterId":           "ClusterId",
		"DisplayName":         "DisplayName",
		"BackupLocation":      "BackupLocation",
		"Location":            "Location",
		"PostgresVersion":     "PostgresVersion",
		"State":               "State",
		"Ram":                 "Ram",
		"Instances":           "Instances",
		"Cores":               "Cores",
		"StorageSize":         "StorageSize",
		"StorageType":         "StorageType",
		"DatacenterId":        "DatacenterId",
		"LanId":               "LanId",
		"Cidr":                "Cidr",
		"MaintenanceWindow":   "MaintenanceWindow",
		"SynchronizationMode": "SynchronizationMode",
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

func getClusters(clusters resources.ClusterList) []resources.ClusterResponse {
	c := make([]resources.ClusterResponse, 0)
	if data, ok := clusters.GetItemsOk(); ok && data != nil {
		for _, d := range *data {
			c = append(c, resources.ClusterResponse{ClusterResponse: d})
		}
	}
	return c
}

func getClustersKVMaps(clusters []resources.ClusterResponse) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		var clusterPrint ClusterPrint
		if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
			clusterPrint.ClusterId = *idOk
		}
		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			if displayNameOk, ok := propertiesOk.GetDisplayNameOk(); ok && displayNameOk != nil {
				clusterPrint.DisplayName = *displayNameOk
			}
			if locationOk, ok := propertiesOk.GetLocationOk(); ok && locationOk != nil {
				clusterPrint.Location = string(*locationOk)
			}
			if backupLocationOk, ok := propertiesOk.GetBackupLocationOk(); ok && backupLocationOk != nil {
				clusterPrint.BackupLocation = string(*backupLocationOk)
			}
			if vdcConnectionsOk, ok := propertiesOk.GetConnectionsOk(); ok && vdcConnectionsOk != nil {
				for _, vdcConnection := range *vdcConnectionsOk {
					if vdcIdOk, ok := vdcConnection.GetDatacenterIdOk(); ok && vdcIdOk != nil {
						clusterPrint.DatacenterId = *vdcIdOk
					}
					if lanIdOk, ok := vdcConnection.GetLanIdOk(); ok && lanIdOk != nil {
						clusterPrint.LanId = *lanIdOk
					}
					if ipAddressOk, ok := vdcConnection.GetCidrOk(); ok && ipAddressOk != nil {
						clusterPrint.Cidr = *ipAddressOk
					}
				}
			}
			if postgresVersionOk, ok := propertiesOk.GetPostgresVersionOk(); ok && postgresVersionOk != nil {
				clusterPrint.PostgresVersion = *postgresVersionOk
			}
			if replicasOk, ok := propertiesOk.GetInstancesOk(); ok && replicasOk != nil {
				clusterPrint.Instances = *replicasOk
			}
			if ramSizeOk, ok := propertiesOk.GetRamOk(); ok && ramSizeOk != nil {
				clusterPrint.Ram = fmt.Sprintf("%vMB", *ramSizeOk)
			}
			if cpuCoreCountOk, ok := propertiesOk.GetCoresOk(); ok && cpuCoreCountOk != nil {
				clusterPrint.Cores = *cpuCoreCountOk
			}
			if storageSizeOk, ok := propertiesOk.GetStorageSizeOk(); ok && storageSizeOk != nil {
				clusterPrint.StorageSize = fmt.Sprintf("%vGB", *storageSizeOk)
			}
			if storageTypeOk, ok := propertiesOk.GetStorageTypeOk(); ok && storageTypeOk != nil {
				clusterPrint.StorageType = string(*storageTypeOk)
			}
			if maintenanceWindowOk, ok := propertiesOk.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				var maintenanceWindow string
				if weekdayOk, ok := maintenanceWindowOk.GetDayOfTheWeekOk(); ok && weekdayOk != nil {
					maintenanceWindow = string(*weekdayOk)
				}
				if timeOk, ok := maintenanceWindowOk.GetTimeOk(); ok && timeOk != nil {
					maintenanceWindow = fmt.Sprintf("%s %s", maintenanceWindow, *timeOk)
				}
				clusterPrint.MaintenanceWindow = maintenanceWindow
			}
			if synchronizationModeOk, ok := propertiesOk.GetSynchronizationModeOk(); ok && synchronizationModeOk != nil {
				clusterPrint.SynchronizationMode = string(*synchronizationModeOk)
			}
		}
		if metadataOk, ok := cluster.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				clusterPrint.State = string(*stateOk)
			}
		}
		o := structs.Map(clusterPrint)
		out = append(out, o)
	}
	return out
}
