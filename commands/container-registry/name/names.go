package name

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
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

	cmd.AddStringFlag("name", "", "", "Name to check availability for", core.RequiredFlagOption())
	return cmd
}

func CmdCheck(c *core.CommandConfig) error {
	name, err := c.Command.Command.Flags().GetString("name")
	if err != nil {
		return err
	}
	res, _ := c.ContainerRegistryServices.Name().Head(name)
	if res.StatusCode == 404 {
		return c.Printer.Print("Name is available.")
	}
	if res.StatusCode == 400 {
		return c.Printer.Print("Name must use only the characters \"a-z\", \"0-9\", or \"-\" " +
			"and starts with a letter and ends with a letter or number " +
			"and is between 3 to 63 characters in length.")
	}
	return c.Printer.Print("Name is already being used.")
}

func PreCmdCheck(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "name")
	if err != nil {
		return err
	}

	return nil
}
