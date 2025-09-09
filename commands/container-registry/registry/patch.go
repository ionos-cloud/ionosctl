package registry

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var patchInput = containerregistry.PatchRegistryInput{}

func RegUpdateCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "update",
			Aliases:    []string{"u", "up"},
			ShortDesc:  "Update the properties of a registry",
			LongDesc:   "Update the \"garbageCollectionSchedule\" time and days of the week for runs of a registry",
			Example:    "ionosctl container-registry registry update --id [REGISTRY_ID]",
			PreCmdRun:  PreCmdUpdate,
			CmdRun:     CmdUpdate,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, "i", "", "Specify the Registry ID", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringSliceFlag(
		FlagRegGCDays, "", []string{}, "Specify the garbage collection schedule days",
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegGCDays,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"Modnday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagRegGCTime, "", "", "Specify the garbage collection schedule time of day")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdUpdate(c *core.CommandConfig) error {
	v := containerregistry.NewWeeklyScheduleWithDefaults()
	id, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	if err != nil {
		return err
	}

	if viper.IsSet(core.GetFlagName(c.NS, "garbage-collection-schedule-days")) {
		days := viper.GetStringSlice(core.GetFlagName(c.NS, "garbage-collection-schedule-days"))
		var daysSdk = []containerregistry.Day{}

		for _, day := range days {
			daysSdk = append(daysSdk, containerregistry.Day(day))
		}

		v.SetDays(daysSdk)
	}

	if viper.IsSet(core.GetFlagName(c.NS, "garbage-collection-schedule-time")) {
		v.Time = viper.GetString(core.GetFlagName(c.NS, "garbage-collection-schedule-time"))
	} else {
		v.SetTime("01:23:00+00:00")
	}

	patchInput.SetGarbageCollectionSchedule(*v)

	reg, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesPatch(context.Background(), id).PatchRegistryInput(patchInput).Execute()
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput(
		"", jsonpaths.ContainerRegistryRegistry, reg, tabheaders.GetHeadersAllDefault(allCols, cols),
	)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func PreCmdUpdate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}
