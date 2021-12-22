package cloudapi_v5

import (
	"context"
	"errors"
	"fmt"
	"go.uber.org/multierr"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/query"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sNodeCmd() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "node",
			Aliases:          []string{"n"},
			Short:            "Kubernetes Node Operations",
			Long:             "The sub-commands of `ionosctl k8s node` allow you to list, get, recreate, delete Kubernetes Nodes.",
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultK8sNodeCols, printer.ColsMessage(defaultK8sNodeCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(k8sCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultK8sNodeCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "node",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Kubernetes Nodes",
		LongDesc:   "Use this command to get a list of existing Kubernetes Nodes.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.K8sNodesFiltersUsage() + "\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id",
		Example:    listK8sNodesExample,
		PreCmdRun:  PreRunK8sNodesList,
		CmdRun:     RunK8sNodeList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv5.ArgK8sNodePoolId, "", "", cloudapiv5.K8sNodePoolId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(list.NS, cloudapiv5.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddIntFlag(cloudapiv5.ArgMaxResults, cloudapiv5.ArgMaxResultsShort, 0, "The maximum number of elements to return")
	list.AddStringFlag(cloudapiv5.ArgOrderBy, "", "", "Limits results to those containing a matching value for a specific property")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv5.ArgFilters, cloudapiv5.ArgFiltersShort, []string{""}, "Limits results to those containing a matching value for a specific property. Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2")
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodesFilters(), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "node",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Kubernetes Node",
		LongDesc:   "Use this command to retrieve details about a specific Kubernetes Node.You can wait for the Node to be in \"ACTIVE\" state using `--wait-for-state` flag together with `--timeout` option.\n\nRequired values to run command:\n\n* K8s Cluster Id\n* K8s NodePool Id\n* K8s Node Id",
		Example:    getK8sNodeExample,
		PreCmdRun:  PreRunK8sClusterNodesIds,
		CmdRun:     RunK8sNodeGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgK8sNodePoolId, "", "", cloudapiv5.K8sNodePoolId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddStringFlag(cloudapiv5.ArgK8sNodeId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sNodeId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodesIds(os.Stderr,
			viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgK8sClusterId)),
			viper.GetString(core.GetFlagName(get.NS, cloudapiv5.ArgK8sNodePoolId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Node to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for Node to be in ACTIVE state [seconds]")

	/*
		Recreate Command
	*/
	recreate := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "recreate",
		Aliases:   []string{"r"},
		ShortDesc: "Recreate a Kubernetes Node",
		LongDesc: `You can recreate a single Kubernetes Node.

Managed Kubernetes starts a process which based on the NodePool's template creates & configures a new Node, waits for status "ACTIVE", and migrates all the Pods from the faulty Node, deleting it once empty. While this operation occurs, the NodePool will have an extra billable "ACTIVE" Node.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id`,
		Example:    recreateK8sNodeExample,
		PreCmdRun:  PreRunK8sClusterNodesIds,
		CmdRun:     RunK8sNodeRecreate,
		InitClient: true,
	})
	recreate.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = recreate.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	recreate.AddStringFlag(cloudapiv5.ArgK8sNodePoolId, "", "", cloudapiv5.K8sNodePoolId, core.RequiredFlagOption())
	_ = recreate.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(recreate.NS, cloudapiv5.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	recreate.AddStringFlag(cloudapiv5.ArgK8sNodeId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sNodeId, core.RequiredFlagOption())
	_ = recreate.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodesIds(os.Stderr,
			viper.GetString(core.GetFlagName(recreate.NS, cloudapiv5.ArgK8sClusterId)),
			viper.GetString(core.GetFlagName(recreate.NS, cloudapiv5.ArgK8sNodePoolId)),
		), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Kubernetes Node",
		LongDesc: `This command deletes a Kubernetes Node within an existing Kubernetes NodePool in a Cluster.

Required values to run command:

* K8s Cluster Id
* K8s NodePool Id
* K8s Node Id`,
		Example:    deleteK8sNodeExample,
		PreCmdRun:  PreRunK8sClusterNodesDelete,
		CmdRun:     RunK8sNodeDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgK8sClusterId, "", "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgK8sNodePoolId, "", "", cloudapiv5.K8sNodePoolId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodePoolId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodePoolsIds(os.Stderr, viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgK8sClusterId))), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgK8sNodeId, "", "", cloudapiv5.K8sNodeId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sNodeId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sNodesIds(os.Stderr,
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgK8sClusterId)),
			viper.GetString(core.GetFlagName(deleteCmd.NS, cloudapiv5.ArgK8sNodePoolId)),
		), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all the Kubernetes Nodes within an existing Kubernetes NodePool in a Cluster.")

	return k8sCmd
}

func PreRunK8sNodesList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgK8sNodePoolId); err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgFilters)) {
		return query.ValidateFilters(c, completer.K8sNodesFilters(), completer.K8sNodesFiltersUsage())
	}
	return nil
}

func PreRunK8sClusterNodesIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgK8sNodePoolId, cloudapiv5.ArgK8sNodeId)
}

func PreRunK8sClusterNodesDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgK8sNodePoolId, cloudapiv5.ArgAll},
		[]string{cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgK8sNodePoolId, cloudapiv5.ArgK8sNodeId},
	)
}

func RunK8sNodeList(c *core.CommandConfig) error {
	c.Printer.Verbose("Listing Nodes from K8s NodePool ID: %v from K8s Cluster ID: %v",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	if !structs.IsZero(listQueryParams) {
		c.Printer.Verbose("Query Parameters set: %v", utils.GetPropertiesKVSet(listQueryParams))
	}
	k8ss, resp, err := c.CloudApiV5Services.K8s().ListNodes(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId)),
		listQueryParams,
	)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePrint(c, getK8sNodes(k8ss)))
}

func RunK8sNodeGet(c *core.CommandConfig) error {
	if err := utils.WaitForState(c, waiter.K8sNodeStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeId))); err != nil {
		return err
	}
	c.Printer.Verbose("Getting K8s Node with ID: %v from K8s NodePool ID: %v from K8s Cluster ID: %v......",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	u, resp, err := c.CloudApiV5Services.K8s().GetNode(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sNodePrint(c, getK8sNode(u)))
}

func RunK8sNodeRecreate(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "recreate k8s node"); err != nil {
		return err
	}
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId))
	k8sNodeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeId))
	c.Printer.Verbose("K8sCluster ID: %v, K8sNodePool ID: %v, K8sNode ID: %v",
		k8sClusterId, k8sNodePoolId, k8sNodeId)
	c.Printer.Verbose("Recreating Node...")
	resp, err := c.CloudApiV5Services.K8s().RecreateNode(k8sClusterId, k8sNodePoolId, k8sNodeId)
	if resp != nil {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print("Status: Command node recreate has been successfully executed")
}

func RunK8sNodeDelete(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))
	nodepoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId))
	nodeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodeId))
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll)) {
		if err := DeleteAllK8sNodes(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s node"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting Node with ID: %v from K8s NodePool ID: %v from K8s Cluster ID: %v...",
			nodeId, nodepoolId, clusterId)
		resp, err := c.CloudApiV5Services.K8s().DeleteNode(clusterId, nodepoolId, nodeId)
		if resp != nil {
			c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	}
}

func DeleteAllK8sNodes(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))
	nodepoolId := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sNodePoolId))
	_ = c.Printer.Print("K8sNodes to be deleted:")
	k8sNodes, _, err := c.CloudApiV5Services.K8s().ListNodes(clusterId, nodepoolId, resources.ListQueryParams{})
	if err != nil {
		return err
	}
	if k8sNodesItems, ok := k8sNodes.GetItemsOk(); ok && k8sNodesItems != nil {
		for _, node := range *k8sNodesItems {
			var messageLog string
			if id, ok := node.GetIdOk(); ok && id != nil {
				messageLog = fmt.Sprintf("K8sNode Id: %v", *id)
			}
			if properties, ok := node.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					messageLog = fmt.Sprintf("%v K8sNode Name: %v", messageLog, *name)
				}
			}
			_ = c.Printer.Print(messageLog)
		}
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the K8sNodes"); err != nil {
			return err
		}
		c.Printer.Verbose("Deleting all the K8sNodes...")
		var multiErr error
		for _, dc := range *k8sNodesItems {
			if id, ok := dc.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Starting deleting Node with ID: %v from K8s NodePool ID: %v from K8s Cluster ID: %v...",
					*id, nodepoolId, clusterId)
				resp, err := c.CloudApiV5Services.K8s().DeleteNode(clusterId, nodepoolId, *id)
				if resp != nil && printer.GetId(resp) != "" {
					c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
				}
				if err != nil {
					multiErr = multierr.Append(multiErr, fmt.Errorf(config.DeleteAllAppendErr, c.Resource, *id, err))
					continue
				} else {
					_ = c.Printer.Print(fmt.Sprintf(config.StatusDeletingAll, c.Resource, *id))
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return err
				}
			}
		}
		if multiErr != nil {
			return multiErr
		}
		return nil
	} else {
		return errors.New("could not get items of K8sNodes")
	}
}

// Output Printing

var defaultK8sNodeCols = []string{"NodeId", "Name", "K8sVersion", "PublicIP", "PrivateIP", "State"}

type K8sNodePrint struct {
	NodeId     string `json:"NodeId,omitempty"`
	Name       string `json:"Name,omitempty"`
	K8sVersion string `json:"K8sVersion,omitempty"`
	PublicIP   string `json:"PublicIP,omitempty"`
	PrivateIP  string `json:"PrivateIP,omitempty"`
	State      string `json:"State,omitempty"`
}

func getK8sNodePrint(c *core.CommandConfig, k8ss []resources.K8sNode) printer.Result {
	r := printer.Result{}
	if c != nil {
		if k8ss != nil {
			r.OutputJSON = k8ss
			r.KeyValue = getK8sNodesKVMaps(k8ss)
			r.Columns = getK8sNodeCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
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
			"PrivateIP":  "PrivateIP",
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
			if priv, ok := properties.GetPrivateIPOk(); ok && priv != nil {
				uPrint.PrivateIP = *priv
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
