package registry

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var regPostProperties = containerregistry.PostRegistryProperties{}

func RegPostCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "registry",
			Verb:      "create",
			Aliases:   []string{"c"},
			ShortDesc: "Create a registry",
			LongDesc: `Create a new Container Registry instance to hold Docker images and OCI artifacts.

The --name becomes the globally-unique hostname prefix and must be available across all IONOS customers, so check it first with 'container-registry name --name <name>'. The --location is fixed at creation and cannot be changed later (use 'container-registry locations' to list valid IDs, e.g. de/txl).

Garbage collection (--garbage-collection-schedule-days / --garbage-collection-schedule-time) is a recurring maintenance run that reclaims storage from untagged and deleted artifacts; pick a low-traffic window. Vulnerability scanning (--vuln-scan) is a paid add-on, enabled by default.

Once the registry is AVAILABLE, authenticate with 'docker login <hostname>' using a token created via 'container-registry token create'.`,
			Example: `# Create a registry with defaults (vulnerability scanning on, random GC window Mon-Fri 10:00-16:00)
ionosctl container-registry registry create --name my-registry --location de/txl

# Create a registry with an explicit weekend GC window and vulnerability scanning disabled
ionosctl container-registry registry create --name my-registry --location de/txl --garbage-collection-schedule-days Saturday,Sunday --garbage-collection-schedule-time "02:00:00Z" --vuln-scan=false`,
			PreCmdRun:  PreCmdPost,
			CmdRun:     CmdPost,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "", "The name of the registry. Becomes the hostname prefix and must be globally unique across all IONOS customers. Lowercase letters, digits and dashes only, 3-63 chars, starting with a letter (regex ^[a-z][-a-z0-9]{1,61}[a-z0-9]$). Check availability with 'container-registry name'", core.RequiredFlagOption(),
	)
	cmd.AddStringFlag(constants.FlagLocation, constants.FlagLocationShort, "", "The location that will host the registry, e.g. de/txl. Fixed at creation - it cannot be changed later. See 'container-registry locations' for valid IDs", core.RequiredFlagOption())
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
		FlagRegGCDays, "", []string{workingDaysOfWeek[rand.Intn(len(workingDaysOfWeek))]}, "Weekly days on which garbage collection runs to reclaim storage from untagged/deleted artifacts. "+
			"Comma-separated full weekday names (Monday...Sunday). Defaults to a random day Mon-Fri",
	)
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagRegGCDays,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagRegGCTime, "", fmt.Sprintf("%02d:00:00Z", hour), "UTC time of day at which garbage collection runs, as an RFC3339 partial-time. "+
		"e.g. \"16:00:00Z\" or \"01:23:00+00:00\". Defaults to a random hour in 10:00-16:00")
	cmd.AddBoolFlag(
		constants.FlagRegistryVulnScan, "", true, "Enable vulnerability scanning of pushed artifacts. This is a paid add-on; enabled by default",
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

	return c.Printer(allCols).Print(reg)
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
