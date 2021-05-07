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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func k8sNode() *builder.Command {
	ctx := context.TODO()
	k8sCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "node",
			Short:            "Kubernetes Node Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl k8s node` + "`" + ` allow you to list, get, recreate, delete Kubernetes Nodes.`,
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultK8sNodeCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(k8sCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultK8sNodeCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodePoolIds, RunK8sNodeList, "list", "List Kubernetes Nodes",
		"Use this command to get a list of existing Kubernetes Nodes.\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id", listK8sNodesExample, true)
	list.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Name(), list.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodesIds, RunK8sNodeGet, "get", "Get a Kubernetes Node",
		"Use this command to retrieve details about a specific Kubernetes Node.You can wait for the Node to be in \"ACTIVE\" state using `+\"`\"+`--wait-for-state`+\"`\"+` flag together with `+\"`\"+`--timeout`+\"`\"+` option.\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id\n* K8s Node Id",
		getK8sNodeExample, true)
	get.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Name(), get.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(config.ArgK8sNodeId, "", "", config.RequiredFlagK8sNodeId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgK8sNodeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodesIds(os.Stderr,
			viper.GetString(builder.GetFlagName(k8sCmd.Name(), get.Name(), config.ArgK8sClusterId)),
			viper.GetString(builder.GetFlagName(k8sCmd.Name(), get.Name(), config.ArgK8sNodePoolId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, "", config.DefaultWait, "Wait for specified Node to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, "", config.K8sTimeoutSeconds, "Timeout option for waiting for Node to be in ACTIVE state [seconds]")

	/*
		Recreate Command
	*/
	recreate := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodesIds, RunK8sNodeRecreate, "recreate", "Recreate a Kubernetes Node",
		`You can recreate a single Kubernetes Node.

Managed Kubernetes starts a process which based on the NodePool's template creates & configures a new Node, waits for status "ACTIVE", and migrates all the Pods from the faulty Node, deleting it once empty. While this operation occurs, the NodePool will have an extra billable "ACTIVE" Node.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id`, recreateK8sNodeExample, true)
	recreate.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = recreate.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	recreate.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = recreate.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Name(), recreate.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	recreate.AddStringFlag(config.ArgK8sNodeId, "", "", config.RequiredFlagK8sNodeId)
	_ = recreate.Command.RegisterFlagCompletionFunc(config.ArgK8sNodeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodesIds(os.Stderr,
			viper.GetString(builder.GetFlagName(k8sCmd.Name(), recreate.Name(), config.ArgK8sClusterId)),
			viper.GetString(builder.GetFlagName(k8sCmd.Name(), recreate.Name(), config.ArgK8sNodePoolId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(ctx, k8sCmd, PreRunK8sClusterNodesIds, RunK8sNodeDelete, "delete", "Delete a Kubernetes Node",
		`This command deletes a Kubernetes Node within an existing Kubernetes NodePool in a Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id`, deleteK8sNodeExample, true)
	deleteCmd.AddStringFlag(config.ArgK8sClusterId, "", "", config.RequiredFlagK8sClusterId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgK8sNodePoolId, "", "", config.RequiredFlagK8sNodePoolId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodePoolsIds(os.Stderr, viper.GetString(builder.GetFlagName(k8sCmd.Name(), deleteCmd.Name(), config.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(config.ArgK8sNodeId, "", "", config.RequiredFlagK8sNodeId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgK8sNodeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getK8sNodesIds(os.Stderr,
			viper.GetString(builder.GetFlagName(k8sCmd.Name(), deleteCmd.Name(), config.ArgK8sClusterId)),
			viper.GetString(builder.GetFlagName(k8sCmd.Name(), deleteCmd.Name(), config.ArgK8sNodePoolId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

	return k8sCmd
}

func PreRunK8sClusterNodesIds(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgK8sClusterId, config.ArgK8sNodePoolId, config.ArgK8sNodeId)
}

func RunK8sNodeList(c *builder.CommandConfig) error {
	k8ss, _, err := c.K8s().ListNodes(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePrint(c, getK8sNodes(k8ss)))
}

func RunK8sNodeGet(c *builder.CommandConfig) error {
	if err := utils.WaitForState(c, GetStateK8sNode, viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeId))); err != nil {
		return err
	}
	u, _, err := c.K8s().GetNode(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePrint(c, getK8sNode(u)))
}

func RunK8sNodeRecreate(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "recreate k8s node"); err != nil {
		return err
	}
	_, err := c.K8s().RecreateNode(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeId)))
	if err != nil {
		return err
	}
	return c.Printer.Print("Status: Command node recreate has been successfully executed")
}

func RunK8sNodeDelete(c *builder.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s node"); err != nil {
		return err
	}
	_, err := c.K8s().DeleteNode(
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodeId)))
	if err != nil {
		return err
	}
	return c.Printer.Print("Status: Command node delete has been successfully executed")
}

// Wait for State

func GetStateK8sNode(c *builder.CommandConfig, objId string) (*string, error) {
	obj, _, err := c.K8s().GetNode(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sClusterId)),
		viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgK8sNodePoolId)), objId)
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

// Output Printing

var defaultK8sNodeCols = []string{"NodeId", "Name", "K8sVersion", "PublicIP", "State"}

type K8sNodePrint struct {
	NodeId     string `json:"NodeId,omitempty"`
	Name       string `json:"Name,omitempty"`
	K8sVersion string `json:"K8sVersion,omitempty"`
	PublicIP   string `json:"PublicIP,omitempty"`
	State      string `json:"State,omitempty"`
}

func getK8sNodePrint(c *builder.CommandConfig, k8ss []resources.K8sNode) printer.Result {
	r := printer.Result{}
	if c != nil {
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sNodesKVMaps(k8ss)
			r.Columns = getK8sNodeCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getK8sNodeCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var k8sCols []string
		columnsMap := map[string]string{
			"NodeId":     "NodeId",
			"Name":       "Name",
			"K8sVersion": "K8sVersion",
			"PublicIP":   "PublicIP",
			"State":      "State",
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
		return defaultK8sNodeCols
	}
}

func getK8sNodes(k8ss resources.K8sNodes) []resources.K8sNode {
	u := make([]resources.K8sNode, 0)
	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.K8sNode{KubernetesNode: item})
		}
	}
	return u
}

func getK8sNode(u *resources.K8sNode) []resources.K8sNode {
	k8ss := make([]resources.K8sNode, 0)
	if u != nil {
		k8ss = append(k8ss, resources.K8sNode{KubernetesNode: u.KubernetesNode})
	}
	return k8ss
}

func getK8sNodesKVMaps(us []resources.K8sNode) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(us))
	for _, u := range us {
		var uPrint K8sNodePrint
		if id, ok := u.GetIdOk(); ok && id != nil {
			uPrint.NodeId = *id
		}
		if properties, ok := u.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				uPrint.Name = *name
			}
			if v, ok := properties.GetK8sVersionOk(); ok && v != nil {
				uPrint.K8sVersion = *v
			}
			if v, ok := properties.GetPublicIPOk(); ok && v != nil {
				uPrint.PublicIP = *v
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

func getK8sNodesIds(outErr io.Writer, clusterId, nodepoolId string) []string {
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
	k8ss, _, err := k8sSvc.ListNodes(clusterId, nodepoolId)
	clierror.CheckError(err, outErr)
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.KubernetesNodes.GetItemsOk(); ok && items != nil {
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
