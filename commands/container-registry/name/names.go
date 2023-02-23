package name

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"strings"
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
	res, _ := c.ContainerRegistryServices.Name().Get(name)
	if strings.Contains(res.Status, "404") {
		fmt.Println("name is available")
	} else if strings.Contains(res.Status, "400") {
		fmt.Println("bad request")
	} else {
		fmt.Println("name is already being used")
	}
	return nil
}

func PreCmdCheck(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "name")
	if err != nil {
		return err
	}

	return nil
}
