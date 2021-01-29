package commands

import (
	"context"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
)

func version() *builder.Command {
	versionCmd := builder.NewCommand(context.TODO(),
		nil,
		RunVersion,
		"version",
		"Show current CLI version",
		"",
		false)

	return versionCmd
}

func RunVersion(c *builder.CommandConfig) error {
	// TODO: to be updated to get the version automatically
	c.Printer.Print("Version: 0.1")
	return nil
}
