package cloudapi_v5

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v5/waiter"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func K8sCmd() *core.Command {
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "k8s",
			Short:            "Kubernetes Operations",
			Long:             "The sub-commands of `ionosctl k8s` allow you to perform Kubernetes Operations.",
			TraverseChildren: true,
		},
	}
	k8sCmd.AddCommand(K8sVersionCmd())
	k8sCmd.AddCommand(K8sClusterCmd())
	k8sCmd.AddCommand(K8sKubeconfigCmd())
	k8sCmd.AddCommand(K8sNodePoolCmd())
	k8sCmd.AddCommand(K8sNodeCmd())

	return k8sCmd
}

func K8sClusterCmd() *core.Command {
	ctx := context.TODO()
	k8sCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cluster",
			Aliases:          []string{"c"},
			Short:            "Kubernetes Cluster Operations",
			Long:             "The sub-commands of `ionosctl k8s cluster` allow you to list, get, create, update, delete Kubernetes Clusters.",
			TraverseChildren: true,
		},
	}
	globalFlags := k8sCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultK8sClusterCols, printer.ColsMessage(allK8sClusterCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(k8sCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = k8sCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allK8sClusterCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace:  "k8s",
		Resource:   "cluster",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Kubernetes Clusters",
		LongDesc:   "Use this command to get a list of existing Kubernetes Clusters.",
		Example:    listK8sClustersExample,
		PreCmdRun:  core.NoPreRun,
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
		Aliases:    []string{"g"},
		ShortDesc:  "Get a Kubernetes Cluster",
		LongDesc:   "Use this command to retrieve details about a specific Kubernetes Cluster.You can wait for the Cluster to be in \"ACTIVE\" state using `--wait-for-state` flag together with `--timeout` option.\n\nRequired values to run command:\n\n* K8s Cluster Id",
		Example:    getK8sClusterExample,
		PreCmdRun:  PreRunK8sClusterId,
		CmdRun:     RunK8sClusterGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Cluster to be in ACTIVE state")
	get.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for Cluster to be in ACTIVE state [seconds]")

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a Kubernetes Cluster",
		LongDesc: `Use this command to create a new Managed Kubernetes Cluster. Regarding the name for the Kubernetes Cluster, the limit is 63 characters following the rule to begin and end with an alphanumeric character with dashes, underscores, dots, and alphanumerics between. Regarding the Kubernetes Version for the Cluster, if not set via flag, it will be used the default one: ` + "`" + `ionosctl k8s version get` + "`" + `.

You can wait for the Cluster to be in "ACTIVE" state using ` + "`" + `--wait-for-state` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.`,
		Example:    createK8sClusterExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunK8sClusterCreate,
		InitClient: true,
	})
	create.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "UnnamedCluster", "The name for the K8s Cluster")
	create.AddStringFlag(cloudapiv5.ArgK8sVersion, "", "", "The K8s version for the Cluster. If not set, the default one will be used")
	create.AddStringFlag(cloudapiv5.ArgS3Bucket, "", "", "S3 Bucket name configured for K8s usage")
	create.AddStringSliceFlag(cloudapiv5.ArgApiSubnets, "", []string{""}, "Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6")
	create.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Cluster creation to be executed")
	create.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for the new Cluster to be in ACTIVE state")
	create.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for Cluster/Request [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "update",
		Aliases:   []string{"u", "up"},
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
	update.AddStringFlag(cloudapiv5.ArgName, cloudapiv5.ArgNameShort, "", "The name for the K8s Cluster")
	update.AddStringFlag(cloudapiv5.ArgK8sVersion, "", "", "The K8s version for the Cluster")
	update.AddStringFlag(cloudapiv5.ArgS3Bucket, "", "", "S3 Bucket name configured for K8s usage")
	update.AddStringSliceFlag(cloudapiv5.ArgApiSubnets, "", []string{""}, "Access to the K8s API server is restricted to these CIDRs. Cluster-internal traffic is not affected by this restriction. If no allowlist is specified, access is not restricted. If an IP without subnet mask is provided, the default value will be used: 32 for IPv4 and 128 for IPv6. This will overwrite the existing ones")
	update.AddStringFlag(cloudapiv5.ArgK8sMaintenanceDay, "", "", "The day of the week for Maintenance Window has the English day format as following: Monday or Saturday")
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}, cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(cloudapiv5.ArgK8sMaintenanceTime, "", "", "The time for Maintenance Window has the HH:mm:ss format as following: 08:00:00")
	update.AddStringFlag(cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddBoolFlag(config.ArgWaitForState, config.ArgWaitForStateShort, config.DefaultWait, "Wait for specified Cluster to be in ACTIVE state after updating")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for Cluster to be in ACTIVE state after updating [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, k8sCmd, core.CommandBuilder{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete a Kubernetes Cluster",
		LongDesc: `This command deletes a Kubernetes cluster. The cluster cannot contain any NodePools when deleting.

You can wait for Request for the Cluster deletion to be executed using ` + "`" + `--wait-for-request` + "`" + ` flag together with ` + "`" + `--timeout` + "`" + ` option.

Required values to run command:

* K8s Cluster Id`,
		Example:    deleteK8sClusterExample,
		PreCmdRun:  PreRunK8sClusterDelete,
		CmdRun:     RunK8sClusterDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(cloudapiv5.ArgK8sClusterId, cloudapiv5.ArgIdShort, "", cloudapiv5.K8sClusterId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv5.ArgK8sClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.K8sClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Cluster deletion to be executed")
	deleteCmd.AddBoolFlag(cloudapiv5.ArgAll, cloudapiv5.ArgAllShort, false, "Delete all the Kubernetes clusters.")
	deleteCmd.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, cloudapiv5.K8sTimeoutSeconds, "Timeout option for waiting for Request [seconds]")

	return k8sCmd
}

func PreRunK8sClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv5.ArgK8sClusterId)
}

func PreRunK8sClusterDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv5.ArgK8sClusterId},
		[]string{cloudapiv5.ArgAll},
	)
}

func RunK8sClusterList(c *core.CommandConfig) error {
	k8ss, resp, err := c.CloudApiV5Services.K8s().ListClusters()
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sClusters(k8ss)))
}

func RunK8sClusterGet(c *core.CommandConfig) error {
	if err := utils.WaitForState(c, waiter.K8sClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))); err != nil {
		return err
	}
	c.Printer.Verbose("K8s cluster with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	u, resp, err := c.CloudApiV5Services.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
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
	c.Printer.Verbose("Creating K8s Cluster...")
	u, resp, err := c.CloudApiV5Services.K8s().CreateCluster(*newCluster)
	if resp != nil {
		c.Printer.Verbose("Request href: %v ", resp.Header.Get("location"))
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if id, ok := u.GetIdOk(); ok && id != nil {
			if err = utils.WaitForState(c, waiter.K8sClusterStateInterrogator, *id); err != nil {
				return err
			}
			if u, _, err = c.CloudApiV5Services.K8s().GetCluster(*id); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new K8s Cluster id")
		}
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, getK8sCluster(u)))
}

func RunK8sClusterUpdate(c *core.CommandConfig) error {
	oldCluster, _, err := c.CloudApiV5Services.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	if err != nil {
		return err
	}
	newCluster := getK8sClusterInfo(oldCluster, c)
	c.Printer.Verbose("Updating K8s cluster with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
	k8sUpd, resp, err := c.CloudApiV5Services.K8s().UpdateCluster(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)), newCluster)
	if resp != nil {
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) {
		if err = utils.WaitForState(c, waiter.K8sClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))); err != nil {
			return err
		}
		if k8sUpd, _, err = c.CloudApiV5Services.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId))); err != nil {
			return err
		}
	}
	return c.Printer.Print(getK8sClusterPrint(nil, c, getK8sCluster(k8sUpd)))
}

func RunK8sClusterDelete(c *core.CommandConfig) error {
	var resp *resources.Response
	var err error
	allFlag := viper.GetBool(core.GetFlagName(c.NS, cloudapiv5.ArgAll))
	if allFlag {
		resp, err = DeleteAllK8sClusters(c)
		if err != nil {
			return err
		}
	} else {
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete k8s cluster"); err != nil {
			return err
		}
		c.Printer.Verbose("Starting deleting K8s cluster with id: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
		resp, err = c.CloudApiV5Services.K8s().DeleteCluster(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sClusterId)))
		if resp != nil {
			c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
	}
	return c.Printer.Print(getK8sClusterPrint(resp, c, nil))
}

func getNewK8sCluster(c *core.CommandConfig) (*resources.K8sClusterForPost, error) {
	var (
		k8sversion string
		err        error
	)
	proper := resources.K8sClusterPropertiesForPost{}
	proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName)))
	c.Printer.Verbose("Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName)))
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion)) {
		k8sversion = viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion))
		c.Printer.Verbose("Property K8sVersion set: %v", k8sversion)
	} else {
		if k8sversion, err = getK8sVersion(c); err != nil {
			return nil, err
		}
	}
	proper.SetK8sVersion(k8sversion)
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgS3Bucket)) {
		s3buckets := make([]ionoscloud.S3Bucket, 0)
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgS3Bucket))
		s3buckets = append(s3buckets, ionoscloud.S3Bucket{
			Name: &name,
		})
		proper.SetS3Buckets(s3buckets)
		c.Printer.Verbose("Property S3Buckets set: %v", name)
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgApiSubnets)) {
		apiSubnets := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv5.ArgApiSubnets))
		proper.SetApiSubnetAllowList(apiSubnets)
		c.Printer.Verbose("Property ApiSubnetAllowList set: %v", apiSubnets)
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
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgName)) {
			n := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgName))
			propertiesUpdated.SetName(n)
			c.Printer.Verbose("Property Name set: %v", n)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				propertiesUpdated.SetName(*name)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion)) {
			v := viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sVersion))
			propertiesUpdated.SetK8sVersion(v)
			c.Printer.Verbose("Property K8sVersion set: %v", v)
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgS3Bucket)) {
			s3buckets := make([]ionoscloud.S3Bucket, 0)
			for _, name := range viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv5.ArgS3Bucket)) {
				s3buckets = append(s3buckets, ionoscloud.S3Bucket{
					Name: &name,
				})
				c.Printer.Verbose("Property S3Buckets set: %v", name)
			}
			propertiesUpdated.SetS3Buckets(s3buckets)
		} else {
			if bucketsOk, ok := properties.GetS3BucketsOk(); ok && bucketsOk != nil {
				propertiesUpdated.SetS3Buckets(*bucketsOk)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgApiSubnets)) {
			apiSubnets := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv5.ArgApiSubnets))
			propertiesUpdated.SetApiSubnetAllowList(apiSubnets)
			c.Printer.Verbose("Property ApiSubnetAllowList set: %v", apiSubnets)
		} else {
			if subnetAllowListOk, ok := properties.GetApiSubnetAllowListOk(); ok && subnetAllowListOk != nil {
				propertiesUpdated.SetApiSubnetAllowList(*subnetAllowListOk)
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceTime)) {
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

func DeleteAllK8sClusters(c *core.CommandConfig) (*resources.Response, error) {
	_ = c.Printer.Print("K8sClusters to be deleted:")
	k8Clusters, resp, err := c.CloudApiV5Services.K8s().ListClusters()
	if err != nil {
		return nil, err
	}
	if k8sClustersItems, ok := k8Clusters.GetItemsOk(); ok && k8sClustersItems != nil {
		for _, k8sCluster := range *k8sClustersItems {
			if id, ok := k8sCluster.GetIdOk(); ok && id != nil {
				_ = c.Printer.Print("K8sCluster Id: " + *id)
			}
			if properties, ok := k8sCluster.GetPropertiesOk(); ok && properties != nil {
				if name, ok := properties.GetNameOk(); ok && name != nil {
					_ = c.Printer.Print(" K8sCluster Name: " + *name)
				}
			}
		}

		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete all the K8sClusters"); err != nil {
			return nil, err
		}
		c.Printer.Verbose("Deleting all the K8sClusters...")

		for _, k8sCluster := range *k8sClustersItems {
			if id, ok := k8sCluster.GetIdOk(); ok && id != nil {
				c.Printer.Verbose("Starting deleting K8sCluster with id: %v...", *id)
				resp, err = c.CloudApiV5Services.K8s().DeleteCluster(*id)
				if resp != nil {
					c.Printer.Verbose("Request Id: %v", printer.GetId(resp))
					c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
				}
				if err != nil {
					return nil, err
				}
				if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
					return nil, err
				}
			}
		}
	}
	return resp, err
}

// Output Printing

var defaultK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "Public", "State", "MaintenanceWindow"}

var allK8sClusterCols = []string{"ClusterId", "Name", "K8sVersion", "State", "MaintenanceWindow", "AvailableUpgradeVersions", "ViableNodePoolVersions", "S3Bucket", "ApiSubnetAllowList"}

type K8sClusterPrint struct {
	ClusterId                string   `json:"ClusterId,omitempty"`
	Name                     string   `json:"Name,omitempty"`
	K8sVersion               string   `json:"K8sVersion,omitempty"`
	AvailableUpgradeVersions []string `json:"AvailableUpgradeVersions,omitempty"`
	ViableNodePoolVersions   []string `json:"ViableNodePoolVersions,omitempty"`
	MaintenanceWindow        string   `json:"MaintenanceWindow,omitempty"`
	State                    string   `json:"State,omitempty"`
	S3Bucket                 []string `json:"S3Bucket,omitempty"`
	ApiSubnetAllowList       []string `json:"ApiSubnetAllowList,omitempty"`
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
			r.Columns = getK8sClusterCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
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
			"S3Bucket":                 "S3Bucket",
			"ApiSubnetAllowList":       "ApiSubnetAllowList",
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
			if apiSubnetAllowListOk, ok := properties.GetApiSubnetAllowListOk(); ok && apiSubnetAllowListOk != nil {
				uPrint.ApiSubnetAllowList = *apiSubnetAllowListOk
			}
			if s3BucketsOk, ok := properties.GetS3BucketsOk(); ok && s3BucketsOk != nil {
				s3Buckets := make([]string, 0)
				for _, s3BucketOk := range *s3BucketsOk {
					if nameOk, ok := s3BucketOk.GetNameOk(); ok && nameOk != nil {
						s3Buckets = append(s3Buckets, *nameOk)
					}
				}
				uPrint.S3Bucket = s3Buckets
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

func getMaintenanceInfo(c *core.CommandConfig, maintenance *resources.K8sMaintenanceWindow) resources.K8sMaintenanceWindow {
	var day, time string
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceDay)) {
		day = viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceDay))
		c.Printer.Verbose("Property DayOfTheWeek of MaintenanceWindow set: %v", day)
	} else {
		if d, ok := maintenance.GetDayOfTheWeekOk(); ok && d != nil {
			day = *d
		}
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceTime)) {
		time = viper.GetString(core.GetFlagName(c.NS, cloudapiv5.ArgK8sMaintenanceTime))
		c.Printer.Verbose("Property Time of MaintenanceWindow set: %v", time)
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
