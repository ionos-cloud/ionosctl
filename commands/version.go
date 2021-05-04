package commands

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
)

func version() *builder.Command {
	versionCmd := builder.NewCommand(context.TODO(), nil, noPreRun, RunVersion, "version", "Show the current version",
		"The `ionosctl version` command displays the version of the ionosctl software.", "", false)

	return versionCmd
}

func RunVersion(c *builder.CommandConfig) error {
	err := c.Printer.Print("You are currently using ionosctl version: " + rootCmd.Command.Version)
	if err != nil {
		return err
	}
	return nil
}
