package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func resource() *core.Command {
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultResourceCols, utils.ColsMessage(defaultResourceCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(resourceCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = resourceCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultResourceCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, resourceCmd, core.CommandBuilder{
		Namespace:  "resource",
		Resource:   "resource",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Resources",
		LongDesc:   "Use this command to get a full list of existing Resources. To sort list by Resource Type, use `ionosctl resource get` command.",
		Example:    listResourcesExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunResourceList,
		InitClient: true,
	})

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
	getRsc.AddStringFlag(config.ArgType, "", "", "The specific Type of Resources to retrieve information about")
	_ = getRsc.Command.RegisterFlagCompletionFunc(config.ArgType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"datacenter", "snapshot", "image", "ipblock", "pcc", "backupunit", "k8s"}, cobra.ShellCompDirectiveNoFileComp
	})
	getRsc.AddStringFlag(config.ArgResourceId, config.ArgIdShort, "", "The ID of the specific Resource to retrieve information about")
	_ = getRsc.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getResourcesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return resourceCmd
}

func PreRunResourceType(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgType)
}

func RunResourceList(c *core.CommandConfig) error {
	resourcesListed, _, err := c.Users().ListResources()
	if err != nil {
		return err
	}
	return c.Printer.Print(getResourcePrint(c, getResources(resourcesListed)))
}

func RunResourceGet(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, config.ArgResourceId)) {
		resourceListed, _, err := c.Users().GetResourceByTypeAndId(
			viper.GetString(core.GetFlagName(c.NS, config.ArgType)),
			viper.GetString(core.GetFlagName(c.NS, config.ArgResourceId)),
		)
		if err != nil {
			return err
		}
		return c.Printer.Print(getResourcePrint(c, getResource(resourceListed)))
	} else {
		resourcesListed, _, err := c.Users().GetResourcesByType(viper.GetString(core.GetFlagName(c.NS, config.ArgType)))
		if err != nil {
			return err
		}
		return c.Printer.Print(getResourcePrint(c, getResources(resourcesListed)))
	}
}

// Group Resources Commands

func groupResource() *core.Command {
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
	listResources.AddStringSliceFlag(config.ArgCols, "", defaultResourceCols, utils.ColsMessage(defaultResourceCols))
	_ = listResources.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultResourceCols, cobra.ShellCompDirectiveNoFileComp
	})
	listResources.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = listResources.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return resourceCmd
}

func RunGroupResourceList(c *core.CommandConfig) error {
	resourcesListed, _, err := c.Groups().ListResources(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getResourcePrint(c, getResourceGroups(resourcesListed)))
}

// Output Printing

var defaultResourceCols = []string{"ResourceId", "Name", "SecAuthProtection", "Type", "State"}

type ResourcePrint struct {
	ResourceId        string `json:"ResourceId,omitempty"`
	Name              string `json:"Name,omitempty"`
	SecAuthProtection bool   `json:"SecAuthProtection,omitempty"`
	Type              string `json:"Type,omitempty"`
	State             string `json:"State,omitempty"`
}

func getResourcePrint(c *core.CommandConfig, res []resources.Resource) printer.Result {
	r := printer.Result{}
	if c != nil {
		if res != nil {
			r.OutputJSON = res
			r.KeyValue = getResourcesKVMaps(res)
			if c.Resource != c.Namespace {
				r.Columns = getResourceCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
			} else {
				r.Columns = getResourceCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
			}
		}
	}
	return r
}

func getResourceCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var groupCols []string
		columnsMap := map[string]string{
			"ResourceId":        "ResourceId",
			"Name":              "Name",
			"SecAuthProtection": "SecAuthProtection",
			"Type":              "Type",
			"State":             "State",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				groupCols = append(groupCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return groupCols
	} else {
		return defaultResourceCols
	}
}

func getResource(res *resources.Resource) []resources.Resource {
	ress := make([]resources.Resource, 0)
	if res != nil {
		ress = append(ress, resources.Resource{Resource: res.Resource})
	}
	return ress
}

func getResources(groups resources.Resources) []resources.Resource {
	u := make([]resources.Resource, 0)
	if items, ok := groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.Resource{Resource: item})
		}
	}
	return u
}

func getResourceGroups(groups resources.ResourceGroups) []resources.Resource {
	u := make([]resources.Resource, 0)
	if items, ok := groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.Resource{Resource: item})
		}
	}
	return u
}

func getResourcesKVMaps(rs []resources.Resource) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(rs))
	for _, r := range rs {
		var rPrint ResourcePrint
		if id, ok := r.GetIdOk(); ok && id != nil {
			rPrint.ResourceId = *id
		}
		if properties, ok := r.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				rPrint.Name = *name
			}
			if sh, ok := properties.GetSecAuthProtectionOk(); ok && sh != nil {
				rPrint.SecAuthProtection = *sh
			}
		}
		if typeResource, ok := r.GetTypeOk(); ok && typeResource != nil {
			rPrint.Type = string(*typeResource)
		}
		if metadata, ok := r.GetMetadataOk(); ok && metadata != nil {
			if state, ok := metadata.GetStateOk(); ok && state != nil {
				rPrint.State = *state
			}
		}
		o := structs.Map(rPrint)
		out = append(out, o)
	}
	return out
}

func getResourcesIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	clierror.CheckError(err, outErr)
	userSvc := resources.NewUserService(clientSvc.Get(), context.TODO())
	res, _, err := userSvc.ListResources()
	clierror.CheckError(err, outErr)
	resIds := make([]string, 0)
	if items, ok := res.Resources.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				resIds = append(resIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return resIds
}
