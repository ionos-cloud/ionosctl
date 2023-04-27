package main

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zon",
		Verb:      "create",
		Aliases:   []string{},
		ShortDesc: "Create a zone",
		Example:   "ionosctl dns zon create ",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			/* TODO: Delete/modify me for --all
			 * err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.Flag<Parent>Id}, []string{constants.ArgAll, constants.Flag<Parent>Id})
			 * if err != nil {
			 * 	return err
			 * }
             * */

			// TODO: If no --all, mark individual flags as required

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			// Implement the actual command logic here
		},
		InitClient: true,
	})



	return cmd
}
