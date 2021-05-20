package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8s() *core.Command {
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "k8s",
			Short:            "Kubernetes Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s` + "`" + ` allow you to list, get, create, update, delete Kubernetes Clusters.`,
			TraverseChildren: true,
		},
	}
	k8sCmd.AddCommand(k8sVersion())
	k8sCmd.AddCommand(k8sCluster())
	k8sCmd.AddCommand(k8sKubeconfig())
	k8sCmd.AddCommand(k8sNodePool())
	k8sCmd.AddCommand(k8sNode())

	return k8sCmd
}

func k8sCluster() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Short:            "Kubernetes Cluster Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s cluster` + "`" + ` allow you to list, get, create, update, delete Kubernetes Clusters.`,
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgFormat, config.ArgFormatShort, defaultK8sClusterCols, "Set of fields to be printed on output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(k8sCmd.Name(), config.ArgFormat), globalFlags.Lookup(config.ArgFormat))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(config.ArgFormat, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allK8sClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "cluster",
		Verb:       "list",
		ShortDesc:  "List Kubernetes Clusters",
		LongDesc:   "Use this command to get a list of existing Kubernetes Clusters.",
		Example:    listK8sClustersExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunK8sClusterList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "cluster",
		Verb:       "get",
		ShortDesc:  "Get a Kubernetes Cluster",
		LongDesc:   "Use this command to retrieve details about a specific Kubernetes Cluster.You can wait for the Cluster to be in \"ACTIVE\" state using `--wait-for-state` flag together with `--timeout` option.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    getK8sClusterExample,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sClusterGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Cluster to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.K8sTimeoutSeconds, "Timeout option for waiting for Cluster to be in ACTIVE state [seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "create",
		ShortDesc: "Create a Kubernetes Cluster",
		LongDesc: `Use this command to create a new Managed Kubernetes Cluster. Regarding the name for the Kubernetes Cluster, the limit is 63 characters following the rule to begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. Regarding the Kubernetes Version for the Cluster, if not set via flag, it will be used the default one: ` + "`" + `ionosctl k8s version get` + "`" + `.

You can wait for the Cluster to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run a command:

* K8s Cluster Name`,
		Example:    createK8sClusterExample,
		PreCmdRun:  PreRunK8sClusterName,
		CmdRun:     RunK8sClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name for the K8s Cluster "+config.RequiredFlag)
	create.AddStringFlag(config.ArgK8sVersion, "", "", "The K8s version for the Cluster. If not set, it will be used the default one")
	create.AddBoolFlag(config.ArgPublic, "", true, "The indicator if the Cluster is public or private")
	create.AddStringFlag(config.ArgGatewayIp, "", "", "The IP address of the gateway used by the Cluster. This is mandatory when `public` is set to `false` and should not be provided otherwise")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Cluster creation to be executed")
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new Cluster to be in ACTIVE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.K8sTimeoutSeconds, "Timeout option for waiting for Cluster/Request [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "update",
		ShortDesc: "Update a Kubernetes Cluster",
		LongDesc: `Use this command to update the name, Kubernetes version, maintenance day and maintenance time of an existing Kubernetes Cluster.

You can wait for the Cluster to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run command:

* K8s Cluster Id`,
		Example:    updateK8sClusterExample,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sClusterUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgName, config.ArgNameShort, "", "The name for the K8s Cluster")
	update.AddStringFlag(config.ArgK8sVersion, "", "", "The K8s version for the Cluster")
	update.AddStringFlag(config.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Cluster to be in ACTIVE state after updating")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.K8sTimeoutSeconds, "Timeout option for waiting for Cluster to be in ACTIVE state after updating [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "delete",
		ShortDesc: "Delete a Kubernetes Cluster",
		LongDesc: `This command deletes a Kubernetes cluster. The cluster cannot contain any NodePools when deleting.

You can wait for Request for the Cluster deletion to be executed using ` + "`" + `--wait-for-request` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run command:

* K8s Cluster Id`,
		Example:    deleteK8sClusterExample,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Cluster deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.K8sTimeoutSeconds, "Timeout option for waiting for Request [seconds]")

	return k8sCmd
}

func PreRunK8sClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgK8sClusterId)
}

func PreRunK8sClusterName(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgName)
}

func RunK8sClusterList(c *core.CommandConfig) error {
	k8ss, _, err := c.K8s().ListClusters()
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sClusters(k8ss)))
}

func RunK8sClusterGet(c *core.CommandConfig) error {
	if err := utils.WaitForState(c, GetStateK8sCluster, viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId))); err != nil {
		return err
	}
	u, _, err := c.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sCluster(u)))
}

func RunK8sClusterCreate(c *core.CommandConfig) error {
	newCluster, err := getNewK8sCluster(c)
	if err != nil {
		return err
	}
	u, resp, err := c.K8s().CreateCluster(*newCluster)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := u.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, GetStateK8sCluster, *id); err != nil {
				return err
			}
			if u, _, err = c.K8s().GetCluster(*id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new K8s Cluster id")
		}
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, getK8sCluster(u)))
}

func RunK8sClusterUpdate(c *core.CommandConfig) error {
	oldCluster, _, err := c.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	newCluster := getK8sClusterInfo(oldCluster, c)
	k8sUpd, _, err := c.K8s().UpdateCluster(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)), newCluster)
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if err = utils.WaitForState(c, GetStateK8sCluster, viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId))); err != nil {
			return err
		}
		if k8sUpd, _, err = c.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId))); err != nil {
			return err
		}
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sCluster(k8sUpd)))
}

func RunK8sClusterDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s cluster"); err != nil {
		return err
	}
	resp, err := c.K8s().DeleteCluster(viper.GetString(core.GetFlagName(c.NS, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, nil))
}

// Wait for State

func GetStateK8sCluster(c *core.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.K8s().GetCluster(objId)
	if err != nil {
		return nil, err
	}
	if metadata, ok := obj.GetMetadataOk(); ok && metadata != nil {
		if state, ok := metadata.GetStateOk(); ok && state != nil {
			return state, nil
		}
	}
	return nil, nil
}

func getNewK8sCluster(c *core.CommandConfig) (*resources.K8sClusterForPost, error) {
	var (
		k8sversion string
		err        error
	)
	proper := resources.K8sClusterPropertiesForPost{}
	proper.SetName(viper.GetString(core.GetFlagName(c.NS, config.ArgName)))
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sVersion)) {
		k8sversion = viper.GetString(core.GetFlagName(c.NS, config.ArgK8sVersion))
	} else {
		if k8sversion, err = getK8sVersion(c); err != nil {
			return nil, err
		}
	}
	proper.SetK8sVersion(k8sversion)
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgPublic)) {
		proper.SetPublic(viper.GetBool(core.GetFlagName(c.NS, config.ArgPublic)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgGatewayIp)) {
		proper.SetGatewayIp(viper.GetString(core.GetFlagName(c.NS, config.ArgGatewayIp)))
	}
	return &resources.K8sClusterForPost{
		KubernetesClusterForPost: ionoscloud.KubernetesClusterForPost{
			Properties: &proper.KubernetesClusterPropertiesForPost,
		},
	}, nil
}

func getK8sClusterInfo(oldUser *resources.K8sCluster, c *core.CommandConfig) resources.K8sClusterForPut {
	propertiesUpdated := resources.K8sClusterPropertiesForPut{}
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgName)) {
			n := viper.GetString(core.GetFlagName(c.NS, config.ArgName))
			propertiesUpdated.SetName(n)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				propertiesUpdated.SetName(*name)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sVersion)) {
			v := viper.GetString(core.GetFlagName(c.NS, config.ArgK8sVersion))
			propertiesUpdated.SetK8sVersion(v)
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaintenanceTime)) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := getMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})
				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}
	}
	return resources.K8sClusterForPut{
		KubernetesClusterForPut: ionoscloud.KubernetesClusterForPut{
			Properties: &propertiesUpdated.KubernetesClusterPropertiesForPut,
		},
	}
}

// Output Printing

var defaultK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "Public", "State", "MaintenanceWindow"}

var allK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "State", "MaintenanceWindow", "AvailableUpgradeVersions", "ViableNodePoolVersions", "Public", "GatewayIp"}

type K8sClusterPrint struct {
	ClusterId                string   `json:"ClusterId,omitempty"`
	Name                     string   `json:"Name,omitempty"`
	K8sVersion               string   `json:"K8sVersion,omitempty"`
	AvailableUpgradeVersions []string `json:"AvailableUpgradeVersions,omitempty"`
	ViableNodePoolVersions   []string `json:"ViableNodePoolVersions,omitempty"`
	MaintenanceWindow        string   `json:"MaintenanceWindow,omitempty"`
	State                    string   `json:"State,omitempty"`
	GatewayIps               string   `json:"GatewayIps,omitempty"`
	Public                   bool     `json:"Public,omitempty"`
}

func getK8sClusterPrint(resp *resources.Response, c *core.CommandConfig, k8ss []resources.K8sCluster) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState))
		}
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sClustersKVMaps(k8ss)
			r.Columns = getK8sClusterCols(core.GetGlobalFlagName(c.Resource, config.ArgFormat), c.Printer.GetStderr())
		}
	}
	return r
}

func getK8sClusterCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var k8sCols []string
		columnsMap := map[string]string{
			"ClusterId":                "ClusterId",
			"Name":                     "Name",
			"K8sVersion":               "K8sVersion",
			"AvailableUpgradeVersions": "AvailableUpgradeVersions",
			"ViableNodePoolVersions":   "ViableNodePoolVersions",
			"MaintenanceWindow":        "MaintenanceWindow",
			"Public":                   "Public",
			"GatewayIps":               "GatewayIps",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				k8sCols = append(k8sCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return k8sCols
	} else {
		return defaultK8sClusterCols
	}
}

func getK8sClusters(k8ss resources.K8sClusters) []resources.K8sCluster {
	u := make([]resources.K8sCluster, 0)
	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.K8sCluster{KubernetesCluster: item})
		}
	}
	return u
}

func getK8sCluster(u *resources.K8sCluster) []resources.K8sCluster {
	k8ss := make([]resources.K8sCluster, 0)
	if u != nil {
		k8ss = append(k8ss, resources.K8sCluster{KubernetesCluster: u.KubernetesCluster})
	}
	return k8ss
}

func getK8sClustersKVMaps(us []resources.K8sCluster) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint K8sClusterPrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.ClusterId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				uPrint.Name = *name
			}
			if v, ok := properties.GetK8sVersionOk(); ok && v != nil {
				uPrint.K8sVersion = *v
			}
			if v, ok := properties.GetAvailableUpgradeVersionsOk(); ok && v != nil {
				uPrint.AvailableUpgradeVersions = *v
			}
			if v, ok := properties.GetViableNodePoolVersionsOk(); ok && v != nil {
				uPrint.ViableNodePoolVersions = *v
			}
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				if day, ok := maintenance.GetDayOfTheWeekOk(); ok && day != nil {
					uPrint.MaintenanceWindow = *day
				}
				if time, ok := maintenance.GetTimeOk(); ok && time != nil {
					uPrint.MaintenanceWindow = fmt.Sprintf("%s %s", uPrint.MaintenanceWindow, *time)
				}
			}
			if pub, ok := properties.GetPublicOk(); ok && pub != nil {
				uPrint.Public = *pub
			}
			if gatewayIps, ok := properties.GetGatewayIpOk(); ok && gatewayIps != nil {
				uPrint.GatewayIps = *gatewayIps
			}
		}
		if meta, ok := u.GetMetadataOk(); ok && meta != nil {
			if state, ok := meta.GetStateOk(); ok && state != nil {
				uPrint.State = *state
			}
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}

func getK8sClustersIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(clientSvc.Get(), context.TODO())
	k8ss, _, err := k8sSvc.ListClusters()
	clierror.CheckError(err, outErr)
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.KubernetesClusters.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				k8ssIds = append(k8ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return k8ssIds
}

func getMaintenanceInfo(c *core.CommandConfig, maintenance *resources.K8sMaintenanceWindow) resources.K8sMaintenanceWindow {
	var day, time string
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaintenanceDay)) {
		day = viper.GetString(core.GetFlagName(c.NS, config.ArgK8sMaintenanceDay))
	} else {
		if d, ok := maintenance.GetDayOfTheWeekOk(); ok && d != nil {
			day = *d
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgK8sMaintenanceTime)) {
		time = viper.GetString(core.GetFlagName(c.NS, config.ArgK8sMaintenanceTime))
	} else {
		if t, ok := maintenance.GetTimeOk(); ok && t != nil {
			time = *t
		}
	}
	return resources.K8sMaintenanceWindow{
		KubernetesMaintenanceWindow: ionoscloud.KubernetesMaintenanceWindow{
			DayOfTheWeek: &day,
			Time:         &time,
		},
	}
}
