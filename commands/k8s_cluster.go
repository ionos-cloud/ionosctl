package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8s() *builder.Command {
	k8sCmd := &builder.Command{
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

func k8sCluster() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Short:            "Kubernetes Cluster Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s cluster` + "`" + ` allow you to list, get, create, update, delete Kubernetes Clusters.`,
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultK8sClusterCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(k8sCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allK8sClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sClusterList, "list", "List Kubernetes Clusters",
		"Use this command to get a list of existing Kubernetes Clusters.", listK8sClustersExample, true)

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterId, RunK8sClusterGet, "get", "Get a Kubernetes Cluster",
		"Use this command to retrieve details about a specific Kubernetes Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		getK8sClusterExample, true)
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterName, RunK8sClusterCreate, "create", "Create a Kubernetes Cluster",
		`Use this command to create a new Managed Kubernetes Cluster. Regarding the name for the Kubernetes Cluster, the limit is 63 characters following the rule to begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between. 

Required values to run a command:

* K8s Cluster Name`, createK8sClusterExample, true)
	create.AddStringFlag(config.ArgK8sClusterName, "", "", "The name for the K8s Cluster "+config.RequiredFlag)
	create.AddStringFlag(config.ArgK8sClusterVersion, "", "1.19.8", "The K8s version for the Cluster")

	/*
		Update Command
	*/
	update := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterId, RunK8sClusterUpdate, "update", "Update a Kubernetes Cluster",
		`Use this command to update the name, Kubernetes version, maintenance day and maintenance time of an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id`, updateK8sClusterExample, true)
	update.AddStringFlag(config.ArgK8sClusterName, "", "", "The name for the K8s Cluster")
	update.AddStringFlag(config.ArgK8sClusterVersion, "", "", "The K8s version for the Cluster")
	update.AddStringFlag(config.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	update.AddStringFlag(config.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterId, RunK8sClusterDelete, "delete", "Delete a Kubernetes Cluster",
		`This command deletes a Kubernetes cluster. The cluster cannot contain any NodePools when deleting.

Required values to run command:

* K8s Cluster Id`, deleteK8sClusterExample, true)
	deleteCmd.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func PreRunK8sClusterId(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterId)
}

func PreRunK8sClusterName(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterName)
}

func RunK8sClusterList(c *builder.CommandConfig) error {
	k8ss, _, err := c.K8s().ListClusters()
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sClusters(k8ss)))
}

func RunK8sClusterGet(c *builder.CommandConfig) error {
	u, _, err := c.K8s().GetCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sCluster(u)))
}

func RunK8sClusterCreate(c *builder.CommandConfig) error {
	n := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterName))
	v := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterVersion))
	newCluster := resources.K8sCluster{
		KubernetesCluster: ionoscloud.KubernetesCluster{
			Properties: &ionoscloud.KubernetesClusterProperties{
				Name:       &n,
				K8sVersion: &v,
			},
		},
	}
	u, resp, err := c.K8s().CreateCluster(newCluster)
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, getK8sCluster(u)))
}

func RunK8sClusterUpdate(c *builder.CommandConfig) error {
	oldCluster, _, err := c.K8s().GetCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	newCluster := getK8sClusterInfo(oldCluster, c)
	k8sUpd, _, err := c.K8s().UpdateCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)), newCluster)
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sCluster(k8sUpd)))
}

func RunK8sClusterDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete K8s cluster"); err != nil {
		return err
	}
	resp, err := c.K8s().DeleteCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, nil))
}

func getK8sClusterInfo(oldUser *resources.K8sCluster, c *builder.CommandConfig) resources.K8sCluster {
	propertiesUpdated := resources.K8sClusterProperties{}
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterName)) {
			n := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterName))
			propertiesUpdated.SetName(n)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				propertiesUpdated.SetName(*name)
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterVersion)) {
			v := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterVersion))
			propertiesUpdated.SetK8sVersion(v)
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceDay)) ||
			viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceTime)) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := getMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})
				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}
	}
	return resources.K8sCluster{
		KubernetesCluster: ionoscloud.KubernetesCluster{
			Properties: &propertiesUpdated.KubernetesClusterProperties,
		},
	}
}

// Output Printing

var defaultK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "State", "MaintenanceWindow"}

var allK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "State", "MaintenanceWindow", "AvailableUpgradeVersions", "ViableNodePoolVersions"}

type K8sClusterPrint struct {
	ClusterId                string   `json:"ClusterId,omitempty"`
	Name                     string   `json:"Name,omitempty"`
	K8sVersion               string   `json:"K8sVersion,omitempty"`
	AvailableUpgradeVersions []string `json:"AvailableUpgradeVersions,omitempty"`
	ViableNodePoolVersions   []string `json:"ViableNodePoolVersions,omitempty"`
	MaintenanceWindow        string   `json:"MaintenanceWindow,omitempty"`
	State                    string   `json:"State,omitempty"`
}

func getK8sClusterPrint(resp *resources.Response, c *builder.CommandConfig, k8ss []resources.K8sCluster) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sClustersKVMaps(k8ss)
			r.Columns = getK8sClusterCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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

func getMaintenanceInfo(c *builder.CommandConfig, maintenance *resources.K8sMaintenanceWindow) resources.K8sMaintenanceWindow {
	var day, time string
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceDay)) {
		day = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceDay))
	} else {
		if d, ok := maintenance.GetDayOfTheWeekOk(); ok && d != nil {
			day = *d
		}
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceTime)) {
		time = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sMaintenanceTime))
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
