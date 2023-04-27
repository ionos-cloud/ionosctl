package record

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func RecordsGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve all records",
		Example:   "ionosctl dns record list ",
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

	cmd.AddStringFlag(filter.zoneId, "", "", "Filter used to fetch only the records that contain specified zoneId")
	cmd.AddStringFlag(filter.name, "", "", "Filter used to fetch only the records that contain specified record name")
	cmd.AddIntFlag(offset, "", 0, "The first element (of the total list of elements) to include in the response. Use together with limit for pagination")
	cmd.AddIntFlag(limit, "", 0, "The maximum number of elements to return. Use together with offset for pagination")
	cmd.AddStringSliceFlag(constants.FlagItems, "", []string{}, "")
	cmd.AddFloat64Flag(constants.FlagLimit, "", 0.0, "Pagination limit")
	cmd.AddFloat64Flag(constants.FlagOffset, "", 0.0, "Pagination offset")

	return cmd
}
