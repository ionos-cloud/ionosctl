package main

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "zon",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a zone",
		Example:   "ionosctl dns zon get --zoneId <String>",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			/* TODO: Delete/modify me for --all
			 * err := core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.Flag<Parent>Id}, []string{constants.ArgAll, constants.Flag<Parent>Id})
			 * if err != nil {
			 * 	return err
			 * }
             * */

			// TODO: If no --all, mark individual flags as required
			err = c.Command.Command.MarkFlagRequired("zoneId")
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			// Implement the actual command logic here
		},
		InitClient: true,
	})


	cmd.AddStringFlag(zoneId, "", "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagDescription, "", "", "The hosted zone is used for..")
	cmd.AddBoolFlag(constants.FlagEnabled, "", false, "Users can activate and deactivate zones")
	cmd.AddStringFlag(constants.FlagZoneName, "", "", "The zone name")
	cmd.AddStringFlag(constants.FlagId, "", "", "The zone ID (UUID)")
	cmd.AddStringFlag(constants.FlagLastModifiedDate, "", "", "The date of the last change formatted as yyyy-MM-dd'T'HH:mm:ss.SSS'Z'")
	cmd.AddStringSliceFlag(constants.FlagNameservers, "", []string{}, "The list of nameservers associated to the zone")
	cmd.AddStringFlag(constants.FlagState, "", "", "The list of possible provisioning states in which DNS resource could be at the specific time")
	cmd.AddStringFlag(constants.FlagCreatedDate, "", "", "The date of creation of the zone formatted as yyyy-MM-dd'T'HH:mm:ss.SSS'Z'")

	return cmd
}
