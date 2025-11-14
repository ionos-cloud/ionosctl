package commands

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultIpConsumerCols = []string{"Ip", "NicId", "ServerId", "DatacenterId", "K8sNodePoolId", "K8sClusterId"}
	allIpConsumerCols     = []string{"Ip", "Mac", "NicId", "ServerId", "ServerName", "DatacenterId", "DatacenterName", "K8sNodePoolId", "K8sClusterId"}
)

func IpconsumerCmd() *core.Command {
	ctx := context.TODO()
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "ipconsumer",
			Aliases:          []string{"ipc"},
			Short:            "Ip Consumer Operations",
			Long:             "The sub-command of `ionosctl ipconsumer` allows you to list information about where each IP address from an IpBlock is being used.",
			TraverseChildren: true,
		},
	}
	globalFlags := resourceCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultIpConsumerCols, tabheaders.ColsMessage(allIpConsumerCols))
	_ = viper.BindPFlag(core.GetFlagName(resourceCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = resourceCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allIpConsumerCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	listResources := core.NewCommand(ctx, resourceCmd, core.CommandBuilder{
		Namespace:  "ipconsumer",
		Resource:   "ipconsumer",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List IpConsumers",
		LongDesc:   "Use this command to get a list of Resources where each IP address from an IpBlock is being used.\n\nRequired values to run command:\n\n* IpBlock Id",
		Example:    listIpConsumersExample,
		PreCmdRun:  PreRunIpBlockId,
		CmdRun:     RunIpConsumersList,
		InitClient: true,
	})
	listResources.AddUUIDFlag(cloudapiv6.ArgIpBlockId, "", "", cloudapiv6.IpBlockId, core.RequiredFlagOption())
	_ = listResources.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgIpBlockId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.IpBlocksIds(), cobra.ShellCompDirectiveNoFileComp
	})
	listResources.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

	return core.WithConfigOverride(resourceCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func RunIpConsumersList(c *core.CommandConfig) error {
	ipBlock, resp, err := c.CloudApiV6Services.IpBlocks().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgIpBlockId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := ipBlock.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting ip block properties")
	}

	ipCons, ok := properties.GetIpConsumersOk()
	if !ok || ipCons == nil {
		return fmt.Errorf("error getting ip consumers")
	}

	ipsConsumers := make([]ionoscloud.IpConsumer, 0)
	for _, ip := range *ipCons {
		ipsConsumers = append(ipsConsumers, ip)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.IpConsumer, ipsConsumers,
		tabheaders.GetHeaders(allIpConsumerCols, defaultIpConsumerCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}
