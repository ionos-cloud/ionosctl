package registry

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var patchInput = sdkgo.PatchRegistryInput{}

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

	cmd.AddStringFlag(FlagRegId, "i", "", "Specify the Registry ID", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdUpdate(c *core.CommandConfig) error {
	id, err := c.Command.Command.Flags().GetString(FlagRegId)
	if err != nil {
		return err
	}
	v := sdkgo.NewWeeklyScheduleWithDefaults()

	if viper.IsSet(core.GetFlagName(c.NS, "garbage-collection-schedule-days")) {
		days := viper.GetStringSlice(core.GetFlagName(c.NS, "garbage-collection-schedule-days"))
		var daysSdk = []sdkgo.Day{}
		for _, day := range days {
			daysSdk = append(daysSdk, sdkgo.Day(day))
		}
		v.SetDays(daysSdk)
	}
	if viper.IsSet(core.GetFlagName(c.NS, "garbage-collection-schedule-time")) {
		*v.Time = viper.GetString(core.GetFlagName(c.NS, "garbage-collection-schedule-time"))
	} else {
		v.SetTime("01:23:00+00:00")
	}

	patchInput.SetGarbageCollectionSchedule(*v)
	reg, _, err := c.ContainerRegistryServices.Registry().Patch(id, patchInput)
	if err != nil {
		return err
	}
	return c.Printer.Print(getRegistryPrint(nil, c, &[]sdkgo.RegistryResponse{reg}, false))
}

func PreCmdUpdate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagRegId)
	if err != nil {
		return err
	}

	return nil
}
