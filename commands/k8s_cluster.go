package commands

import (
	"context"
	"errors"
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

func k8sCluster() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "k8s-cluster",
			Short:            "K8s Cluster Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s-cluster` + "`" + ` allow you to list, get, create, update, delete K8s Clusters.`,
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultK8sClusterCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(k8sCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	builder.NewCommand(ctx, k8sCmd, noPreRun, RunK8sClusterList, "list", "List K8s Clusters",
		"Use this command to get a list of existing K8s Clusters.", listK8sClustersExample, true)

	get := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterIdValidate, RunK8sClusterGet, "get", "Get a K8s Cluster",
		"Use this command to retrieve details about a specific K8s Cluster.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		getK8sClusterExample, true)
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	create := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNameValidate, RunK8sClusterCreate, "create", "Create a K8s Cluster",
		`Use this command to create a new Managed Kubernetes Cluster.

Required values to run a command:

* K8s Cluster Name`, createK8sClusterExample, true)
	create.AddStringFlag(config.ArgK8sClusterName, "", "", "The name for the K8s Cluster "+config.RequiredFlag)
	create.AddStringFlag(config.ArgK8sClusterVersion, "", "1.19.8", "The K8s version for the Cluster")

	update := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterIdValidate, RunK8sClusterUpdate, "update", "Update a K8s Cluster",
		`Use this command to update the name, version and other properties of an existing Kubernetes Cluster.

Required values to run command:

* K8s Cluster Id`, updateK8sClusterExample, true)
	update.AddStringFlag(config.ArgK8sClusterName, "", "", "The name for the K8s Cluster")
	update.AddStringFlag(config.ArgK8sClusterVersion, "", "", "The K8s version for the Cluster")
	update.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	deleteCmd := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterIdValidate, RunK8sClusterDelete, "delete", "Delete a K8s Cluster",
		`This command deletes a Kubernetes cluster. The cluster cannot contain any node pools when deleting.

Required values to run command:

* K8s Cluster Id`, deleteK8sClusterExample, true)
	deleteCmd.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func PreRunK8sClusterIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterId)
}

func PreRunK8sClusterNameValidate(c *builder.PreCommandConfig) error {
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
	newProperties := getK8sClusterInfo(oldCluster, c)
	newCluster := resources.K8sCluster{
		KubernetesCluster: ionoscloud.KubernetesCluster{
			Properties: &newProperties.KubernetesClusterProperties,
		},
	}
	k8sUpd, _, err := c.K8s().UpdateCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)), newCluster)
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sCluster(k8sUpd)))
}

func RunK8sClusterDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete K8s cluster")
	if err != nil {
		return err
	}
	resp, err := c.K8s().DeleteCluster(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, nil))
}

func getK8sClusterInfo(oldUser *resources.K8sCluster, c *builder.CommandConfig) *resources.K8sClusterProperties {
	var n, v string
	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterName)) {
			n = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterName))
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				n = *name
			}
		}
		if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterVersion)) {
			v = viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterVersion))
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				v = *vers
			}
		}
	}
	return &resources.K8sClusterProperties{
		KubernetesClusterProperties: ionoscloud.KubernetesClusterProperties{
			Name:       &n,
			K8sVersion: &v,
		},
	}
}

// Output Printing

var defaultK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "AvailableUpgradeVersions", "ViableNodePoolVersions", "State"}

type K8sClusterPrint struct {
	ClusterId                string   `json:"ClusterId,omitempty"`
	Name                     string   `json:"Name,omitempty"`
	K8sVersion               string   `json:"K8sVersion,omitempty"`
	AvailableUpgradeVersions []string `json:"AvailableUpgradeVersions,omitempty"`
	ViableNodePoolVersions   []string `json:"ViableNodePoolVersions,omitempty"`
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
