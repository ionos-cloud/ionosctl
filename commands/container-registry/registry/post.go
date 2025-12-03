package registry

import (
	"context"
	"fmt"
	"math/rand"
	"time"

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

var regPostProperties = containerregistry.PostRegistryProperties{}

func RegPostCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "create",
			Aliases:    []string{"c"},
			ShortDesc:  "Create a registry",
			LongDesc:   "Create a registry to hold container images or OCI compliant artifacts",
			Example:    "ionosctl container-registry registry create --name NAME --location LOCATION",
			PreCmdRun:  PreCmdPost,
			CmdRun:     CmdPost,
			InitClient: true,
		},
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "", "Specify the name of the registry", core.RequiredFlagOption(),
	)
	cmd.AddStringFlag(constants.FlagLocation, constants.FlagLocationShort, "", "Specify the location of the registry", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagLocation,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return getLocForAutoComplete(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hour := 10 + r.Intn(7) // Random hour 10-16
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	cmd.AddStringSliceFlag(
		FlagRegGCDays, "", []string{workingDaysOfWeek[rand.Intn(len(workingDaysOfWeek))]}, "Specify the garbage collection schedule days. "+
			"Defaults to a random day during Mon-Fri, during the hours 10:00-16:00",
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegGCDays,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagRegGCTime, "", fmt.Sprintf("%02d:00:00Z", hour), "Specify the garbage collection schedule time of day using RFC3339 format. "+
		"i.e. \"16:00:00Z\". Defaults to a random day during Mon-Fri, during the hours 10:00-16:00")
	cmd.AddBoolFlag(
		constants.FlagRegistryVulnScan, "", true, "Enable/disable vulnerability scanning (this is a paid add-on)",
	)

	return cmd
}

func PreCmdPost(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagLocation)
	if err != nil {
		return err
	}

	return nil
}

func CmdPost(c *core.CommandConfig) error {
	var name, location string

	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return err
	}

	location, err = c.Command.Command.Flags().GetString(constants.FlagLocation)
	if err != nil {
		return err
	}

	v := containerregistry.NewWeeklyScheduleWithDefaults()

	days := viper.GetStringSlice(core.GetFlagName(c.NS, FlagRegGCDays))
	var daysSdk = []containerregistry.Day{}

	for _, day := range days {
		daysSdk = append(daysSdk, containerregistry.Day(day))
	}

	v.SetDays(daysSdk)
	v.Time = viper.GetString(core.GetFlagName(c.NS, FlagRegGCTime))

	feat := containerregistry.NewRegistryFeaturesWithDefaults()
	featEnabled := viper.GetBool(core.GetFlagName(c.NS, constants.FlagRegistryVulnScan))
	feat.SetVulnerabilityScanning(containerregistry.FeatureVulnerabilityScanning{Enabled: featEnabled})

	regPostProperties.SetName(name)
	regPostProperties.SetLocation(location)
	regPostProperties.SetGarbageCollectionSchedule(*v)
	regPostProperties.SetFeatures(*feat)

	regPostInput := containerregistry.NewPostRegistryInputWithDefaults()
	regPostInput.SetProperties(regPostProperties)

	reg, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesPost(context.Background()).PostRegistryInput(*regPostInput).Execute()
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

func getLocForAutoComplete() []string {
	var locations []string
	locs, _, _ := client.Must().RegistryClient.LocationsApi.LocationsGet(context.Background()).Execute()
	list := locs.GetItems()

	for _, item := range list {
		locations = append(locations, item.GetId())
	}

	return locations
}
