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

var regPutProperties = containerregistry.PostRegistryProperties{}

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

	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "", "Specify the name of the registry", core.RequiredFlagOption(),
	)
	cmd.AddStringFlag(constants.FlagLocation, "", "", "Specify the location of the registry", core.RequiredFlagOption())

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
	cmd.AddBoolFlag(
		constants.FlagRegistryVulnScan, "", true, "Enable/disable (?) vulnerability scanning (this is a paid add-on)",
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
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

	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return err
	}

	location, err = c.Command.Command.Flags().GetString(constants.FlagLocation)
	if err != nil {
		return err
	}

	id, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
	if err != nil {
		return err
	}

	v := containerregistry.NewWeeklyScheduleWithDefaults()

	if viper.IsSet(core.GetFlagName(c.NS, FlagRegGCDays)) {
		days := viper.GetStringSlice(core.GetFlagName(c.NS, FlagRegGCDays))
		var daysSdk = []containerregistry.Day{}

		for _, day := range days {
			daysSdk = append(daysSdk, containerregistry.Day(day))
		}

		v.SetDays(daysSdk)
	}

	if viper.IsSet(core.GetFlagName(c.NS, FlagRegGCTime)) {
		v.Time = viper.GetString(core.GetFlagName(c.NS, FlagRegGCTime))
	} else {
		v.SetTime("01:23:00+00:00")
	}

	feat := containerregistry.NewRegistryFeaturesWithDefaults()
	featEnabled := viper.GetBool(core.GetFlagName(c.NS, constants.FlagRegistryVulnScan))
	feat.SetVulnerabilityScanning(containerregistry.FeatureVulnerabilityScanning{Enabled: featEnabled})

	regPutProperties.SetName(name)
	regPutProperties.SetLocation(location)
	regPutProperties.SetGarbageCollectionSchedule(*v)
	regPostProperties.SetFeatures(*feat)

	var putInput = containerregistry.PutRegistryInput{}
	putInput.SetProperties(regPutProperties)

	reg, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesPut(context.Background(), id).PutRegistryInput(putInput).Execute()
	if err != nil {
		return err
	}

	regPrint := containerregistry.NewRegistryResponseWithDefaults()
	regPrint.SetProperties(reg.GetProperties())

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

func PreCmdPut(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(
		c.Command, c.NS, constants.FlagRegistryId, constants.FlagName, constants.FlagLocation,
	)
	if err != nil {
		return err
	}

	return nil
}
