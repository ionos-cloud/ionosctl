package name

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func RegNamesCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "registry",
			Verb:      "names",
			Aliases:   []string{"check", "name", "n"},
			ShortDesc: "Check whether a registry name is available",
			LongDesc: `Check whether a desired registry name is still free. Registry names are globally unique across all IONOS customers because they form the public hostname, so this check should be run before 'container-registry registry create'.

A valid name is 3-63 characters, lowercase letters/digits/dashes only, starting with a letter and ending with a letter or digit (regex ^[a-z][-a-z0-9]{1,61}[a-z0-9]$). The command reports whether the name is available, already taken, or invalid.`,
			Example:    "ionosctl container-registry name --name my-registry",
			PreCmdRun:  PreCmdCheck,
			CmdRun:     CmdCheck,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The desired registry name to check. Must be 3-63 chars, lowercase letters/digits/dashes, starting with a letter (regex ^[a-z][-a-z0-9]{1,61}[a-z0-9]$)", core.RequiredFlagOption())
	return cmd
}

func CmdCheck(c *core.CommandConfig) error {
	name, err := c.Command.Command.Flags().GetString(constants.FlagName)
	if err != nil {
		return err
	}

	res, _ := client.Must().RegistryClient.NamesApi.NamesCheckUsage(context.Background(), name).Execute()
	if res.StatusCode == 404 {
		c.Msg("Name is available")
		return nil
	}

	if res.StatusCode == 400 {
		c.Msg("Name must use only the characters \"a-z\", \"0-9\", or \"-\" " +
			"and starts with a letter and ends with a letter or number " +
			"and is between 3 to 63 characters in length.")
		return nil
	}

	c.Msg("Name is already being used.")
	return nil
}

func PreCmdCheck(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName)
	if err != nil {
		return err
	}

	return nil
}
