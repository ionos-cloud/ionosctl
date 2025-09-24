package name

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
)

func RegNamesCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "registry",
			Verb:       "names",
			Aliases:    []string{"check", "name", "n"},
			ShortDesc:  "Check if a Registry Name is available",
			LongDesc:   "Check if a Registry Name is available",
			Example:    "ionosctl container-registry name",
			PreCmdRun:  PreCmdCheck,
			CmdRun:     CmdCheck,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name to check availability for", core.RequiredFlagOption())
	return cmd
}

func CmdCheck(c *core.CommandConfig) error {
	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return err
	}

	res, _ := client.Must().RegistryClient.NamesApi.NamesCheckUsage(context.Background(), name).Execute()
	if res.StatusCode == 404 {
		fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Name is available"))
		return nil
	}

	if res.StatusCode == 400 {
		fmt.Fprintf(
			c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput(
				"Name must use only the characters \"a-z\", \"0-9\", or \"-\" "+
					"and starts with a letter and ends with a letter or number "+
					"and is between 3 to 63 characters in length.",
			),
		)
		return nil
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Name is already being used."))
	return nil
}

func PreCmdCheck(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
	if err != nil {
		return err
	}

	return nil
}
