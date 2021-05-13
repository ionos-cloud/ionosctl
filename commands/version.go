package commands

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/core"
)

func version() *core.Command {
	ctx := context.TODO()
	versionCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace:  "version",
		Resource:   "version",
		Verb:       "version",
		ShortDesc:  "Show the current version",
		LongDesc:   "The `ionosctl version` command displays the version of the ionosctl software.",
		Example:    "",
		PreCmdRun:  noPreRun,
		CmdRun:     RunVersion,
		InitClient: false,
	})

	return versionCmd
}

func RunVersion(c *core.CommandConfig) error {
	err := c.Printer.Print("You are currently using ionosctl version: " + rootCmd.Command.Version)
	if err != nil {
		return err
	}
	return nil
}
