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

	cmd.AddStringFlag("name", "n", "", "Specify name of the certificate", core.RequiredFlagOption())
	cmd.AddStringFlag("location", "", "", "Specify the certificate itself", core.RequiredFlagOption())

	cmd.AddStringFlag("registry-id", "i", "", "Registry ID", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringSliceFlag(
		"garbage-collection-schedule-days", "", []string{}, "Specify the garbage collection schedule days",
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"garbage-collection-schedule-days",
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"Modnday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag("garbage-collection-schedule-time", "", "", "Specify the garbage collection schedule time of day")

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

	name, err := c.Command.Command.Flags().GetString("name")
	if err != nil {
		return err
	}
	location, err = c.Command.Command.Flags().GetString("location")
	if err != nil {
		return err
	}

	id, err := c.Command.Command.Flags().GetString("registry-id")
	if err != nil {
		return err
	}
	v := sdkgo.NewWeeklyScheduleWithDefaults()

	if viper.IsSet(core.GetFlagName(c.NS, "garbage-collection-schedule-days")) {
		days := viper.GetStringSlice(core.GetFlagName(c.NS, "garbage-collection-schedule-days"))
		var daysSdk = []sdkgo.Day{}
		for _, day := range days {
			// TODO: remove this default value when it will work with nil
			daysSdk = append(daysSdk, sdkgo.Day(day))
		}
		v.SetDays(daysSdk)
	} else {
		// TODO: remove this default value when it will work with nil
		v.SetDays([]sdkgo.Day{"Monday"})
	}
	if viper.IsSet(core.GetFlagName(c.NS, "garbage-collection-schedule-time")) {
		*v.Time = viper.GetString(core.GetFlagName(c.NS, "garbage-collection-schedule-time"))
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
	err := core.CheckRequiredFlags(c.Command, c.NS, "registry-id", "name", "location")
	if err != nil {
		return err
	}

	return nil
}
