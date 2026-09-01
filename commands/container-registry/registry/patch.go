package registry

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var patchInput = containerregistry.PatchRegistryInput{}

func RegUpdateCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "registry",
			Verb:      "update",
			Aliases:   []string{"u", "up"},
			ShortDesc: "Update a registry's garbage-collection schedule (PATCH)",
			LongDesc: `Update the garbage-collection schedule of an existing registry (HTTP PATCH). Garbage collection is the recurring maintenance run that reclaims storage occupied by untagged and deleted artifacts.

Only the days and time of the schedule can be changed here; the registry name, location and features are not touched. Set --garbage-collection-schedule-days to the weekly run days and --garbage-collection-schedule-time to the UTC time of day.`,
			Example: `# Move garbage collection to Monday
ionosctl container-registry registry update --id REGISTRY_ID --garbage-collection-schedule-days Monday

# Run garbage collection on the weekend, early morning UTC
ionosctl container-registry registry update --id REGISTRY_ID --garbage-collection-schedule-days Saturday,Sunday --garbage-collection-schedule-time "03:00:00Z"`,
			PreCmdRun:  PreCmdUpdate,
			CmdRun:     CmdUpdate,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, "i", "", "The unique ID of the registry to update", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringSliceFlag(
		FlagRegGCDays, "", []string{}, "Weekly days on which garbage collection runs. Comma-separated full weekday names (Monday...Sunday)",
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegGCDays,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagRegGCTime, "", "", "UTC time of day for garbage collection, as an RFC3339 partial-time, e.g. \"01:23:00+00:00\"")

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

	return c.Printer(allCols).Print(reg)
}

func PreCmdUpdate(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}
