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
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Resources Commands

func resource() *builder.Command {
	ctx := context.TODO()
	resourceCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "resource",
			Short:            "Resource Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl resource` + "`" + ` allows you to list, get Resources.`,
			TraverseChildren: true,
		},
	}
	globalFlags := resourceCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultResourceCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(resourceCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Resources Command
	*/
	builder.NewCommand(ctx, resourceCmd, noPreRun, RunResourcesList, "list", "List Resources",
		"Use this command to get a list of Resources.", listResourcesExample, true)

	/*
		Get Resource Command
	*/
	getRsc := builder.NewCommand(ctx, resourceCmd, PreRunResourceTypeValidate, RunResourcesGet, "get", "Get all Resources of a Type or a specific Resource Type",
		"Use this command to get all Resources of a Type or a specific Resource Type.\n\nRequired values to run command:\n\n* Resource Type",
		getResourceExample, true)
	getRsc.AddStringFlag(config.ArgResourceType, "", "", "The specific Type of Resources to retrieve information about")
	_ = getRsc.Command.RegisterFlagCompletionFunc(config.ArgResourceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"datacenter", "snapshot", "image", "ipblock", "pcc", "backupunit", "k8s"}, cobra.ShellCompDirectiveNoFileComp
	})
	getRsc.AddStringFlag(config.ArgResourceId, "", "", "The ID of the specific Resource to retrieve information about")
	_ = getRsc.Command.RegisterFlagCompletionFunc(config.ArgResourceId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getResourcesIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	return resourceCmd
}

func PreRunResourceTypeValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgResourceType)
}

func RunResourcesList(c *builder.CommandConfig) error {
	resourcesListed, _, err := c.Users().ListResources()
	if err != nil {
		return err
	}
	return c.Printer.Print(getResourcePrint(nil, c, getResources(resourcesListed)))
}

func RunResourcesGet(c *builder.CommandConfig) error {
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceId)) {
		resourceListed, _, err := c.Users().GetResourceByTypeAndId(
			viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceType)),
			viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceId)),
		)
		if err != nil {
			return err
		}
		return c.Printer.Print(getResourcePrint(nil, c, getResource(resourceListed)))
	} else {
		resourcesListed, _, err := c.Users().GetResourcesByType(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgResourceType)))
		if err != nil {
			return err
		}
		return c.Printer.Print(getResourcePrint(nil, c, getResources(resourcesListed)))
	}
}

// Group Resources Commands

func resourceGroup(groupCmd *builder.Command) {
	ctx := context.TODO()

	/*
		List Resources Command
	*/
	listResources := builder.NewCommand(ctx, groupCmd, PreRunGroupIdValidate, RunGroupListResources, "list-resources", "List Resources from a Group",
		"Use this command to get a list of Resources assigned to a Group. To see more details about Resources under a specific User, use `ionosctl user` commands.\n\nRequired values to run command:\n\n* Group Id",
		listGroupResourcesExample, true)
	listResources.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = listResources.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	return
}

func RunGroupListResources(c *builder.CommandConfig) error {
	resourcesListed, _, err := c.Groups().ListResources(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgGroupId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getResourcePrint(nil, c, getResourceGroups(resourcesListed)))
}

// Output Printing

var defaultResourceCols = []string{"ResourceId", "Name", "SecAuthProtection", "Type"}

type ResourcePrint struct {
	ResourceId        string `json:"ResourceId,omitempty"`
	Name              string `json:"Name,omitempty"`
	SecAuthProtection bool   `json:"SecAuthProtection,omitempty"`
	Type              string `json:"Type,omitempty"`
}

func getResourcePrint(resp *resources.Response, c *builder.CommandConfig, groups []resources.Resource) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if groups != nil {
			r.OutputJSON = groups
			r.KeyValue = getResourcesKVMaps(groups)
			r.Columns = getResourceCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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
		viper.GetString(config.ArgServerUrl),
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
