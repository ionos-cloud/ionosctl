package registry

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var regPutProperties = sdkgo.PostRegistryProperties{}

func RegReplaceCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "replace",
			Aliases:    []string{"r", "rep"},
			ShortDesc:  "Replace a registry",
			LongDesc:   "Create/replace a registry to hold container images or OCI compliant artifacts",
			Example:    "ionosctl container-registry registry replace --id [REGISTRY_ID] --name [REGISTRY_NAME] --location [REGISTRY_LOCATION]",
			PreCmdRun:  PreCmdPut,
			CmdRun:     CmdPut,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(FlagName, "n", "", "Specify the name of the registry", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagLocation, "", "", "Specify the location of the registry", core.RequiredFlagOption())

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

func CmdPut(c *core.CommandConfig) error {
	var name, location string

	name, err := c.Command.Command.Flags().GetString(FlagName)
	if err != nil {
		return err
	}
	location, err = c.Command.Command.Flags().GetString(FlagLocation)
	if err != nil {
		return err
	}

	id, err := c.Command.Command.Flags().GetString(FlagRegId)
	if err != nil {
		return err
	}
	v := sdkgo.NewWeeklyScheduleWithDefaults()

	if viper.IsSet(core.GetFlagName(c.NS, FlagRegGCDays)) {
		days := viper.GetStringSlice(core.GetFlagName(c.NS, FlagRegGCDays))
		var daysSdk = []sdkgo.Day{}
		for _, day := range days {
			daysSdk = append(daysSdk, sdkgo.Day(day))
		}
		v.SetDays(daysSdk)
	}
	if viper.IsSet(core.GetFlagName(c.NS, FlagRegGCTime)) {
		*v.Time = viper.GetString(core.GetFlagName(c.NS, FlagRegGCTime))
	} else {
		v.SetTime("01:23:00+00:00")
	}

	regPutProperties.SetName(name)
	regPutProperties.SetLocation(location)
	regPutProperties.SetGarbageCollectionSchedule(*v)

	var putInput = sdkgo.PutRegistryInput{}
	putInput.SetProperties(regPutProperties)

	reg, _, err := c.ContainerRegistryServices.Registry().Put(id, putInput)
	if err != nil {
		return err
	}

	regPrint := sdkgo.NewRegistryResponseWithDefaults()
	regPrint.SetProperties(*reg.GetProperties())

	return c.Printer.Print(getRegistryPrint(nil, c, &[]sdkgo.RegistryResponse{}, false))
}

func PreCmdPut(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagRegId, FlagName, FlagLocation)
	if err != nil {
		return err
	}

	return nil
}
