package cluster

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "kafka",
			Resource:  "cluster",
			Verb:      "list",
			Aliases:   []string{"ls"},
			ShortDesc: "Retrieve all clusters using pagination and optional filters",
			Example:   `ionosctl kafka c list --location de/txl`,
			PreCmdRun: func(c *core.PreCommandConfig) error {
				if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagLocation); err != nil {
					return err
				}

				return nil
			},
			CmdRun: func(c *core.CommandConfig) error {
				return listClusters(c)
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagFilterName, "", "", "Filter used to fetch only the records that contain specified name.")
	cmd.AddSetFlag(
		constants.FlagFilterState, "", "", []string{"AVAILABLE", "BUSY", "DEPLOYING", "UPDATING", "FAILED_UPDATING", "FAILED", "DESTROYING"},
		"Filter used to fetch only the records that contain specified state.",
	)

	return cmd
}

func listClusters(c *core.CommandConfig) error {
	ls, err := completer.Clusters(
		func(req kafka.ApiClustersGetRequest) (kafka.ApiClustersGetRequest, error) {
			if fn := core.GetFlagName(c.NS, constants.FlagFilterState); viper.IsSet(fn) {
				req = req.FilterState(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagFilterName); viper.IsSet(fn) {
				req = req.FilterName(viper.GetString(fn))
			}
			return req, nil
		},
	)
	if err != nil {
		return fmt.Errorf("failed listing kafka clusters: %w", err)
	}

	return c.Printer(allCols).Prefix("items").Print(ls)
}
