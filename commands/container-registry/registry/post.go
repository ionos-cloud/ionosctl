package registry

import (
	"context"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
)

var regPostProperties = sdkgo.PostRegistryProperties{}

func RegPostCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "create",
			Aliases:    []string{"c"},
			ShortDesc:  "Create a registry",
			LongDesc:   "Create a registry to hold container images or OCI compliant artifacts",
			Example:    "ionosctl container-registry registry create",
			PreCmdRun:  PreCmdPost,
			CmdRun:     CmdPost,
			InitClient: true,
		},
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag("name", "n", "", "Specify name of the certificate", core.RequiredFlagOption())
	cmd.AddStringFlag("location", "", "", "Specify the certificate itself", core.RequiredFlagOption())

	return cmd
}

func PreCmdPost(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "name", "location")
	if err != nil {
		return err
	}

	return nil
}

func CmdPost(c *core.CommandConfig) error {
	var name, location string

	name, err := c.Command.Command.Flags().GetString("name")
	if err != nil {
		return err
	}
	location, err = c.Command.Command.Flags().GetString("location")
	if err != nil {
		return err
	}

	v := sdkgo.NewWeeklyScheduleWithDefaults()
	v.SetTime("01:23:00+00:00")
	v.SetDays([]sdkgo.Day{"Monday"})

	regPostProperties.SetName(name)
	regPostProperties.SetLocation(location)
	regPostProperties.SetGarbageCollectionSchedule(*v)

	regPostInput := sdkgo.NewPostRegistryInputWithDefaults()

	regPostInput.SetProperties(regPostProperties)

	reg, _, err := c.ContainerRegistryServices.Registry().Post(*regPostInput)
	if err != nil {
		return err
	}

	regPrint := sdkgo.NewRegistryResponseWithDefaults()
	regPrint.SetProperties(*reg.GetProperties())

	return c.Printer.Print(getRegistryPrint(nil, c, &[]sdkgo.RegistryResponse{*regPrint}))
}
