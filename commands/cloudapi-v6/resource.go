package commands

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultResourceCols = []string{"ResourceId", "Name", "SecAuthProtection", "Type", "State"}
)

func ResourceCmd() *core.Command {
	ctx := context.TODO()
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "resource",
			Aliases:          []string{"res"},
			Short:            "Resource Operations",
			Long:             "The sub-commands of `ionosctl resource` allow you to list, get Resources.",
			TraverseChildren: true,
		},
	}
	globalFlags := resourceCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultResourceCols, tabheaders.ColsMessage(defaultResourceCols))
	_ = viper.BindPFlag(core.GetFlagName(resourceCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = resourceCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultResourceCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, resourceCmd, core.CommandBuilder{
		Namespace:  "resource",
		Resource:   "resource",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Resources",
		LongDesc:   "Use this command to get a full list of existing Resources. To sort list by Resource Type, use `ionosctl resource get` command.",
		Example:    listResourcesExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunResourceList,
		InitClient: true,
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

	/*
		Get Command
	*/
	getRsc := core.NewCommand(ctx, resourceCmd, core.CommandBuilder{
		Namespace:  "resource",
		Resource:   "resource",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get all Resources of a Type or a specific Resource Type",
		LongDesc:   "Use this command to get all Resources of a Type or a specific Resource Type using its Type and ID.\n\nRequired values to run command:\n\n* Type",
		Example:    getResourceExample,
		PreCmdRun:  PreRunResourceType,
		CmdRun:     RunResourceGet,
		InitClient: true,
	})
	getRsc.AddStringFlag(constants.FlagType, "", "", "The specific Type of Resources to retrieve information about", core.RequiredFlagOption())
	_ = getRsc.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"datacenter", "snapshot", "image", "ipblock", "pcc", "backupunit", "k8s"}, cobra.ShellCompDirectiveNoFileComp
	})
	getRsc.AddUUIDFlag(cloudapiv6.ArgResourceId, cloudapiv6.ArgIdShort, "", "The ID of the specific Resource to retrieve information about")
	_ = getRsc.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ResourcesIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return core.WithConfigOverride(resourceCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func PreRunResourceType(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagType)
}

func RunResourceList(c *core.CommandConfig) error {
	resourcesListed, resp, err := c.CloudApiV6Services.Users().ListResources()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Resource, resourcesListed.Resources,
		tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunResourceGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Resource with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId))))

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)) {
		resourceListed, resp, err := c.CloudApiV6Services.Users().GetResourceByTypeAndId(
			viper.GetString(core.GetFlagName(c.NS, constants.FlagType)),
			viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgResourceId)),
		)
		if resp != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
		}
		if err != nil {
			return err
		}

		out, err := jsontabwriter.GenerateOutput("", jsonpaths.Resource, resourceListed.Resource,
			tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
		return nil
	}

	resourcesListed, resp, err := c.CloudApiV6Services.Users().GetResourcesByType(viper.GetString(core.GetFlagName(c.NS, constants.FlagType)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Resource, resourcesListed.Resources,
		tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

// Group Resources Commands

func GroupResourceCmd() *core.Command {
	ctx := context.TODO()
	resourceCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "resource",
			Aliases:          []string{"res"},
			Short:            "Group Resource Operations",
			Long:             "The sub-command of `ionosctl group resource` allows you to list Resources from a Group.",
			TraverseChildren: true,
		},
	}

	/*
		List Resources Command
	*/
	listResources := core.NewCommand(ctx, resourceCmd, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "resource",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Resources from a Group",
		LongDesc:   "Use this command to get a list of Resources assigned to a Group. To see more details about existing Resources, use `ionosctl resource` commands.\n\nRequired values to run command:\n\n* Group Id",
		Example:    listGroupResourcesExample,
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupResourceList,
		InitClient: true,
	})
	listResources.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	listResources.AddStringSliceFlag(constants.ArgCols, "", defaultResourceCols, tabheaders.ColsMessage(defaultResourceCols))
	_ = listResources.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultResourceCols, cobra.ShellCompDirectiveNoFileComp
	})
	listResources.AddUUIDFlag(cloudapiv6.ArgGroupId, "", "", cloudapiv6.GroupId, core.RequiredFlagOption())
	_ = listResources.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.GroupsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return core.WithConfigOverride(resourceCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}

func RunGroupResourceList(c *core.CommandConfig) error {
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Listing Resources from Group with ID: %v...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))))

	resourcesListed, resp, err := c.CloudApiV6Services.Groups().ListResources(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)), listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Resource, resourcesListed.ResourceGroups,
		tabheaders.GetHeadersAllDefault(defaultResourceCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
