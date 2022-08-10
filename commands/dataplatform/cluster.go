package dataplatform

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/multierr"

	"github.com/fatih/structs"
	cloudapiv6completer "github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/dataplatform/completer"
	"github.com/ionos-cloud/ionosctl/commands/dataplatform/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	dp "github.com/ionos-cloud/ionosctl/services/dataplatform"
	"github.com/ionos-cloud/ionosctl/services/dataplatform/resources"
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
			Short:            "Data Platform Cluster Operations",
			Long:             "The sub-commands of `ionosctl dataplatform cluster` allow you to manage the Data Platform Clusters under your account.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dataplatform",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Data Platform Clusters",
		LongDesc:   "Use this command to retrieve a list of Data Platform Clusters provisioned under your account. You can filter the result based on Data Platform Name using `--name` option.",
		Example:    listClusterExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunClusterList,
		InitClient: true,
	})
	list.AddStringFlag(dp.ArgName, dp.ArgNameShort, "", "Response filter to list only the Data Platform Clusters that contain the specified name in the DisplayName field. The value is case insensitive")
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	list.AddStringSliceFlag(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace:  "dataplatform",
		Resource:   "cluster",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Data Platform Cluster",
		Example:    getClusterExample,
		LongDesc:   "Use this command to retrieve details about a Data Platform Cluster by using its ID.\n\nRequired values to run command:\n\n* Cluster Id",
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterGet,
		InitClient: true,
	})
	get.AddStringFlag(dp.ArgClusterId, dp.ArgIdShort, "", dp.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dp.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state [seconds]")
	get.AddStringSliceFlag(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Data Platform Cluster",
		LongDesc: `Use this command to create a new Data Platform Cluster. You must set the unique ID of the Datacenter, and the Name of the Cluster. If the other options are not set, the default values will be used. 

Required values to run command:

* Datacenter Id
* Name`,
		Example:    createClusterExample,
		PreCmdRun:  PreRunClusterCreate,
		CmdRun:     RunClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(dp.ArgDatacenterId, dp.ArgDatacenterIdShort, "", "The UUID of the virtual data center (VDC) the cluster is provisioned.", core.RequiredFlagOption())
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dp.ArgName, dp.ArgNameShort, "", "The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.", core.RequiredFlagOption())
	create.AddStringFlag(dp.ArgVersion, dp.ArgVersionShort, "1.1.0", "The Data Platform version of your Cluster")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataPlatformVersions(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	create.AddStringFlag(dp.ArgMaintenanceTime, dp.ArgMaintenanceTimeShort, "", "Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	create.AddStringFlag(dp.ArgMaintenanceDay, dp.ArgMaintenanceDayShort, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dp.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	create.AddStringSliceFlag(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = create.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
		ShortDesc: "Update a Data Platform Cluster",
		LongDesc: `Use this command to update attributes of a Data Platform Cluster.

Required values to run command:

* Cluster Id`,
		Example:    updateClusterExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunClusterUpdate,
		InitClient: true,
	})
	update.AddStringFlag(dp.ArgClusterId, dp.ArgIdShort, "", dp.ClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dp.ArgName, dp.ArgNameShort, "", "The name of your cluster. Must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.")
	update.AddStringFlag(dp.ArgVersion, dp.ArgVersionShort, "", "The Data Platform version of your cluster")
	_ = update.Command.RegisterFlagCompletionFunc(dp.ArgVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.DataPlatformVersions(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(dp.ArgMaintenanceTime, dp.ArgMaintenanceTimeShort, "", "Time at which the maintenance should start. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	update.AddStringFlag(dp.ArgMaintenanceDay, dp.ArgMaintenanceDayShort, "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is starting time of a weekly 4 hour-long window, during which maintenance might occur in hh:mm:ss format")
	_ = create.Command.RegisterFlagCompletionFunc(dp.ArgMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be in AVAILABLE state")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dp.DefaultClusterTimeout, "Timeout option for Cluster to be in AVAILABLE state[seconds]")
	update.AddStringSliceFlag(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, clusterCmd, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Data Platform Cluster",
		LongDesc: `Use this command to delete a specified Data Platform Cluster from your account. You can wait for the cluster to be deleted with the wait-for-deletion option.

Required values to run command:

* Cluster Id`,
		Example:    deleteClusterExample,
		PreCmdRun:  PreRunClusterDelete,
		CmdRun:     RunClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(dp.ArgClusterId, dp.ArgIdShort, "", dp.ClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(dp.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgAll, config.ArgAllShort, false, "Delete all Clusters")
	deleteCmd.AddStringFlag(dp.ArgName, dp.ArgNameShort, "", "Delete all Clusters after filtering based on name. It does not require an exact match. Can be used with --all flag")
	deleteCmd.AddBoolFlag(config.ArgWaitForDelete, config.ArgWaitForStateShort, config.DefaultWait, "Wait for Cluster to be completely removed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, dp.DefaultClusterTimeout, "Timeout option for Cluster to be completely removed[seconds]")
	deleteCmd.AddStringSliceFlag(config.ArgCols, "", defaultClusterCols, printer.ColsMessage(allClusterCols))
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	return clusterCmd
}

func PreRunClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, dp.ArgClusterId)
}

func PreRunClusterDelete(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{dp.ArgClusterId}, []string{config.ArgAll})
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgName)) {
		if !viper.IsSet(core.GetFlagName(c.NS, config.ArgAll)) {
			return errors.New("error: name flag can to be used with the --all flag")
		}
	}
	return nil
}

func PreRunClusterCreate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, dp.ArgDatacenterId, dp.ArgName)
	if err != nil {
		return err
	}
	// Validate Flags
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgVersion)) {
		if len(viper.GetString(core.GetFlagName(c.NS, dp.ArgVersion))) > 32 {
			return errors.New("version string has to have a maximum of 32 characters")
		}
	}
	if len(viper.GetString(core.GetFlagName(c.NS, dp.ArgName))) > 63 {
		return errors.New("name string has to have a maximum of 63 characters")
	}
	return nil
}

func RunClusterList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Clusters...")
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgName)) {
		c.Printer.Verbose("Filtering after Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, dp.ArgName)))
	}
	clusters, _, err := c.DataPlatformServices.Clusters().List(viper.GetString(core.GetFlagName(c.NS, dp.ArgName)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(nil, c, getClusters(clusters)))
}

func RunClusterGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Cluster ID: %v", viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)))
	c.Printer.Verbose("Getting Cluster...")
	if err := utils.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))); err != nil {
		return err
	}
	cluster, _, err := c.DataPlatformServices.Clusters().Get(viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getClusterPrint(nil, c, []resources.ClusterResponseData{cluster}))
}

func RunClusterCreate(c *core.CommandConfig) error {
	input, err := getCreateClusterRequest(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Creating Cluster...")
	cluster, resp, err := c.DataPlatformServices.Clusters().Create(*input)
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := cluster.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.ClusterStateInterrogator, *id); err != nil {
				return err
			}
			if cluster, _, err = c.DataPlatformServices.Clusters().Get(*id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new Cluster Id")
		}
	}
	return c.Printer.Print(getClusterPrint(resp, c, []resources.ClusterResponseData{cluster}))
}

func RunClusterUpdate(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))
	c.Printer.Verbose("Cluster ID: %v", clusterId)
	input, err := getPatchClusterRequest(c)
	if err != nil {
		return err
	}
	c.Printer.Verbose("Updating Cluster...")
	item, resp, err := c.DataPlatformServices.Clusters().Update(clusterId, *input)
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		c.Printer.Verbose("Wait 10 seconds before checking state...")
		// Sleeping 10 seconds to make sure the cluster is in BUSY state. This will be removed in future releases.
		time.Sleep(10 * time.Second)
		if err = utils.WaitForState(c, waiter.ClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))); err != nil {
			return err
		}
	}
	return c.Printer.Print(getClusterPrint(resp, c, []resources.ClusterResponseData{item}))
}

func RunClusterDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgAll)) {
		if err := ClusterDeleteAll(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		clusterId := viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))
		c.Printer.Verbose("Cluster ID: %v", clusterId)
		if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete cluster with id: %v", clusterId)); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting Cluster...")
		_, resp, err := c.DataPlatformServices.Clusters().Delete(clusterId)
		if err != nil {
			return err
		}
		if err = utils.WaitForDelete(c, waiter.ClusterDeleteInterrogator, viper.GetString(core.GetFlagName(c.NS, dp.ArgClusterId))); err != nil {
			return err
		}
		return c.Printer.Print(getClusterPrint(resp, c, nil))
	}
}

func ClusterDeleteAll(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting all Clusters...")
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgName)) {
		c.Printer.Verbose("Filtering based on Cluster Name: %v", viper.GetString(core.GetFlagName(c.NS, dp.ArgName)))
	}
	clusters, _, err := c.DataPlatformServices.Clusters().List(viper.GetString(core.GetFlagName(c.NS, dp.ArgName)))
	if err != nil {
		return err
	}
	if dataOk, ok := clusters.GetItemsOk(); ok && dataOk != nil {
		if len(*dataOk) > 0 {
			_ = c.Printer.Print("Clusters to be deleted:")
			for _, cluster := range *dataOk {
				var log string
				if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
					if nameOk, ok := propertiesOk.GetNameOk(); ok && nameOk != nil {
						log = fmt.Sprintf("Cluster Name: %s", *nameOk)
					}
				}
				if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
					log = fmt.Sprintf("%s; Cluster ID: %s", log, *idOk)
				}
				c.Printer.Print(log)
			}
			if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete ALL clusters"); err != nil {
				return err
			}
			c.Printer.Verbose("Deleting all the Clusters...")
			var multiErr error
			for _, cluster := range *dataOk {
				if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
					c.Printer.Verbose("Cluster ID: %v", *idOk)
					c.Printer.Verbose("Deleting Cluster...")
					_, _, err = c.DataPlatformServices.Clusters().Delete(*idOk)
					if err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *idOk, err))
						continue
					} else {
						_ = c.Printer.Print(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *idOk))
					}
					if err = utils.WaitForDelete(c, waiter.ClusterDeleteInterrogator, *idOk); err != nil {
						multiErr = multierr.Append(multiErr, fmt.Errorf(config.WaitDeleteAllAppendErr, c.Resource, *idOk, err))
						continue
					}
				}
			}
			if multiErr != nil {
				return multiErr
			}
			return nil
		} else {
			return errors.New("no Clusters found")
		}
	} else {
		return errors.New("could not get items of Clusters")
	}
}

func getCreateClusterRequest(c *core.CommandConfig) (*resources.CreateClusterRequest, error) {
	inputCluster := resources.CreateClusterRequest{}
	input := sdkgo.CreateClusterProperties{}
	// Setting Attributes
	datacenterId := viper.GetString(core.GetFlagName(c.NS, dp.ArgDatacenterId))
	c.Printer.Verbose("Datacenter ID: %v", datacenterId)
	input.SetDatacenterId(datacenterId)

	name := viper.GetString(core.GetFlagName(c.NS, dp.ArgName))
	c.Printer.Verbose("Name: %v", name)
	input.SetName(name)

	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgVersion)) {
		version := viper.GetString(core.GetFlagName(c.NS, dp.ArgVersion))
		c.Printer.Verbose("Data Platform Version: %v", version)
		input.SetDataPlatformVersion(version)
	}

	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) ||
		viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) {
		maintenanceWindow := sdkgo.MaintenanceWindow{}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceTime))
			c.Printer.Verbose("MaintenanceWindow - Time: %v", maintenanceTime)
			maintenanceWindow.SetTime(maintenanceTime)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceDay))
			c.Printer.Verbose("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay)
			maintenanceWindow.SetDayOfTheWeek(maintenanceDay)
		}
		input.SetMaintenanceWindow(maintenanceWindow)
	}

	inputCluster.SetProperties(input)
	return &inputCluster, nil
}

func getPatchClusterRequest(c *core.CommandConfig) (*resources.PatchClusterRequest, error) {
	inputCluster := resources.PatchClusterRequest{}
	input := sdkgo.PatchClusterProperties{}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgName)) {
		name := viper.GetString(core.GetFlagName(c.NS, dp.ArgName))
		c.Printer.Verbose("Name: %v", name)
		input.SetName(name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgVersion)) {
		version := viper.GetString(core.GetFlagName(c.NS, dp.ArgVersion))
		c.Printer.Verbose("Data Platform Version: %v", version)
		input.SetDataPlatformVersion(version)
	}
	if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) ||
		viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) {
		maintenanceWindow := sdkgo.MaintenanceWindow{}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceTime)) {
			maintenanceTime := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceTime))
			c.Printer.Verbose("MaintenanceWindow - Time: %v", maintenanceTime)
			maintenanceWindow.SetTime(maintenanceTime)
		}
		if viper.IsSet(core.GetFlagName(c.NS, dp.ArgMaintenanceDay)) {
			maintenanceDay := viper.GetString(core.GetFlagName(c.NS, dp.ArgMaintenanceDay))
			c.Printer.Verbose("MaintenanceWindow - DayOfTheWeek: %v", maintenanceDay)
			maintenanceWindow.SetDayOfTheWeek(maintenanceDay)
		}
		input.SetMaintenanceWindow(maintenanceWindow)
	}
	inputCluster.SetProperties(input)
	return &inputCluster, nil
}

// Output Printing

var (
	defaultClusterCols = []string{"ClusterId", "Name", "DataPlatformVersion", "MaintenanceWindow", "DatacenterId", "State"}
	allClusterCols     = []string{"ClusterId", "Name", "DataPlatformVersion", "MaintenanceWindow", "DatacenterId", "State"}
)

type ClusterPrint struct {
	ClusterId           string `json:"ClusterId,omitempty"`
	Name                string `json:"Name,omitempty"`
	DataPlatformVersion string `json:"DataPlatformVersion,omitempty"`
	MaintenanceWindow   string `json:"MaintenanceWindow,omitempty"`
	DatacenterId        string `json:"DatacenterId,omitempty"`
	State               string `json:"State,omitempty"`
}

func getClusterPrint(resp *resources.Response, c *core.CommandConfig, dcs []resources.ClusterResponseData) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getClustersKVMaps(dcs)
			r.Columns = getClusterCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
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
		"ClusterId":           "ClusterId",
		"Name":                "Name",
		"DataPlatformVersion": "DataPlatformVersion",
		"MaintenanceWindow":   "MaintenanceWindow",
		"DatacenterId":        "DatacenterId",
		"State":               "State",
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

func getClusters(clusters resources.ClusterListResponseData) []resources.ClusterResponseData {
	c := make([]resources.ClusterResponseData, 0)
	if data, ok := clusters.GetItemsOk(); ok && data != nil {
		for _, d := range *data {
			c = append(c, resources.ClusterResponseData{ClusterResponseData: d})
		}
	}
	return c
}

func getClustersKVMaps(clusters []resources.ClusterResponseData) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(clusters))
	for _, cluster := range clusters {
		var clusterPrint ClusterPrint
		if idOk, ok := cluster.GetIdOk(); ok && idOk != nil {
			clusterPrint.ClusterId = *idOk
		}
		if propertiesOk, ok := cluster.GetPropertiesOk(); ok && propertiesOk != nil {
			if datacenterIdOk, ok := propertiesOk.GetDatacenterIdOk(); ok && datacenterIdOk != nil {
				clusterPrint.DatacenterId = *datacenterIdOk
			}

			if nameOk, ok := propertiesOk.GetNameOk(); ok && nameOk != nil {
				clusterPrint.Name = *nameOk
			}

			if versionOk, ok := propertiesOk.GetDataPlatformVersionOk(); ok && versionOk != nil {
				clusterPrint.DataPlatformVersion = *versionOk
			}

			if maintenanceWindowOk, ok := propertiesOk.GetMaintenanceWindowOk(); ok && maintenanceWindowOk != nil {
				var maintenanceWindow string
				if weekdayOk, ok := maintenanceWindowOk.GetDayOfTheWeekOk(); ok && weekdayOk != nil {
					maintenanceWindow = *weekdayOk
				}
				if timeOk, ok := maintenanceWindowOk.GetTimeOk(); ok && timeOk != nil {
					maintenanceWindow = fmt.Sprintf("%s %s", maintenanceWindow, *timeOk)
				}
				clusterPrint.MaintenanceWindow = maintenanceWindow
			}

		}
		if metadataOk, ok := cluster.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				clusterPrint.State = *stateOk
			}
		}
		o := structs.Map(clusterPrint)
		out = append(out, o)
	}
	return out
}
