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

var regPutProperties = containerregistry.PostRegistryProperties{}

func RegReplaceCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "registry",
			Verb:      "replace",
			Aliases:   []string{"r", "rep"},
			ShortDesc: "Replace a registry's mutable properties (PUT)",
			LongDesc: `Replace the properties of an existing registry (HTTP PUT). Unlike 'update' (PATCH), this sends a full properties object, so any garbage-collection or feature field you omit is reset to its default rather than preserved.

Note: --name and --location identify an existing registry and cannot be changed by this call (the location is immutable; renaming a registry is not supported). Passing values that differ from the current ones will be rejected by the API. To only tweak the garbage-collection schedule, prefer 'container-registry registry update'.`,
			Example: `# Replace a registry, setting a single GC day and a fixed UTC GC time
ionosctl container-registry registry replace --id REGISTRY_ID --name my-registry --location de/txl --garbage-collection-schedule-days Monday --garbage-collection-schedule-time "01:00:00Z"

# Replace and disable vulnerability scanning
ionosctl container-registry registry replace --id REGISTRY_ID --name my-registry --location de/txl --vuln-scan=false`,
			PreCmdRun:  PreCmdPut,
			CmdRun:     CmdPut,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "", "The registry name. Must match the existing registry's name - it cannot be renamed via this call", core.RequiredFlagOption(),
	)
	cmd.AddStringFlag(constants.FlagLocation, "", "", "The registry location, e.g. de/txl. Must match the existing location - it is immutable", core.RequiredFlagOption())

	cmd.AddStringFlag(constants.FlagRegistryId, "i", "", "The unique ID of the registry to replace", core.RequiredFlagOption())
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
				"Modnday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagRegGCTime, "", "", "UTC time of day for garbage collection, as an RFC3339 partial-time, e.g. \"01:23:00+00:00\"")
	cmd.AddBoolFlag(
		constants.FlagRegistryVulnScan, "", true, "Enable vulnerability scanning of pushed artifacts. This is a paid add-on",
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

	return c.Printer(allCols).Print(reg)
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
